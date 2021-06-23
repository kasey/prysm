package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateVector struct {
	*types.ValueVector
}

func (g *generateVector) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateVector) variableSizeSSZ(fieldname string) string {
	/*
	if !g.ElementValue.IsVariableSized() {
		return nil
	}
	elementGenerator := newMethodGenerator(g.ElementValue)
	return jen.Lit(g.Size).Op("*").Add(elementGenerator.variableSizeSSZ(fieldName))
	 */
	return ""
}

var _ methodGenerator = &generateVector{}