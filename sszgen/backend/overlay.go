package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateOverlay struct {
	*types.ValueOverlay
}

func (g *generateOverlay) GenerateSizeSSZ() *generatedCode {
	return nil
}

func (g *generateOverlay) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ methodGenerator = &generateOverlay{}
