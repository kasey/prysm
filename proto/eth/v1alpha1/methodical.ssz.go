package eth

import (
	"fmt"
	ssz "github.com/ferranbt/fastssz"
	prysmaticlabs_eth2_types "github.com/prysmaticlabs/eth2-types"
	prysmaticlabs_go_bitfield "github.com/prysmaticlabs/go-bitfield"
)

func (c *Attestation) XXSizeSSZ() int {
	size := 228

	return size
}
func (c *Attestation) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *Attestation) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 228

	// Field 0: AggregationBits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AggregationBits) * 1

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(AttestationData)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}
func (c *Attestation) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 228 {
		return ssz.ErrSize
	}

	s1 := buf[4:132]   // c.Data
	s2 := buf[132:228] // c.Signature

	v0 := ssz.ReadOffset(buf[0:4]) // c.AggregationBits
	if v0 < 228 {
		return ssz.ErrInvalidVariableOffset
	}
	if v0 > size {
		return ssz.ErrOffset
	}
	s0 := buf[v0:] // c.AggregationBits

	// Field 0: AggregationBits
	if err = ssz.ValidateBitlist(s0, 2048); err != nil {
		return err
	}
	c.AggregationBits = append([]byte{}, prysmaticlabs_go_bitfield.Bitlist(s0)...)

	// Field 1: Data
	c.Data = new(AttestationData)
	if err = c.Data.UnmarshalSSZ(s1); err != nil {
		return err
	}

	// Field 2: Signature
	c.Signature = append([]byte{}, s2...)
	return err
}
func (c *AttestationData) XXSizeSSZ() int {
	size := 128

	return size
}
func (c *AttestationData) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *AttestationData) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 1: CommitteeIndex
	dst = ssz.MarshalUint64(dst, uint64(c.CommitteeIndex))

	// Field 2: BeaconBlockRoot
	if len(c.BeaconBlockRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BeaconBlockRoot...)

	// Field 3: Source
	if c.Source == nil {
		c.Source = new(Checkpoint)
	}
	if dst, err = c.Source.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 4: Target
	if c.Target == nil {
		c.Target = new(Checkpoint)
	}
	if dst, err = c.Target.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}
func (c *AttestationData) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 128 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]    // c.Slot
	s1 := buf[8:16]   // c.CommitteeIndex
	s2 := buf[16:48]  // c.BeaconBlockRoot
	s3 := buf[48:88]  // c.Source
	s4 := buf[88:128] // c.Target

	// Field 0: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s0))

	// Field 1: CommitteeIndex
	c.CommitteeIndex = prysmaticlabs_eth2_types.CommitteeIndex(ssz.UnmarshallUint64(s1))

	// Field 2: BeaconBlockRoot
	c.BeaconBlockRoot = append([]byte{}, s2...)

	// Field 3: Source
	c.Source = new(Checkpoint)
	if err = c.Source.UnmarshalSSZ(s3); err != nil {
		return err
	}

	// Field 4: Target
	c.Target = new(Checkpoint)
	if err = c.Target.UnmarshalSSZ(s4); err != nil {
		return err
	}
	return err
}
func (c *AttesterSlashing) XXSizeSSZ() int {
	size := 8
	if c.Attestation_1 == nil {
		c.Attestation_1 = new(IndexedAttestation)
	}
	size += c.Attestation_1.SizeSSZ()
	if c.Attestation_2 == nil {
		c.Attestation_2 = new(IndexedAttestation)
	}
	size += c.Attestation_2.SizeSSZ()
	return size
}
func (c *AttesterSlashing) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *AttesterSlashing) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 8

	// Field 0: Attestation_1
	if c.Attestation_1 == nil {
		c.Attestation_1 = new(IndexedAttestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Attestation_1.SizeSSZ()

	// Field 1: Attestation_2
	if c.Attestation_2 == nil {
		c.Attestation_2 = new(IndexedAttestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Attestation_2.SizeSSZ()

	// Field 0: Attestation_1
	if dst, err = c.Attestation_1.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Attestation_2
	if dst, err = c.Attestation_2.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}
func (c *AttesterSlashing) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return ssz.ErrSize
	}

	v0 := ssz.ReadOffset(buf[0:4]) // c.Attestation_1
	if v0 < 8 {
		return ssz.ErrInvalidVariableOffset
	}
	if v0 > size {
		return ssz.ErrOffset
	}
	v1 := ssz.ReadOffset(buf[4:8]) // c.Attestation_2
	if v1 > size || v1 < v0 {
		return ssz.ErrOffset
	}
	s0 := buf[v0:v1] // c.Attestation_1
	s1 := buf[v1:]   // c.Attestation_2

	// Field 0: Attestation_1
	c.Attestation_1 = new(IndexedAttestation)
	if err = c.Attestation_1.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Attestation_2
	c.Attestation_2 = new(IndexedAttestation)
	if err = c.Attestation_2.UnmarshalSSZ(s1); err != nil {
		return err
	}
	return err
}
func (c *BeaconBlock) XXSizeSSZ() int {
	size := 84
	if c.Body == nil {
		c.Body = new(BeaconBlockBody)
	}
	size += c.Body.SizeSSZ()
	return size
}
func (c *BeaconBlock) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *BeaconBlock) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 84

	// Field 0: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 1: ProposerIndex
	dst = ssz.MarshalUint64(dst, uint64(c.ProposerIndex))

	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ParentRoot...)

	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.StateRoot...)

	// Field 4: Body
	if c.Body == nil {
		c.Body = new(BeaconBlockBody)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Body.SizeSSZ()

	// Field 4: Body
	if dst, err = c.Body.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}
