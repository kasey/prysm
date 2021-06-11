package types

type ValueVector struct {
	ElementValue ValRep
	Size int
}

func (vv *ValueVector) TypeName() string {
	return "[]" + vv.ElementValue.TypeName()
}

var _ ValRep = &ValueVector{}