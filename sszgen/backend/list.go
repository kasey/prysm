package backend

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateList struct {
	*types.ValueList
	targetPackage string
}

var generateListGenerateUnmarshalValueFixedTmpl = `{
	if len({{.SliceName}}) % {{.ElementSize}} != 0 {
		return fmt.Errorf("misaligned bytes: {{.FieldName}} length is %d, which is not a multiple of {{.ElementSize}}", len({{.SliceName}}))
	}
	numElem := len({{.SliceName}}) / {{.ElementSize}}
	if numElem > {{ .MaxSize }} {
		return fmt.Errorf("ssz-max exceeded: {{.FieldName}} has %d elements, ssz-max is {{.MaxSize}}", numElem)
	}
	for {{.LoopVar}} := 0; {{.LoopVar}} < numElem; {{.LoopVar}}++ {
		var tmp {{.TypeName}}
		{{.Initializer}}
		tmpSlice := {{.SliceName}}[{{.LoopVar}}*{{.NestedFixedSize}}:(1+{{.LoopVar}})*{{.NestedFixedSize}}]
	{{.NestedUnmarshal}}
		{{.FieldName}} = append({{.FieldName}}, tmp)
	}
}`

var generateListGenerateUnmarshalValueVariableTmpl = `{
// empty lists are zero length, so make sure there is room for an offset
// before attempting to unmarshal it
if len({{.SliceName}}) > 3 {
	firstOffset := ssz.ReadOffset({{.SliceName}}[0:4])
	if firstOffset % 4 != 0 {
			return fmt.Errorf("misaligned list bytes: when decoding {{.FieldName}}, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
	}
	listLen := firstOffset / 4
	if listLen > {{.MaxSize}} {
			return fmt.Errorf("ssz-max exceeded: {{.FieldName}} has %d elements, ssz-max is {{.MaxSize}}", listLen)
	}
	listOffsets := make([]uint64, listLen)
	for {{.LoopVar}} := 0; uint64({{.LoopVar}}) < listLen; {{.LoopVar}}++ {
		listOffsets[{{.LoopVar}}] = ssz.ReadOffset({{.SliceName}}[{{.LoopVar}}*4:({{.LoopVar}}+1)*4])
	}
	for {{.LoopVar}} := 0; {{.LoopVar}} < len(listOffsets); {{.LoopVar}}++ {
			var tmp {{.TypeName}}
			{{.Initializer}}
			var tmpSlice []byte
			if {{.LoopVar}}+1 == len(listOffsets) {
				tmpSlice = {{.SliceName}}[listOffsets[{{.LoopVar}}]:]
			} else {
				tmpSlice = {{.SliceName}}[listOffsets[{{.LoopVar}}]:listOffsets[{{.LoopVar}}+1]]
			}
		{{.NestedUnmarshal}}
			{{.FieldName}} = append({{.FieldName}}, tmp)
	}
}
}`