func (c *BeaconBlock) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 84 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]   // c.Slot
	s1 := buf[8:16]  // c.ProposerIndex
	s2 := buf[16:48] // c.ParentRoot
	s3 := buf[48:80] // c.StateRoot

	v4 := ssz.ReadOffset(buf[80:84]) // c.Body
	if v4 < 84 {
		return ssz.ErrInvalidVariableOffset
	}
	if v4 > size {
		return ssz.ErrOffset
	}
	s4 := buf[v4:] // c.Body

	// Field 0: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s0))

	// Field 1: ProposerIndex
	c.ProposerIndex = prysmaticlabs_eth2_types.ValidatorIndex(ssz.UnmarshallUint64(s1))

	// Field 2: ParentRoot
	c.ParentRoot = append([]byte{}, s2...)

	// Field 3: StateRoot
	c.StateRoot = append([]byte{}, s3...)

	// Field 4: Body
	c.Body = new(BeaconBlockBody)
	if err = c.Body.UnmarshalSSZ(s4); err != nil {
		return err
	}
	return err
}
func (c *BeaconBlockHeader) XXSizeSSZ() int {
	size := 112

	return size
}
func (c *BeaconBlockHeader) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *BeaconBlockHeader) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 1: ProposerIndex
	dst = ssz.MarshalUint64(dst, uint64(c.ProposerIndex))

	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ParentRoot...)

	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.StateRoot...)

	// Field 4: BodyRoot
	if len(c.BodyRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BodyRoot...)

	return dst, err
}
func (c *BeaconBlockHeader) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 112 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]    // c.Slot
	s1 := buf[8:16]   // c.ProposerIndex
	s2 := buf[16:48]  // c.ParentRoot
	s3 := buf[48:80]  // c.StateRoot
	s4 := buf[80:112] // c.BodyRoot

	// Field 0: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s0))

	// Field 1: ProposerIndex
	c.ProposerIndex = prysmaticlabs_eth2_types.ValidatorIndex(ssz.UnmarshallUint64(s1))

	// Field 2: ParentRoot
	c.ParentRoot = append([]byte{}, s2...)

	// Field 3: StateRoot
	c.StateRoot = append([]byte{}, s3...)

	// Field 4: BodyRoot
	c.BodyRoot = append([]byte{}, s4...)
	return err
}
func (c *Checkpoint) XXSizeSSZ() int {
	size := 40

	return size
}
func (c *Checkpoint) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *Checkpoint) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Epoch
	dst = ssz.MarshalUint64(dst, uint64(c.Epoch))

	// Field 1: Root
	if len(c.Root) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Root...)

	return dst, err
}
func (c *Checkpoint) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 40 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]  // c.Epoch
	s1 := buf[8:40] // c.Root

	// Field 0: Epoch
	c.Epoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s0))

	// Field 1: Root
	c.Root = append([]byte{}, s1...)
	return err
}
func (c *Deposit) XXSizeSSZ() int {
	size := 1240

	return size
}
func (c *Deposit) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *Deposit) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Proof
	if len(c.Proof) != 33 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Proof {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(Deposit_Data)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}
