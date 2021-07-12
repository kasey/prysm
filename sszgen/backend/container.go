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

var generateMarshalValueContainerTmpl = `if {{.FieldName}} == nil {
	{{.FieldName}} = new({{.TypeName}})
}
if dst, err = {{.FieldName}}.MarshalSSZTo(dst); err != nil {
	return nil, err
}`

func (g *generateContainer) generateUnmarshalValue(fieldName string, sliceName string) string {
	t := `if err = %s.UnmarshalSSZ(%s); err != nil {
		return err
	}`
	return fmt.Sprintf(t, fieldName, sliceName)
}

func (g *generateContainer) generateFixedMarshalValue(fieldName string) string {
	tmpl, err := template.New("generateMarshalValueContainerTmpl").Parse(generateMarshalValueContainerTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	typeName := g.TypeName()
	if g.targetPackage != g.PackagePath() {
		typeName = fmt.Sprintf("%s.%s", importAlias(g.PackagePath()), g.TypeName())
	}
	tmpl.Execute(buf, struct{
		FieldName string
		TypeName string
	}{
		FieldName: fieldName,
		TypeName: typeName,
	})
	return string(buf.Bytes())
}

func (g *generateContainer) generateVariableMarshalValue(fieldName string) string {
	return g.generateFixedMarshalValue(fieldName)
}

func (g *generateContainer) variableSizeSSZ(fieldName string) string {
	return fmt.Sprintf("%s.SizeSSZ()", fieldName)
}


var sizeBodyTmpl = `func ({{.Receiver}} {{.Type}}) SizeSSZ() (size int) {
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

	variableComputations := make([]string, 0)
	for _, c := range g.Contents {
		vg := newValueGenerator(c.Value, g.targetPackage)
		if !c.Value.IsVariableSized() {
			continue
		}
		cv := vg.variableSizeSSZ(fmt.Sprintf("%s.%s", receiverName, c.Key))
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
		FixedSize: g.FixedSize(),
		VariableSize: "\n" + strings.Join(variableComputations, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents),
	}
}

var marshalBodyTmpl = `func ({{.Receiver}} {{.Type}}) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

func ({{.Receiver}} {{.Type}}) MarshalSSZTo(dst []byte) ([]byte, error) {
	offset := {{.FixedSize -}}
{{- .ValueMarshaling }}
{{- .VariableValueMarshaling }}
	return dst, nil
}`

func (g *generateContainer) GenerateMarshalSSZ() *generatedCode {
	sizeTmpl, err := template.New("GenerateMarshalSSZ").Parse(marshalBodyTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)

	marshalValueBlocks := make([]string, 0)
	marshalVariableValueBlocks := make([]string, 0)
	for i, c := range g.Contents {
		mg := newValueGenerator(c.Value, g.targetPackage)
		fieldName := fmt.Sprintf("%s.%s", receiverName, c.Key)
		mv := mg.generateFixedMarshalValue(fieldName)
		marshalValueBlocks = append(marshalValueBlocks, fmt.Sprintf("\n\t// Field %d: %s", i, c.Key))
		marshalValueBlocks = append(marshalValueBlocks, "\t" + mv)
		if !c.Value.IsVariableSized() {
			continue
		}
		_, ok := mg.(variableMarshaller)
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

	sizeTmpl.Execute(buf, struct{
		Receiver string
		Type string
		FixedSize int
		ValueMarshaling string
		VariableValueMarshaling string
	}{
		Receiver: receiverName,
		Type: fmt.Sprintf("*%s", g.TypeName()),
		FixedSize: g.FixedSize(),
		ValueMarshaling: "\n" + strings.Join(marshalValueBlocks, "\n"),
		VariableValueMarshaling: "\n" + strings.Join(marshalVariableValueBlocks, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents),
	}
}

var generateUnmarshalSSZTmpl = `func ({{.Receiver}} {{.Type}}) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size {{ .SizeInequality }} {{ .FixedSize }} {
		return ssz.ErrSize
	}

	{{ .SliceDeclaration }}
	{{ .VariableOffsetDeclarations }}
	{{ .VariableOffsetValidation }}
	{{ .VariableSliceDeclarations }}
{{ .ValueUnmarshaling }}
	return err
}`

func (g *generateContainer) GenerateUnmarshalSSZ() *generatedCode {
	sizeInequality := "!="
	if g.IsVariableSized() {
		sizeInequality = "<"
	}
	//unmarshalVariableBlocks := make([]string, 0)
	offsets := make([]string, len(g.Contents))
	validations := make([]string, 0)
	unmarshalBlocks := make([]string, 0)
	begin := 0
	end := 0
	slices := make([]string, 0)
	variableOffsets := make([]int, 0)
	for i, c := range g.Contents {
		begin = end
		end += c.Value.FixedSize()
		sliceName := fmt.Sprintf("s%d", i)
		if c.Value.IsVariableSized() {
			offsets = append(offsets, fmt.Sprintf("v%d = ssz.ReadOffset(buf[%d:%d])", i, begin, end))
			var prevBoundCheck string
			if len(variableOffsets) == 0 {
				validations = append(validations, fmt.Sprintf("if v%d < %d {\n\treturn ssz.ErrInvalidVariableOffset\n}", i, g.FixedSize()))
			} else {
				prevBoundCheck = fmt.Sprintf("|| v%d > v%d", variableOffsets[len(variableOffsets)-1], i)
			}
			validations = append(validations, fmt.Sprintf("if v%d > size %s{\n\treturn ssz.ErrOffset\n}", i, prevBoundCheck))
			variableOffsets = append(variableOffsets, i)
		} else {
			slices = append(slices, fmt.Sprintf("%s := buf[%d:%d]", sliceName, begin, end))
		}

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

	variableSlices := make([]string, 0)
	for i := 0; i < len(variableOffsets); i++ {
		current := fmt.Sprintf("v%d", variableOffsets[i])
		next := ""
		if i + 1 < len(variableOffsets) {
			next = fmt.Sprintf("v%d", variableOffsets[i+1])
		}
		variableSlices = append(variableSlices , fmt.Sprintf("s%d := buf[%s:%s]", variableOffsets[i], current, next))
	}
	unmTmpl, err := template.New("GenerateUnmarshalSSZTmpl").Parse(generateUnmarshalSSZTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	unmTmpl.Execute(buf, struct{
		Receiver string
		Type string
		SizeInequality string
		FixedSize int
		SliceDeclaration string
		VariableOffsetDeclarations string
		VariableOffsetValidation string
		VariableSliceDeclarations string
		ValueUnmarshaling string
	}{
		Receiver: receiverName,
		Type: fmt.Sprintf("*%s", g.TypeName()),
		SizeInequality: sizeInequality,
		FixedSize: g.FixedSize(),
		SliceDeclaration: strings.Join(slices, "\n"),
		VariableOffsetDeclarations: strings.Join(offsets, "\n"),
		VariableOffsetValidation: strings.Join(validations, "\n"),
		VariableSliceDeclarations: strings.Join(variableSlices, "\n"),
		ValueUnmarshaling: strings.Join(unmarshalBlocks, "\n"),
	})
	return &generatedCode{
		blocks:  []string{string(buf.Bytes())},
		imports: extractImportsFromContainerFields(g.Contents),
	}
}

func (g *generateContainer) initializeValue(fieldName string) string {
	fqType := g.TypeName()
	if g.targetPackage != g.PackagePath() {
		fqType = importAlias(g.PackagePath()) + "." + fqType
	}
	return fmt.Sprintf("new(%s)", fullyQualifiedTypeName(g.ValueContainer, g.targetPackage))
}

var _ methodGenerator = &generateContainer{}
var _ valueGenerator = &generateContainer{}
var _ valueInitializer = &generateContainer{}