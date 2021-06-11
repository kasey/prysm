package types

type ValueByte struct {
	Name string
}

func (vb *ValueByte) TypeName() string {
	return vb.Name
}

var _ ValRep = &ValueByte{}