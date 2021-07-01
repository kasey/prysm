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
	sizeTmpl, err := template.New("generateFixedMarshalValue").Parse(marshalBodyTmpl)
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
var _ valueGenerator = &generateContainer{}
