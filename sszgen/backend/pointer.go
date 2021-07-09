package backend

import (
	"fmt"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generatePointer struct {
	*types.ValuePointer
	targetPackage string
}

func (g *generatePointer) generateFixedMarshalValue(fieldName string) string {
	gg := newValueGenerator(g.Referent, g.targetPackage)
	return gg.generateFixedMarshalValue(fieldName)
}

func (g *generatePointer) generateUnmarshalValue(fieldName string, s string) string {
	gg := newValueGenerator(g.Referent, g.targetPackage)
	return gg.generateUnmarshalValue(fieldName, "")
}

func (g *generatePointer) initializeValue(fieldName string) string {
	gg := newValueGenerator(g.Referent, g.targetPackage)
	iv, ok := gg.(valueInitializer)
	if ok {
		return iv.initializeValue(fieldName)
	}
	return ""
}

func (g *generatePointer) generateVariableMarshalValue(fieldName string) string {
	gg := newValueGenerator(g.Referent, g.targetPackage)
	vm, ok := gg.(variableMarshaller)
	if !ok {
		panic(fmt.Sprintf("variable size type does not implement variableMarshaller: %v", g.Referent))
	}
	return vm.generateVariableMarshalValue(fieldName)
}

func (g *generatePointer) variableSizeSSZ(fieldName string) string {
	gg := newValueGenerator(g.Referent, g.targetPackage)
	return gg.variableSizeSSZ(fieldName)
}

var _ valueGenerator = &generatePointer{}
