package backend

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateVector struct {
	*types.ValueVector
	targetPackage string
}

func (g *generateVector) generateUnmarshalValue(fieldName string, sliceName string) string {
	gg := newValueGenerator(g.ElementValue, g.targetPackage)
	switch g.ElementValue.(type) {
	case *types.ValueByte:
		return fmt.Sprintf("%s = append([]byte{}, %s...)", fieldName, sliceName)
	default:
		loopVar := "i"
		if fieldName[0:1] == "i" && monoCharacter(fieldName) {
			loopVar = fieldName + "i"
		}
		t := `{
	var tmp {{ .TypeName }}
	for {{ .LoopVar }} := 0; {{ .LoopVar }} < {{ .FixedSize }}; {{ .LoopVar }} ++ {
		tmpSlice := {{ .SliceName }}[{{ .LoopVar }}*{{ .NestedFixedSize }}:(1+{{ .LoopVar }})*{{ .NestedFixedSize }}]
{{ .NestedUnmarshal }}
		{{ .FieldName }} = append({{ .FieldName }}, tmp)
	}
}`
		tmpl, err := template.New("tmplgenerateUnmarshalValueDefault").Parse(t)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewBuffer(nil)
		vr := types.ValRep(g.ValueVector)
		err = tmpl.Execute(buf, struct{
			TypeName string
			SliceName string
			FixedSize int
			NestedFixedSize int
			LoopVar string
			NestedUnmarshal string
			FieldName string
		}{
			TypeName: fullyQualifiedTypeName(vr, g.targetPackage),
			SliceName: sliceName,
			FixedSize: g.FixedSize(),
			NestedFixedSize: g.ElementValue.FixedSize(),
			LoopVar: loopVar,
			NestedUnmarshal: gg.generateUnmarshalValue("tmp", "tmpSlice"),
			FieldName: fieldName,
		})
		if err != nil {
			panic(err)
		}
		return string(buf.Bytes())
	}
}

var tmplGenerateMarshalValueVector = `if len({{.FieldName}}) != {{.Size}} {
	return nil, ssz.ErrBytesLength
}
{{.MarshalValue}}`

func (g *generateVector) generateFixedMarshalValue(fieldName string) string {
	mvTmpl, err := template.New("tmplGenerateMarshalValueVector").Parse(tmplGenerateMarshalValueVector)
	if err != nil {
		panic(err)
	}
	var marshalValue string
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
		internal := gg.generateFixedMarshalValue(nestedFieldName)
		marshalValue = fmt.Sprintf(t, nestedFieldName, fieldName, internal)
	}
	buf := bytes.NewBuffer(nil)
	err = mvTmpl.Execute(buf, struct{
		FieldName string
		Size int
		MarshalValue string
	}{
		FieldName: fieldName,
		Size: g.Size,
		MarshalValue: marshalValue,
	})
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

func monoCharacter(s string) bool {
	ch := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] == ch {
			continue
		}
		return false
	}
	return true
}

func (g *generateVector) variableSizeSSZ(fieldName string) string {
	if !g.ElementValue.IsVariableSized() {
		return fmt.Sprintf("len(%s) * %d", fieldName, g.ElementValue.FixedSize())
	}
	return ""
}

var _ valueGenerator = &generateVector{}