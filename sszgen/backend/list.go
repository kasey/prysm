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

var generateListGenerateUnmarshalValueTmpl = `{
	if len({{.SliceName}}) % {{.ElementSize}} != 0 {
		return fmt.Errorf("misaligned bytes: %s length is %d, which is not a multiple of %d", "{{.FieldName}}", len({{.SliceName}}), {{.ElementSize}})
	}
	numElem := len({{.SliceName}}) / {{.ElementSize}}
	if numElem > {{ .MaxSize }} {
		return fmt.Errorf("ssz-max exceeded: %s has %d elements, ssz-max is %d", "{{.FieldName}}", numElem, {{.MaxSize}})
	}
	for {{.LoopVar}} := 0; {{.LoopVar}} < numElem; {{.LoopVar}}++ {
		var tmp {{.TypeName}}
		{{.Initializer}}
		tmpSlice := {{.SliceName}}[{{.LoopVar}}*{{.NestedFixedSize}}:(1+{{.LoopVar}})*{{.NestedFixedSize}}]
{{.NestedUnmarshal}}
		{{.FieldName}} = append({{.FieldName}}, tmp)
	}
}`

func (g *generateList) generateUnmarshalValue(fieldName string, sliceName string) string {
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
	tmpl, err := template.New("generateListGenerateUnmarshalValueTmpl").Parse(generateListGenerateUnmarshalValueTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, struct{
		SliceName string
		ElementSize int
		FieldName string
		MaxSize int
		TypeName string
		LoopVar string
		Initializer string
		NestedFixedSize int
		NestedUnmarshal string
	}{
		SliceName: sliceName,
		ElementSize: g.ElementValue.FixedSize(),
		FieldName: fieldName,
		MaxSize: g.MaxSize,
		TypeName: fullyQualifiedTypeName(g.ElementValue, g.targetPackage),
		LoopVar: loopVar,
		Initializer: initializer,
		NestedFixedSize: g.ElementValue.FixedSize(),
		NestedUnmarshal: gg.generateUnmarshalValue("tmp", "tmpSlice"),
	})
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

func (g *generateList) generateFixedMarshalValue(fieldName string) string {
	tmpl := `dst = ssz.WriteOffset(dst, offset)
offset += %s
`
	//gg := newValueGenerator(g.ElementValue)
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
		SizeComputation: gg.variableSizeSSZ(fieldName),
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
