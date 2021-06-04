package sszgen

import "github.com/prysmaticlabs/prysm/sszgen/types"

type GeneratedCode struct {
	methods map[string]string
	imports map[string]string
}
var _ SSZSatisfier = &GeneratedCode{}

func (gc *GeneratedCode) Imports() map[string]string {
	return gc.imports
}

func (gc *GeneratedCode) Methods() map[string]string {
	return gc.methods
}

func Render(val types.ValRep) (*GeneratedCode, error) {
	return nil, nil
}
