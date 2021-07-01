package backend

import (
	"fmt"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateUint struct {
	*types.ValueUint
	targetPackage string
}

func (g *generateUint) coerce() func(string) string {
	return func(fieldName string) string {
		return fmt.Sprintf("%s(%s)", g.TypeName(), fieldName)
	}
}

func (g *generateUint) generateFixedMarshalValue(fieldName string) string {
	return fmt.Sprintf("dst = ssz.MarshalUint%d(dst, %s)", g.Size, fieldName)
}

func (g *generateUint) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ valueGenerator = &generateUint{}
