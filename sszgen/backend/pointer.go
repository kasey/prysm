package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generatePointer struct {
	*types.ValuePointer
}

func (g *generatePointer) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generatePointer) variableSizeSSZ(fieldName string) string {
	gg := newMethodGenerator(g.Referent)
	return gg.variableSizeSSZ(fieldName)
}

var _ methodGenerator = &generatePointer{}
