package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateBool struct {
	*types.ValueBool
}

func (g *generateBool) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateBool) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ methodGenerator = &generateBool{}
