package types

type UintSize int

const (
	Uint8 = 8
	Uint16 = 16
	Uint32 = 32
	Uint64 = 64
	Uint128 = 128
	Uint256 = 256
)

type ValueUint struct {
	Name string
	Size UintSize
}

func (vu *ValueUint) TypeName() string {
	return vu.Name
}

var _ ValRep = &ValueUint{}