func (g *generateList) generateUnmarshalVariableValue(fieldName string, sliceName string) string {
	loopVar := "i"
	if fieldName[0:1] == "i" && monoCharacter(fieldName) {
		loopVar = fieldName + "i"
	}
	gg := newValueGenerator(g.ElementValue, g.targetPackage)
	vi, ok := gg.(valueInitializer)
	var initializer string
	if ok {
		initializer = vi.initializeValue("tmp")
		if initializer != "" {
			initializer = "tmp = " + initializer
		}
	}
	tmpl, err := template.New("generateListGenerateUnmarshalValueVariableTmpl").Parse(generateListGenerateUnmarshalValueVariableTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, struct{
		LoopVar string
		SliceName string
		ElementSize int
		TypeName string
		FieldName string
		MaxSize int
		Initializer string
		NestedFixedSize int
		NestedUnmarshal string
	}{
		LoopVar: loopVar,
		SliceName: sliceName,
		ElementSize: g.ElementValue.FixedSize(),
		TypeName: fullyQualifiedTypeName(g.ElementValue, g.targetPackage),
		FieldName: fieldName,
		MaxSize: g.MaxSize,
		Initializer: initializer,
		NestedFixedSize: g.ElementValue.FixedSize(),
		NestedUnmarshal: gg.generateUnmarshalValue("tmp", "tmpSlice"),
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (g *generateList) generateUnmarshalFixedValue(fieldName string, sliceName string) string {
	loopVar := "i"
	if fieldName[0:1] == "i" && monoCharacter(fieldName) {
		loopVar = fieldName + "i"
	}
	gg := newValueGenerator(g.ElementValue, g.targetPackage)
	vi, ok := gg.(valueInitializer)
	var initializer string
	if ok {
		initializer = vi.initializeValue("tmp")
		if initializer != "" {
			initializer = "tmp = " + initializer
		}
	}
	tmpl, err := template.New("generateListGenerateUnmarshalValueFixedTmpl").Parse(generateListGenerateUnmarshalValueFixedTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, struct{
		LoopVar string
		SliceName string
		ElementSize int
		TypeName string
		FieldName string
		MaxSize int
		Initializer string
		NestedFixedSize int
		NestedUnmarshal string
	}{
		LoopVar: loopVar,
		SliceName: sliceName,
		ElementSize: g.ElementValue.FixedSize(),
		TypeName: fullyQualifiedTypeName(g.ElementValue, g.targetPackage),
		FieldName: fieldName,
		MaxSize: g.MaxSize,
		Initializer: initializer,
		NestedFixedSize: g.ElementValue.FixedSize(),
		NestedUnmarshal: gg.generateUnmarshalValue("tmp", "tmpSlice"),
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (g *generateList) generateUnmarshalValue(fieldName string, sliceName string) string {
	if g.ElementValue.IsVariableSized() {
		return g.generateUnmarshalVariableValue(fieldName, sliceName)
	} else {
		return g.generateUnmarshalFixedValue(fieldName, sliceName)

	}
}

func (g *generateList) generateFixedMarshalValue(fieldName string) string {
	tmpl := `dst = ssz.WriteOffset(dst, offset)
offset += %s
`
	offset := g.variableSizeSSZ(fieldName)

	return fmt.Sprintf(tmpl, offset)
}

var variableSizedListTmpl = `func() int {
	s := 0
	for _, o := range {{ .FieldName }} {
		s += 4
		s += {{ .SizeComputation }}
	}
	return s
}()`

func (g *generateList) variableSizeSSZ(fieldName string) string {
	if !g.ElementValue.IsVariableSized() {
		return fmt.Sprintf("len(%s) * %d", fieldName, g.ElementValue.FixedSize())
	}

	gg := newValueGenerator(g.ElementValue, g.targetPackage)
	vslTmpl, err := template.New("variableSizedListTmpl").Parse(variableSizedListTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	err = vslTmpl.Execute(buf, struct{
		FieldName string
		SizeComputation string
	}{
		FieldName: fieldName,
		SizeComputation: gg.variableSizeSSZ("o"),
	})
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

var generateVariableMarshalValueTmpl = `if len({{ .FieldName }}) > {{ .MaxSize }} {
		return nil, ssz.ErrListTooBig
}

for _, o := range {{ .FieldName }} {
		if len(o) != {{ .ElementSize }} {
				return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o) 
}`

var tmplVariableOffsetManagement = `{
	offset = 4 * len({{.FieldName}})
	for _, {{.NestedFieldName}} := range {{.FieldName}} {
		dst = ssz.WriteOffset(dst, offset)
		offset += {{.SizeComputation}}
	}
}
`

func variableOffsetManagement(vg valueGenerator, fieldName, nestedFieldName string) string {
	vomt, err := template.New("tmplVariableOffsetManagement").Parse(tmplVariableOffsetManagement)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	err = vomt.Execute(buf, struct{
		FieldName string
		NestedFieldName string
		SizeComputation string
	}{
		FieldName: fieldName,
		NestedFieldName: nestedFieldName,
		SizeComputation: vg.variableSizeSSZ(nestedFieldName),
	})
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

var tmplGenerateMarshalValueList = `if len({{.FieldName}}) > {{.MaxSize}} {
	return nil, ssz.ErrListTooBig
}
{{.OffsetManagement}}{{.MarshalValue}}`

func (g *generateList) generateVariableMarshalValue(fieldName string) string {
	mvTmpl, err := template.New("tmplGenerateMarshalValueList").Parse(tmplGenerateMarshalValueList)
	if err != nil {
		panic(err)
	}
	var marshalValue string
	var offsetMgmt string
	switch g.ElementValue.(type) {
	case *types.ValueByte:
		marshalValue = fmt.Sprintf("dst = append(dst, %s...)", fieldName)
	default:
		nestedFieldName := "o"
		if fieldName[0:1] == "o" && monoCharacter(fieldName) {
			nestedFieldName = fieldName + "o"
		}
		t := `for _, %s := range %s {
	%s
}`
		gg := newValueGenerator(g.ElementValue, g.targetPackage)
		var internal string
		if g.ElementValue.IsVariableSized() {
			vm, ok := gg.(variableMarshaller)
			if !ok {
				panic(fmt.Sprintf("variable size type does not implement variableMarshaller: %v", g.ElementValue))
			}
			internal = vm.generateVariableMarshalValue(nestedFieldName)
			offsetMgmt = variableOffsetManagement(gg, fieldName, nestedFieldName)
		} else {
			internal = gg.generateFixedMarshalValue(nestedFieldName)
		}
		marshalValue = fmt.Sprintf(t, nestedFieldName, fieldName, internal)
	}
	buf := bytes.NewBuffer(nil)
	err = mvTmpl.Execute(buf, struct{
		FieldName string
		MaxSize int
		MarshalValue string
		OffsetManagement string
	}{
		FieldName: fieldName,
		MaxSize: g.MaxSize,
		MarshalValue: marshalValue,
		OffsetManagement: offsetMgmt,
	})
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

var _ valueGenerator = &generateList{}
