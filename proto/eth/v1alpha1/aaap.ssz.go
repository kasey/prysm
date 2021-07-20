package eth

import (
	ssz "github.com/ferranbt/fastssz"
	prysmaticlabs_eth2_types "github.com/prysmaticlabs/eth2-types"
)

func (c *AggregateAttestationAndProof) XXSizeSSZ() int {
	size := 108
	if c.Aggregate == nil {
		c.Aggregate = new(Attestation)
	}
	size += c.Aggregate.SizeSSZ()
	return size
}
func (c *AggregateAttestationAndProof) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *AggregateAttestationAndProof) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 108

	// Field 0: AggregatorIndex
	dst = ssz.MarshalUint64(dst, uint64(c.AggregatorIndex))

	// Field 1: Aggregate
	if c.Aggregate == nil {
		c.Aggregate = new(Attestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Aggregate.SizeSSZ()

	// Field 2: SelectionProof
	if len(c.SelectionProof) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.SelectionProof...)

	// Field 1: Aggregate
	if dst, err = c.Aggregate.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}
func (c *AggregateAttestationAndProof) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 108 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]    // c.AggregatorIndex
	s2 := buf[12:108] // c.SelectionProof

	v1 := ssz.ReadOffset(buf[8:12]) // c.Aggregate
	if v1 < 108 {
		return ssz.ErrInvalidVariableOffset
	}
	if v1 > size {
		return ssz.ErrOffset
	}
	s1 := buf[v1:] // c.Aggregate

	// Field 0: AggregatorIndex
	c.AggregatorIndex = prysmaticlabs_eth2_types.ValidatorIndex(ssz.UnmarshallUint64(s0))

	// Field 1: Aggregate
	c.Aggregate = new(Attestation)
	if err = c.Aggregate.UnmarshalSSZ(s1); err != nil {
		return err
	}

	// Field 2: SelectionProof
	c.SelectionProof = append([]byte{}, s2...)
	return err
}
