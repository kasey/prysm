package types

type ValueList struct {
	Name string
	ElementValue ValRep
	MaxSize int
	isReference bool
}

func (vl *ValueList) TypeName() string {
	return "[]" + vl.ElementValue.TypeName()
}

func (vl *ValueList) IsReference() bool {
	return vl.isReference
}

var _ ValRep = &ValueList{}
