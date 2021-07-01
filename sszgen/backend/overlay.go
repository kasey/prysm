package backend

import (
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateOverlay struct {
	*types.ValueOverlay
	targetPackage string
}

func (g *generateOverlay) generateFixedMarshalValue(fieldName string) string {
	gg := newValueGenerator(g.Underlying, g.targetPackage)
	uc, ok := gg.(coercer)
	if ok {
		return gg.generateFixedMarshalValue(uc.coerce()(fieldName))
	}
	return gg.generateFixedMarshalValue(fieldName)
}

func (g *generateOverlay) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ valueGenerator = &generateOverlay{}
