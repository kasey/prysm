package types

type ValueUnion struct {
	Name string
}

func (vu *ValueUnion) TypeName() string {
	return vu.Name
}

var _ ValRep = &ValueUnion{}