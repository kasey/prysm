package backend

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

const receiverName = "c"

type generateContainer struct {
	*types.ValueContainer
	targetPackage string
}

func (g *generateContainer) generateUnmarshalValue(fieldName string, sliceName string) string {
	t := `if err = %s.UnmarshalSSZ(%s); err != nil {
		return err
	}`
	return fmt.Sprintf(t, fieldName, sliceName)
}

func (g *generateContainer) generateFixedMarshalValue(fieldName string) string {
	if g.IsVariableSized() {
		return fmt.Sprintf(`dst = ssz.WriteOffset(dst, offset)
offset += %s.SizeSSZ()`, fieldName)
	}
	return g.generateDelegateFieldMarshalSSZ(fieldName)
}

var generateMarshalValueContainerTmpl = `
`

// method that generates code which calls the MarshalSSZ method of the field
func (g *generateContainer) generateDelegateFieldMarshalSSZ(fieldName string) string {
	return fmt.Sprintf(`if dst, err = %s.MarshalSSZTo(dst); err != nil {
		return nil, err
	}`, fieldName)
}

func (g *generateContainer) generateVariableMarshalValue(fieldName string) string {
	return g.generateDelegateFieldMarshalSSZ(fieldName)
}

func (g *generateContainer) variableSizeSSZ(fieldName string) string {
	return fmt.Sprintf("%s.SizeSSZ()", fieldName)
}


var sizeBodyTmpl = `func ({{.Receiver}} {{.Type}}) XXSizeSSZ() (int) {
	size := {{.FixedSize}}
	{{- .VariableSize }}
	return size
}`

