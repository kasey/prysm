package types

type ValueBool struct {
	Name string
}

func (vb *ValueBool) TypeName() string {
	return vb.Name
}

var _ ValRep = &ValueBool{}