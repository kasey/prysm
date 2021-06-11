package types

type ValuePointer struct {
	Referent ValRep
}

func (vp *ValuePointer) TypeName() string {
	return "*" + vp.Referent.TypeName()
}

var _ ValRep = &ValuePointer{}