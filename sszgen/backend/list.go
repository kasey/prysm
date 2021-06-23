package backend

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateList struct {
	*types.ValueList
}

func (g *generateList) GenerateSizeSSZ() *generatedCode {
	return nil
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

	gg := newMethodGenerator(g.ElementValue)
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

var _ methodGenerator = &generateList{}
