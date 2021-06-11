package types

type ValueContainer struct {
	Name string
	Contents map[string]ValRep
}

func (vc *ValueContainer) TypeName() string {
	return vc.Name
}

var _ ValRep = &ValueContainer{}