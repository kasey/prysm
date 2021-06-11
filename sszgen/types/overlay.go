package types

type ValueOverlay struct {
	Name string
	Underlying ValRep
}

func (vo *ValueOverlay) TypeName() string {
	return vo.Name
}

var _ ValRep = &ValueOverlay{}