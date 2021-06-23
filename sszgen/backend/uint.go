package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateUint struct {
	*types.ValueUint
}

func (g *generateUint) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateUint) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ methodGenerator = &generateUint{}
