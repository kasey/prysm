package backend

import (
	"fmt"
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generateBool struct {
	*types.ValueBool
	targetPackage string
}

func (g *generateBool) coerce() func(string) string {
	return func(fieldName string) string {
		return fmt.Sprintf("%s(%s)", g.TypeName(), fieldName)
	}
}

func (g *generateBool) generateFixedMarshalValue(fieldName string) string {
	return ""
}

func (g *generateBool) generateUnmarshalValue(fieldName string, s string) string {
	return ""
}

func (g *generateBool) variableSizeSSZ(fieldname string) string {
	return ""
}

var _ valueGenerator = &generateBool{}
