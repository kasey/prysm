package backend

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateContainer struct {
	*types.ValueContainer
}

const receiverName = "c"

func (g *generateContainer) variableSizeSSZ(fieldName string) string {
	return fmt.Sprintf("%s.SizeSSZ()", fieldName)
}

var generateBodyTmpl = `func ({{.Receiver}} {{.Type}}) SizeSSZ() (size int) {
	size := {{.FixedSize}}
	{{- .VariableSize }}
	return size
}`

func (g *generateContainer) GenerateSizeSSZ() *generatedCode {
	sizeTmpl, err := template.New("GenerateSizeSSZ").Parse(generateBodyTmpl)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)

	variableComputations := make([]string, 0)
	for _, c := range g.Contents {
		cg := newMethodGenerator(c.Value)
		if !c.Value.IsVariableSized() {
			continue
		}
		cv := cg.variableSizeSSZ(fmt.Sprintf("%s.%s", receiverName, c.Key))
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
	}
}

/*
		jen.Id("size").Op("+=").Id("len").Call(jen.Id("b").Dot("HistoricalRoots")).Op("*").Lit(32),
		jen.Id("size").Op("+=").Id("len").Call(jen.Id("b").Dot("Eth1DataVotes")).Op("*").Lit(72),
		jen.Id("size").Op("+=").Id("len").Call(jen.Id("b").Dot("Validators")).Op("*").Lit(121),
		jen.Id("size").Op("+=").Id("len").Call(jen.Id("b").Dot("Balances")).Op("*").Lit(8),
		jen.For(jen.Id("ii").Op(":=").Lit(0),
			jen.Id("ii").Op("<").Id("len").Call(jen.Id("b").Dot("PreviousEpochAttestations")),
			jen.Id("ii").Op("++")).
			Block(jen.Id("size").Op("+=").Lit(4),
				jen.Id("size").Op("+=").Id("b").Dot("PreviousEpochAttestations").Index(jen.Id("ii")).Dot("SizeSSZ").Call()),
		jen.For(jen.Id("ii").Op(":=").Lit(0),
			jen.Id("ii").Op("<").Id("len").Call(jen.Id("b").Dot("CurrentEpochAttestations")),
			jen.Id("ii").Op("++")).Block(jen.Id("size").Op("+=").Lit(4), jen.Id("size").Op("+=").Id("b").Dot("CurrentEpochAttestations").Index(jen.Id("ii")).Dot("SizeSSZ").Call()),
		)
 */

var _ methodGenerator = &generateContainer{}