func (c *Deposit) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 1240 {
		return ssz.ErrSize
	}

	s0 := buf[0:1056]    // c.Proof
	s1 := buf[1056:1240] // c.Data

	// Field 0: Proof
	{
		var tmp []byte
		for i := 0; i < 33; i++ {
			tmpSlice := s0[i*32 : (1+i)*32]
			tmp = append([]byte{}, tmpSlice...)
			c.Proof = append(c.Proof, tmp)
		}
	}

	// Field 1: Data
	c.Data = new(Deposit_Data)
	if err = c.Data.UnmarshalSSZ(s1); err != nil {
		return err
	}
	return err
}
func (c *Eth1Data) XXSizeSSZ() int {
	size := 72

	return size
}
func (c *Eth1Data) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *Eth1Data) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: DepositRoot
	if len(c.DepositRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.DepositRoot...)

	// Field 1: DepositCount
	dst = ssz.MarshalUint64(dst, c.DepositCount)

	// Field 2: BlockHash
	if len(c.BlockHash) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BlockHash...)

	return dst, err
}
func (c *Eth1Data) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 72 {
		return ssz.ErrSize
	}

	s0 := buf[0:32]  // c.DepositRoot
	s1 := buf[32:40] // c.DepositCount
	s2 := buf[40:72] // c.BlockHash

	// Field 0: DepositRoot
	c.DepositRoot = append([]byte{}, s0...)

	// Field 1: DepositCount
	c.DepositCount = ssz.UnmarshallUint64(s1)

	// Field 2: BlockHash
	c.BlockHash = append([]byte{}, s2...)
	return err
}
func (c *IndexedAttestation) XXSizeSSZ() int {
	size := 228
	size += len(c.AttestingIndices) * 8
	return size
}
func (c *IndexedAttestation) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *IndexedAttestation) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 228

	// Field 0: AttestingIndices
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AttestingIndices) * 8

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(AttestationData)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: AttestingIndices
	if len(c.AttestingIndices) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.AttestingIndices {
		dst = ssz.MarshalUint64(dst, o)
	}
	return dst, err
}
func (c *IndexedAttestation) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 228 {
		return ssz.ErrSize
	}

	s1 := buf[4:132]   // c.Data
	s2 := buf[132:228] // c.Signature

	v0 := ssz.ReadOffset(buf[0:4]) // c.AttestingIndices
	if v0 < 228 {
		return ssz.ErrInvalidVariableOffset
	}
	if v0 > size {
		return ssz.ErrOffset
	}
	s0 := buf[v0:] // c.AttestingIndices

	// Field 0: AttestingIndices
	{
		if len(s0)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.AttestingIndices length is %d, which is not a multiple of 8", len(s0))
		}
		numElem := len(s0) / 8
		if numElem > 2048 {
			return fmt.Errorf("ssz-max exceeded: c.AttestingIndices has %d elements, ssz-max is 2048", numElem)
		}
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := s0[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.AttestingIndices = append(c.AttestingIndices, tmp)
		}
	}

	// Field 1: Data
	c.Data = new(AttestationData)
	if err = c.Data.UnmarshalSSZ(s1); err != nil {
		return err
	}

	// Field 2: Signature
	c.Signature = append([]byte{}, s2...)
	return err
}
func (c *ProposerSlashing) XXSizeSSZ() int {
	size := 416

	return size
}
func (c *ProposerSlashing) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *ProposerSlashing) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Header_1
	if c.Header_1 == nil {
		c.Header_1 = new(SignedBeaconBlockHeader)
	}
	if dst, err = c.Header_1.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Header_2
	if c.Header_2 == nil {
		c.Header_2 = new(SignedBeaconBlockHeader)
	}
	if dst, err = c.Header_2.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}
func (c *ProposerSlashing) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 416 {
		return ssz.ErrSize
	}

	s0 := buf[0:208]   // c.Header_1
	s1 := buf[208:416] // c.Header_2

	// Field 0: Header_1
	c.Header_1 = new(SignedBeaconBlockHeader)
	if err = c.Header_1.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Header_2
	c.Header_2 = new(SignedBeaconBlockHeader)
	if err = c.Header_2.UnmarshalSSZ(s1); err != nil {
		return err
	}
	return err
}
func (c *SignedAggregateAttestationAndProof) XXSizeSSZ() int {
	size := 100
	if c.Message == nil {
		c.Message = new(AggregateAttestationAndProof)
	}
	size += c.Message.SizeSSZ()
	return size
}
func (c *SignedAggregateAttestationAndProof) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *SignedAggregateAttestationAndProof) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 100

	// Field 0: Message
	if c.Message == nil {
		c.Message = new(AggregateAttestationAndProof)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Message.SizeSSZ()

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: Message
	if dst, err = c.Message.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}
func (c *SignedAggregateAttestationAndProof) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 100 {
		return ssz.ErrSize
	}

	s1 := buf[4:100] // c.Signature

	v0 := ssz.ReadOffset(buf[0:4]) // c.Message
	if v0 < 100 {
		return ssz.ErrInvalidVariableOffset
	}
	if v0 > size {
		return ssz.ErrOffset
	}
	s0 := buf[v0:] // c.Message

	// Field 0: Message
	c.Message = new(AggregateAttestationAndProof)
	if err = c.Message.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = append([]byte{}, s1...)
	return err
}
func (c *SignedBeaconBlock) XXSizeSSZ() int {
	size := 100
	if c.Block == nil {
		c.Block = new(BeaconBlock)
	}
	size += c.Block.SizeSSZ()
	return size
}
func (c *SignedBeaconBlock) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *SignedBeaconBlock) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 100

	// Field 0: Block
	if c.Block == nil {
		c.Block = new(BeaconBlock)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Block.SizeSSZ()

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: Block
	if dst, err = c.Block.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}
