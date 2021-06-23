package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateByte struct {
	*types.ValueByte
}

func (g *generateByte) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateByte) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ methodGenerator = &generateByte{}
