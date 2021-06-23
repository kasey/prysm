package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateUnion struct {
	*types.ValueUnion
}

func (g *generateUnion) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateUnion) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ methodGenerator = &generateUnion{}