func (c *SignedBeaconBlock) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 100 {
		return ssz.ErrSize
	}

	s1 := buf[4:100] // c.Signature

	v0 := ssz.ReadOffset(buf[0:4]) // c.Block
	if v0 < 100 {
		return ssz.ErrInvalidVariableOffset
	}
	if v0 > size {
		return ssz.ErrOffset
	}
	s0 := buf[v0:] // c.Block

	// Field 0: Block
	c.Block = new(BeaconBlock)
	if err = c.Block.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = append([]byte{}, s1...)
	return err
}
func (c *SignedBeaconBlockHeader) XXSizeSSZ() int {
	size := 208

	return size
}
func (c *SignedBeaconBlockHeader) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *SignedBeaconBlockHeader) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Header
	if c.Header == nil {
		c.Header = new(BeaconBlockHeader)
	}
	if dst, err = c.Header.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}
func (c *SignedBeaconBlockHeader) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 208 {
		return ssz.ErrSize
	}

	s0 := buf[0:112]   // c.Header
	s1 := buf[112:208] // c.Signature

	// Field 0: Header
	c.Header = new(BeaconBlockHeader)
	if err = c.Header.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = append([]byte{}, s1...)
	return err
}
func (c *SignedVoluntaryExit) XXSizeSSZ() int {
	size := 112

	return size
}
func (c *SignedVoluntaryExit) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *SignedVoluntaryExit) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Exit
	if c.Exit == nil {
		c.Exit = new(VoluntaryExit)
	}
	if dst, err = c.Exit.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}
func (c *SignedVoluntaryExit) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 112 {
		return ssz.ErrSize
	}

	s0 := buf[0:16]   // c.Exit
	s1 := buf[16:112] // c.Signature

	// Field 0: Exit
	c.Exit = new(VoluntaryExit)
	if err = c.Exit.UnmarshalSSZ(s0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = append([]byte{}, s1...)
	return err
}
func (c *Validator) XXSizeSSZ() int {
	size := 121

	return size
}
func (c *Validator) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *Validator) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PublicKey...)

	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.WithdrawalCredentials...)

	// Field 2: EffectiveBalance
	dst = ssz.MarshalUint64(dst, c.EffectiveBalance)

	// Field 3: Slashed
	dst = ssz.MarshalBool(dst, c.Slashed)

	// Field 4: ActivationEligibilityEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.ActivationEligibilityEpoch))

	// Field 5: ActivationEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.ActivationEpoch))

	// Field 6: ExitEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.ExitEpoch))

	// Field 7: WithdrawableEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.WithdrawableEpoch))

	return dst, err
}
func (c *Validator) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 121 {
		return ssz.ErrSize
	}

	s0 := buf[0:48]    // c.PublicKey
	s1 := buf[48:80]   // c.WithdrawalCredentials
	s2 := buf[80:88]   // c.EffectiveBalance
	s3 := buf[88:89]   // c.Slashed
	s4 := buf[89:97]   // c.ActivationEligibilityEpoch
	s5 := buf[97:105]  // c.ActivationEpoch
	s6 := buf[105:113] // c.ExitEpoch
	s7 := buf[113:121] // c.WithdrawableEpoch

	// Field 0: PublicKey
	c.PublicKey = append([]byte{}, s0...)

	// Field 1: WithdrawalCredentials
	c.WithdrawalCredentials = append([]byte{}, s1...)

	// Field 2: EffectiveBalance
	c.EffectiveBalance = ssz.UnmarshallUint64(s2)

	// Field 3: Slashed
	c.Slashed = ssz.UnmarshalBool(s3)

	// Field 4: ActivationEligibilityEpoch
	c.ActivationEligibilityEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s4))

	// Field 5: ActivationEpoch
	c.ActivationEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s5))

	// Field 6: ExitEpoch
	c.ExitEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s6))

	// Field 7: WithdrawableEpoch
	c.WithdrawableEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s7))
	return err
}
func (c *VoluntaryExit) XXSizeSSZ() int {
	size := 16

	return size
}
func (c *VoluntaryExit) XXMarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

func (c *VoluntaryExit) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Epoch
	dst = ssz.MarshalUint64(dst, uint64(c.Epoch))

	// Field 1: ValidatorIndex
	dst = ssz.MarshalUint64(dst, uint64(c.ValidatorIndex))

	return dst, err
}
func (c *VoluntaryExit) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]  // c.Epoch
	s1 := buf[8:16] // c.ValidatorIndex

	// Field 0: Epoch
	c.Epoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s0))

	// Field 1: ValidatorIndex
	c.ValidatorIndex = prysmaticlabs_eth2_types.ValidatorIndex(ssz.UnmarshallUint64(s1))
	return err
}