func (g *generateContainer) GenerateSizeSSZ() *generatedCode {
	sizeTmpl, err := template.New("GenerateSizeSSZ").Parse(sizeBodyTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)

	fixedSize := 0
	variableComputations := make([]string, 0)
	for _, c := range g.Contents {
		vg := newValueGenerator(c.Value, g.targetPackage)
		fixedSize += c.Value.FixedSize()
		if !c.Value.IsVariableSized() {
			continue
		}
		fieldName := fmt.Sprintf("%s.%s", receiverName, c.Key)
		vi, ok := vg.(valueInitializer)
		if ok {
			ini := vi.initializeValue(fieldName)
			if ini != "" {
				variableComputations = append(variableComputations, fmt.Sprintf("if %s == nil {\n\t%s = %s\n}", fieldName, fieldName, ini))
			}
		}
		cv := vg.variableSizeSSZ(fieldName)
		if cv != "" {
			variableComputations = append(variableComputations, fmt.Sprintf("\tsize += %s", cv))
		}
	}

	sizeTmpl.Execute(buf, struct{
		Receiver string
		Type string
		FixedSize int
		VariableSize string
	}{
		Receiver: receiverName,
		Type: fmt.Sprintf("*%s", g.TypeName()),
		FixedSize: fixedSize,
		VariableSize: "\n" + strings.Join(variableComputations, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents, g.targetPackage),
	}
}

var marshalBodyTmpl = `func ({{.Receiver}} {{.Type}}) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ({{.Receiver}})
}

func ({{.Receiver}} {{.Type}}) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
{{- .OffsetDeclaration -}}
{{- .ValueMarshaling }}
{{- .VariableValueMarshaling }}
	return dst, err
}`

func (g *generateContainer) GenerateMarshalSSZ() *generatedCode {
	sizeTmpl, err := template.New("GenerateMarshalSSZ").Parse(marshalBodyTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)

	marshalValueBlocks := make([]string, 0)
	marshalVariableValueBlocks := make([]string, 0)
	offset := 0
	for i, c := range g.Contents {
		// only lists need the offset variable
		mg := newValueGenerator(c.Value, g.targetPackage)
		fieldName := fmt.Sprintf("%s.%s", receiverName, c.Key)
		marshalValueBlocks = append(marshalValueBlocks, fmt.Sprintf("\n\t// Field %d: %s", i, c.Key))
		vi, ok := mg.(valueInitializer)
		if ok {
			ini := vi.initializeValue(fieldName)
			if ini != "" {
				marshalValueBlocks = append(marshalValueBlocks , fmt.Sprintf("if %s == nil {\n\t%s = %s\n}", fieldName, fieldName, ini))
			}
		}
		mv := mg.generateFixedMarshalValue(fieldName)
		marshalValueBlocks = append(marshalValueBlocks, "\t" + mv)
		offset += c.Value.FixedSize()
		if !c.Value.IsVariableSized() {
			continue
		}
		_, ok = mg.(variableMarshaller)
		if !ok {
			continue
		}
		vm := mg.(variableMarshaller)
		vmc := vm.generateVariableMarshalValue(fieldName)
		if vmc != "" {
			marshalVariableValueBlocks = append(marshalVariableValueBlocks, fmt.Sprintf("\n\t// Field %d: %s", i, c.Key))
			marshalVariableValueBlocks = append(marshalVariableValueBlocks, "\t" + vmc)
		}
	}
	// only set the offset declaration if we need it
	// otherwise we'll have an unused variable (syntax error)
	offsetDeclaration := ""
	if g.IsVariableSized() {
		// if there are any variable sized values in the container, we'll need to set this offset declaration
		// so it gets rendered to the top of the marshal method
		offsetDeclaration = fmt.Sprintf("\noffset := %d\n", offset)
	}

	sizeTmpl.Execute(buf, struct{
		Receiver string
		Type string
		OffsetDeclaration string
		ValueMarshaling string
		VariableValueMarshaling string
	}{
		Receiver: receiverName,
		Type: fmt.Sprintf("*%s", g.TypeName()),
		OffsetDeclaration: offsetDeclaration,
		ValueMarshaling: "\n" + strings.Join(marshalValueBlocks, "\n"),
		VariableValueMarshaling: "\n" + strings.Join(marshalVariableValueBlocks, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents, g.targetPackage),
	}
}

var generateUnmarshalSSZTmpl = `func ({{.Receiver}} {{.Type}}) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size {{ .SizeInequality }} {{ .FixedOffset }} {
		return ssz.ErrSize
	}

	{{ .SliceDeclaration }}
{{ .ValueUnmarshaling }}
	return err
}`

type unmarshalStep struct {
	valRep types.ValRep
	fieldNumber int
	fieldName string
	beginByte int
	endByte int
	previousVariable *unmarshalStep
	nextVariable *unmarshalStep
}

type unmarshalStepSlice []*unmarshalStep

func (us *unmarshalStep) fixedSize() int {
	return us.valRep.FixedSize()
}

func (us *unmarshalStep) variableOffset(outerFixedSize int) string {
	o := fmt.Sprintf("v%d := ssz.ReadOffset(buf[%d:%d]) // %s", us.fieldNumber, us.beginByte, us.endByte, us.fieldName)
	if us.previousVariable == nil {
		o += fmt.Sprintf("\nif v%d < %d {\n\treturn ssz.ErrInvalidVariableOffset\n}", us.fieldNumber, outerFixedSize)
		o += fmt.Sprintf("\nif v%d > size {\n\treturn ssz.ErrOffset\n}", us.fieldNumber)
	} else {
		o += fmt.Sprintf("\nif v%d > size || v%d < v%d {\n\treturn ssz.ErrOffset\n}", us.fieldNumber, us.fieldNumber, us.previousVariable.fieldNumber)
	}
	return o
}

func (us *unmarshalStep) slice() string {
	if us.valRep.IsVariableSized() {
		if us.nextVariable == nil {
			return fmt.Sprintf("s%d := buf[v%d:]\t\t// %s", us.fieldNumber, us.fieldNumber, us.fieldName)
		}
		return fmt.Sprintf("s%d := buf[v%d:v%d]\t\t// %s", us.fieldNumber, us.fieldNumber, us.nextVariable.fieldNumber, us.fieldName)
	}
	return fmt.Sprintf("s%d := buf[%d:%d]\t\t// %s", us.fieldNumber, us.beginByte, us.endByte, us.fieldName)
}

func (steps unmarshalStepSlice) fixedSlices() string {
	slices := make([]string, 0)
	for _, s := range steps {
		if s.valRep.IsVariableSized() {
			continue
		}
		slices = append(slices, s.slice())
	}
	return strings.Join(slices, "\n")
}

func (steps unmarshalStepSlice)  variableSlices(outerSize int) string {
	validate := make([]string, 0)
	assign := make([]string, 0)
	for _, s := range steps {
		if !s.valRep.IsVariableSized() {
			continue
		}
		validate = append(validate, s.variableOffset(outerSize))
		assign = append(assign, s.slice())
	}
	return strings.Join(append(validate, assign...), "\n")
}

func (g *generateContainer) unmarshalSteps() unmarshalStepSlice{
	ums := make([]*unmarshalStep, 0)
	var begin, end int
	var prevVariable *unmarshalStep
	for i, c := range g.Contents {
		begin = end
		end += c.Value.FixedSize()
		um := &unmarshalStep{
			valRep: c.Value,
			fieldNumber: i,
			fieldName: fmt.Sprintf("%s.%s", receiverName, c.Key),
			beginByte: begin,
			endByte: end,
		}
		if c.Value.IsVariableSized() {
			if prevVariable != nil {
				um.previousVariable = prevVariable
				prevVariable.nextVariable = um
			}
			prevVariable = um
		}

		ums = append(ums, um)
	}
	return ums
}

func (g *generateContainer) GenerateUnmarshalSSZ() *generatedCode {
	sizeInequality := "!="
	if g.IsVariableSized() {
		sizeInequality = "<"
	}
	ums := g.unmarshalSteps()
	unmarshalBlocks := make([]string, 0)
	for i, c := range g.Contents {
		unmarshalBlocks = append(unmarshalBlocks, fmt.Sprintf("\n\t// Field %d: %s", i, c.Key))
		mg := newValueGenerator(c.Value, g.targetPackage)
		fieldName := fmt.Sprintf("%s.%s", receiverName, c.Key)

		vi, ok := mg.(valueInitializer)
		if ok {
			ini := vi.initializeValue(fieldName)
			if ini != "" {
				unmarshalBlocks = append(unmarshalBlocks, fmt.Sprintf("%s = %s", fieldName, ini))
			}
		}

		sliceName := fmt.Sprintf("s%d", i)
		mv := mg.generateUnmarshalValue(fieldName, sliceName)
		if mv != "" {
			//unmarshalBlocks = append(unmarshalBlocks, fmt.Sprintf("\t%s = %s", fieldName, mv))
			unmarshalBlocks = append(unmarshalBlocks, mv)
		}

		/*
				if !c.Value.IsVariableSized() {
					continue
				}
		_, ok := mg.(variableUnmarshaller)
		if !ok {
			continue
		}
		vm := mg.(variableUnmarshaller)
		vmc := vm.generateVariableUnmarshalValue(fieldName)
		if vmc != "" {
			unmarshalVariableBlocks = append(unmarshalVariableBlocks, fmt.Sprintf("\n\t// Field %d: %s", i, c.Key))
			unmarshalVariableBlocks = append(unmarshalVariableBlocks, "\t" + vmc)
		}
		 */
	}

	sliceDeclarations := strings.Join([]string{ums.fixedSlices(), "", ums.variableSlices(g.fixedOffset())}, "\n")
	unmTmpl, err := template.New("GenerateUnmarshalSSZTmpl").Parse(generateUnmarshalSSZTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	unmTmpl.Execute(buf, struct{
		Receiver string
		Type string
		SizeInequality string
		FixedOffset int
		SliceDeclaration string
		ValueUnmarshaling string
	}{
		Receiver: receiverName,
		Type: fmt.Sprintf("*%s", g.TypeName()),
		SizeInequality: sizeInequality,
		FixedOffset: g.fixedOffset(),
		SliceDeclaration: sliceDeclarations,
		ValueUnmarshaling: strings.Join(unmarshalBlocks, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents, g.targetPackage),
	}
}

func (g *generateContainer) fixedOffset() int {
	offset := 0
	for _, c := range g.Contents {
		offset += c.Value.FixedSize()
	}
	return offset
}

func (g *generateContainer) initializeValue(fieldName string) string {
	fqType := g.TypeName()
	if g.targetPackage != g.PackagePath() {
		fqType = importAlias(g.PackagePath()) + "." + fqType
	}
	return fmt.Sprintf("new(%s)", fullyQualifiedTypeName(g.ValueContainer, g.targetPackage))
}

func containsList(v types.ValRep) bool {
	switch t := v.(type) {
	case *types.ValueContainer:
		// we only care about top-level lists, so
		// we don't want to look deeper into containers
		return false
	case *types.ValueList:
		return true
	case *types.ValueOverlay:
		return containsList(t.Underlying)
	case *types.ValuePointer:
		return containsList(t.Referent)
	}
	return false
}

var _ methodGenerator = &generateContainer{}
var _ valueGenerator = &generateContainer{}
var _ valueInitializer = &generateContainer{}