package ethereum_beacon_p2p_v1

import (
	"fmt"
	ssz "github.com/ferranbt/fastssz"
	prysmaticlabs_eth2_types "github.com/prysmaticlabs/eth2-types"
	prysmaticlabs_go_bitfield "github.com/prysmaticlabs/go-bitfield"
	prysmaticlabs_prysm_proto_eth_v1alpha1 "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
)

func (c *BeaconBlocksByRangeRequest) XXSizeSSZ() int {
	size := 24

	return size
}
func (c *BeaconBlocksByRangeRequest) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *BeaconBlocksByRangeRequest) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: StartSlot
	dst = ssz.MarshalUint64(dst, uint64(c.StartSlot))

	// Field 1: Count
	dst = ssz.MarshalUint64(dst, c.Count)

	// Field 2: Step
	dst = ssz.MarshalUint64(dst, c.Step)

	return dst, err
}
func (c *BeaconBlocksByRangeRequest) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 24 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]   // c.StartSlot
	s1 := buf[8:16]  // c.Count
	s2 := buf[16:24] // c.Step

	// Field 0: StartSlot
	c.StartSlot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s0))

	// Field 1: Count
	c.Count = ssz.UnmarshallUint64(s1)

	// Field 2: Step
	c.Step = ssz.UnmarshallUint64(s2)
	return err
}
func (c *BeaconBlocksByRangeRequest) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconBlocksByRangeRequest) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: StartSlot
	hh.PutUint64(uint64(c.StartSlot))
	// Field 1: Count
	hh.PutUint64(c.Count)
	// Field 2: Step
	hh.PutUint64(c.Step)
	hh.Merkleize(indx)
	return nil
}
func (c *DepositMessage) XXSizeSSZ() int {
	size := 88

	return size
}
func (c *DepositMessage) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *DepositMessage) XXMarshalSSZTo(dst []byte) ([]byte, error) {
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

	// Field 2: Amount
	dst = ssz.MarshalUint64(dst, c.Amount)

	return dst, err
}
func (c *DepositMessage) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 88 {
		return ssz.ErrSize
	}

	s0 := buf[0:48]  // c.PublicKey
	s1 := buf[48:80] // c.WithdrawalCredentials
	s2 := buf[80:88] // c.Amount

	// Field 0: PublicKey
	c.PublicKey = make([]byte, 0, 48)
	c.PublicKey = append(c.PublicKey, s0...)

	// Field 1: WithdrawalCredentials
	c.WithdrawalCredentials = make([]byte, 0, 32)
	c.WithdrawalCredentials = append(c.WithdrawalCredentials, s1...)

	// Field 2: Amount
	c.Amount = ssz.UnmarshallUint64(s2)
	return err
}
func (c *DepositMessage) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *DepositMessage) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PublicKey)
	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.WithdrawalCredentials)
	// Field 2: Amount
	hh.PutUint64(c.Amount)
	hh.Merkleize(indx)
	return nil
}
func (c *ENRForkID) XXSizeSSZ() int {
	size := 16

	return size
}
func (c *ENRForkID) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *ENRForkID) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: CurrentForkDigest
	if len(c.CurrentForkDigest) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentForkDigest...)

	// Field 1: NextForkVersion
	if len(c.NextForkVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.NextForkVersion...)

	// Field 2: NextForkEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.NextForkEpoch))

	return dst, err
}
func (c *ENRForkID) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	s0 := buf[0:4]  // c.CurrentForkDigest
	s1 := buf[4:8]  // c.NextForkVersion
	s2 := buf[8:16] // c.NextForkEpoch

	// Field 0: CurrentForkDigest
	c.CurrentForkDigest = make([]byte, 0, 4)
	c.CurrentForkDigest = append(c.CurrentForkDigest, s0...)

	// Field 1: NextForkVersion
	c.NextForkVersion = make([]byte, 0, 4)
	c.NextForkVersion = append(c.NextForkVersion, s1...)

	// Field 2: NextForkEpoch
	c.NextForkEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s2))
	return err
}
func (c *ENRForkID) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *ENRForkID) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: CurrentForkDigest
	if len(c.CurrentForkDigest) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentForkDigest)
	// Field 1: NextForkVersion
	if len(c.NextForkVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.NextForkVersion)
	// Field 2: NextForkEpoch
	hh.PutUint64(uint64(c.NextForkEpoch))
	hh.Merkleize(indx)
	return nil
}
func (c *MetaDataV0) XXSizeSSZ() int {
	size := 16

	return size
}
func (c *MetaDataV0) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *MetaDataV0) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: SeqNumber
	dst = ssz.MarshalUint64(dst, c.SeqNumber)

	// Field 1: Attnets
	if len([]byte(c.Attnets)) != 8 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.Attnets)...)

	return dst, err
}
func (c *MetaDataV0) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]  // c.SeqNumber
	s1 := buf[8:16] // c.Attnets

	// Field 0: SeqNumber
	c.SeqNumber = ssz.UnmarshallUint64(s0)

	// Field 1: Attnets
	c.Attnets = make([]byte, 0, 8)
	c.Attnets = append(c.Attnets, prysmaticlabs_go_bitfield.Bitvector64(s1)...)
	return err
}
func (c *MetaDataV0) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *MetaDataV0) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: SeqNumber
	hh.PutUint64(c.SeqNumber)
	// Field 1: Attnets
	if len([]byte(c.Attnets)) != 8 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.Attnets))
	hh.Merkleize(indx)
	return nil
}
func (c *MetaDataV1) XXSizeSSZ() int {
	size := 80

	return size
}
func (c *MetaDataV1) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *MetaDataV1) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: SeqNumber
	dst = ssz.MarshalUint64(dst, c.SeqNumber)

	// Field 1: Attnets
	if len([]byte(c.Attnets)) != 8 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.Attnets)...)

	// Field 2: Syncnets
	if len([]byte(c.Syncnets)) != 64 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.Syncnets)...)

	return dst, err
}
func (c *MetaDataV1) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 80 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]   // c.SeqNumber
	s1 := buf[8:16]  // c.Attnets
	s2 := buf[16:80] // c.Syncnets

	// Field 0: SeqNumber
	c.SeqNumber = ssz.UnmarshallUint64(s0)

	// Field 1: Attnets
	c.Attnets = make([]byte, 0, 8)
	c.Attnets = append(c.Attnets, prysmaticlabs_go_bitfield.Bitvector64(s1)...)

	// Field 2: Syncnets
	c.Syncnets = make([]byte, 0, 64)
	c.Syncnets = append(c.Syncnets, prysmaticlabs_go_bitfield.Bitvector512(s2)...)
	return err
}
func (c *MetaDataV1) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *MetaDataV1) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: SeqNumber
	hh.PutUint64(c.SeqNumber)
	// Field 1: Attnets
	if len([]byte(c.Attnets)) != 8 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.Attnets))
	// Field 2: Syncnets
	if len([]byte(c.Syncnets)) != 64 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.Syncnets))
	hh.Merkleize(indx)
	return nil
}
func (c *Fork) XXSizeSSZ() int {
	size := 16

	return size
}
func (c *Fork) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *Fork) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PreviousVersion
	if len(c.PreviousVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PreviousVersion...)

	// Field 1: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentVersion...)

	// Field 2: Epoch
	dst = ssz.MarshalUint64(dst, uint64(c.Epoch))

	return dst, err
}
func (c *Fork) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	s0 := buf[0:4]  // c.PreviousVersion
	s1 := buf[4:8]  // c.CurrentVersion
	s2 := buf[8:16] // c.Epoch

	// Field 0: PreviousVersion
	c.PreviousVersion = make([]byte, 0, 4)
	c.PreviousVersion = append(c.PreviousVersion, s0...)

	// Field 1: CurrentVersion
	c.CurrentVersion = make([]byte, 0, 4)
	c.CurrentVersion = append(c.CurrentVersion, s1...)

	// Field 2: Epoch
	c.Epoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s2))
	return err
}
func (c *Fork) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Fork) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PreviousVersion
	if len(c.PreviousVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PreviousVersion)
	// Field 1: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentVersion)
	// Field 2: Epoch
	hh.PutUint64(uint64(c.Epoch))
	hh.Merkleize(indx)
	return nil
}
func (c *ForkData) XXSizeSSZ() int {
	size := 36

	return size
}
func (c *ForkData) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *ForkData) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentVersion...)

	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.GenesisValidatorsRoot...)

	return dst, err
}
func (c *ForkData) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 36 {
		return ssz.ErrSize
	}

	s0 := buf[0:4]  // c.CurrentVersion
	s1 := buf[4:36] // c.GenesisValidatorsRoot

	// Field 0: CurrentVersion
	c.CurrentVersion = make([]byte, 0, 4)
	c.CurrentVersion = append(c.CurrentVersion, s0...)

	// Field 1: GenesisValidatorsRoot
	c.GenesisValidatorsRoot = make([]byte, 0, 32)
	c.GenesisValidatorsRoot = append(c.GenesisValidatorsRoot, s1...)
	return err
}
func (c *ForkData) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *ForkData) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentVersion)
	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.GenesisValidatorsRoot)
	hh.Merkleize(indx)
	return nil
}
func (c *HistoricalBatch) XXSizeSSZ() int {
	size := 524288

	return size
}
func (c *HistoricalBatch) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *HistoricalBatch) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: BlockRoots
	if len(c.BlockRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.BlockRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 1: StateRoots
	if len(c.StateRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.StateRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	return dst, err
}
func (c *HistoricalBatch) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 524288 {
		return ssz.ErrSize
	}

	s0 := buf[0:262144]      // c.BlockRoots
	s1 := buf[262144:524288] // c.StateRoots

	// Field 0: BlockRoots
	{
		var tmp []byte
		c.BlockRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s0[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.BlockRoots[i] = tmp
		}
	}

	// Field 1: StateRoots
	{
		var tmp []byte
		c.StateRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s1[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.StateRoots[i] = tmp
		}
	}
	return err
}
func (c *HistoricalBatch) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *HistoricalBatch) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: BlockRoots
	{
		if len(c.BlockRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.BlockRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 1: StateRoots
	{
		if len(c.StateRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.StateRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	hh.Merkleize(indx)
	return nil
}
func (c *Status) XXSizeSSZ() int {
	size := 84

	return size
}
func (c *Status) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *Status) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: ForkDigest
	if len(c.ForkDigest) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ForkDigest...)

	// Field 1: FinalizedRoot
	if len(c.FinalizedRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.FinalizedRoot...)

	// Field 2: FinalizedEpoch
	dst = ssz.MarshalUint64(dst, uint64(c.FinalizedEpoch))

	// Field 3: HeadRoot
	if len(c.HeadRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.HeadRoot...)

	// Field 4: HeadSlot
	dst = ssz.MarshalUint64(dst, uint64(c.HeadSlot))

	return dst, err
}
func (c *Status) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 84 {
		return ssz.ErrSize
	}

	s0 := buf[0:4]   // c.ForkDigest
	s1 := buf[4:36]  // c.FinalizedRoot
	s2 := buf[36:44] // c.FinalizedEpoch
	s3 := buf[44:76] // c.HeadRoot
	s4 := buf[76:84] // c.HeadSlot

	// Field 0: ForkDigest
	c.ForkDigest = make([]byte, 0, 4)
	c.ForkDigest = append(c.ForkDigest, s0...)

	// Field 1: FinalizedRoot
	c.FinalizedRoot = make([]byte, 0, 32)
	c.FinalizedRoot = append(c.FinalizedRoot, s1...)

	// Field 2: FinalizedEpoch
	c.FinalizedEpoch = prysmaticlabs_eth2_types.Epoch(ssz.UnmarshallUint64(s2))

	// Field 3: HeadRoot
	c.HeadRoot = make([]byte, 0, 32)
	c.HeadRoot = append(c.HeadRoot, s3...)

	// Field 4: HeadSlot
	c.HeadSlot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s4))
	return err
}
func (c *Status) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Status) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: ForkDigest
	if len(c.ForkDigest) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ForkDigest)
	// Field 1: FinalizedRoot
	if len(c.FinalizedRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.FinalizedRoot)
	// Field 2: FinalizedEpoch
	hh.PutUint64(uint64(c.FinalizedEpoch))
	// Field 3: HeadRoot
	if len(c.HeadRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.HeadRoot)
	// Field 4: HeadSlot
	hh.PutUint64(uint64(c.HeadSlot))
	hh.Merkleize(indx)
	return nil
}
func (c *BeaconState) XXSizeSSZ() int {
	size := 2687377
	size += len(c.HistoricalRoots) * 32
	size += len(c.Eth1DataVotes) * 72
	size += len(c.Validators) * 121
	size += len(c.Balances) * 8
	size += func() int {
		s := 0
		for _, o := range c.PreviousEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	size += func() int {
		s := 0
		for _, o := range c.CurrentEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	return size
}
func (c *BeaconState) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *BeaconState) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 2687377

	// Field 0: GenesisTime
	dst = ssz.MarshalUint64(dst, c.GenesisTime)

	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.GenesisValidatorsRoot...)

	// Field 2: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 3: Fork
	if c.Fork == nil {
		c.Fork = new(Fork)
	}
	if dst, err = c.Fork.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 4: LatestBlockHeader
	if c.LatestBlockHeader == nil {
		c.LatestBlockHeader = new(prysmaticlabs_prysm_proto_eth_v1alpha1.BeaconBlockHeader)
	}
	if dst, err = c.LatestBlockHeader.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 5: BlockRoots
	if len(c.BlockRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.BlockRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 6: StateRoots
	if len(c.StateRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.StateRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 7: HistoricalRoots
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.HistoricalRoots) * 32

	// Field 8: Eth1Data
	if c.Eth1Data == nil {
		c.Eth1Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
	}
	if dst, err = c.Eth1Data.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 9: Eth1DataVotes
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Eth1DataVotes) * 72

	// Field 10: Eth1DepositIndex
	dst = ssz.MarshalUint64(dst, c.Eth1DepositIndex)

	// Field 11: Validators
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Validators) * 121

	// Field 12: Balances
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Balances) * 8

	// Field 13: RandaoMixes
	if len(c.RandaoMixes) != 65536 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.RandaoMixes {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 14: Slashings
	if len(c.Slashings) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Slashings {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochAttestations
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.PreviousEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 16: CurrentEpochAttestations
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.CurrentEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.JustificationBits)...)

	// Field 18: PreviousJustifiedCheckpoint
	if c.PreviousJustifiedCheckpoint == nil {
		c.PreviousJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.PreviousJustifiedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 19: CurrentJustifiedCheckpoint
	if c.CurrentJustifiedCheckpoint == nil {
		c.CurrentJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.CurrentJustifiedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 20: FinalizedCheckpoint
	if c.FinalizedCheckpoint == nil {
		c.FinalizedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.FinalizedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 7: HistoricalRoots
	if len(c.HistoricalRoots) > 16777216 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.HistoricalRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 9: Eth1DataVotes
	if len(c.Eth1DataVotes) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Eth1DataVotes {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 11: Validators
	if len(c.Validators) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Validators {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 12: Balances
	if len(c.Balances) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Balances {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochAttestations
	if len(c.PreviousEpochAttestations) > 4096 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.PreviousEpochAttestations)
		for _, o := range c.PreviousEpochAttestations {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.PreviousEpochAttestations {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 16: CurrentEpochAttestations
	if len(c.CurrentEpochAttestations) > 4096 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.CurrentEpochAttestations)
		for _, o := range c.CurrentEpochAttestations {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.CurrentEpochAttestations {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}
	return dst, err
}
func (c *BeaconState) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 2687377 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]              // c.GenesisTime
	s1 := buf[8:40]             // c.GenesisValidatorsRoot
	s2 := buf[40:48]            // c.Slot
	s3 := buf[48:64]            // c.Fork
	s4 := buf[64:176]           // c.LatestBlockHeader
	s5 := buf[176:262320]       // c.BlockRoots
	s6 := buf[262320:524464]    // c.StateRoots
	s8 := buf[524468:524540]    // c.Eth1Data
	s10 := buf[524544:524552]   // c.Eth1DepositIndex
	s13 := buf[524560:2621712]  // c.RandaoMixes
	s14 := buf[2621712:2687248] // c.Slashings
	s17 := buf[2687256:2687257] // c.JustificationBits
	s18 := buf[2687257:2687297] // c.PreviousJustifiedCheckpoint
	s19 := buf[2687297:2687337] // c.CurrentJustifiedCheckpoint
	s20 := buf[2687337:2687377] // c.FinalizedCheckpoint

	v7 := ssz.ReadOffset(buf[524464:524468]) // c.HistoricalRoots
	if v7 < 2687377 {
		return ssz.ErrInvalidVariableOffset
	}
	if v7 > size {
		return ssz.ErrOffset
	}
	v9 := ssz.ReadOffset(buf[524540:524544]) // c.Eth1DataVotes
	if v9 > size || v9 < v7 {
		return ssz.ErrOffset
	}
	v11 := ssz.ReadOffset(buf[524552:524556]) // c.Validators
	if v11 > size || v11 < v9 {
		return ssz.ErrOffset
	}
	v12 := ssz.ReadOffset(buf[524556:524560]) // c.Balances
	if v12 > size || v12 < v11 {
		return ssz.ErrOffset
	}
	v15 := ssz.ReadOffset(buf[2687248:2687252]) // c.PreviousEpochAttestations
	if v15 > size || v15 < v12 {
		return ssz.ErrOffset
	}
	v16 := ssz.ReadOffset(buf[2687252:2687256]) // c.CurrentEpochAttestations
	if v16 > size || v16 < v15 {
		return ssz.ErrOffset
	}
	s7 := buf[v7:v9]    // c.HistoricalRoots
	s9 := buf[v9:v11]   // c.Eth1DataVotes
	s11 := buf[v11:v12] // c.Validators
	s12 := buf[v12:v15] // c.Balances
	s15 := buf[v15:v16] // c.PreviousEpochAttestations
	s16 := buf[v16:]    // c.CurrentEpochAttestations

	// Field 0: GenesisTime
	c.GenesisTime = ssz.UnmarshallUint64(s0)

	// Field 1: GenesisValidatorsRoot
	c.GenesisValidatorsRoot = make([]byte, 0, 32)
	c.GenesisValidatorsRoot = append(c.GenesisValidatorsRoot, s1...)

	// Field 2: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s2))

	// Field 3: Fork
	c.Fork = new(Fork)
	if err = c.Fork.XXUnmarshalSSZ(s3); err != nil {
		return err
	}

	// Field 4: LatestBlockHeader
	c.LatestBlockHeader = new(prysmaticlabs_prysm_proto_eth_v1alpha1.BeaconBlockHeader)
	if err = c.LatestBlockHeader.XXUnmarshalSSZ(s4); err != nil {
		return err
	}

	// Field 5: BlockRoots
	{
		var tmp []byte
		c.BlockRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s5[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.BlockRoots[i] = tmp
		}
	}

	// Field 6: StateRoots
	{
		var tmp []byte
		c.StateRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s6[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.StateRoots[i] = tmp
		}
	}

	// Field 7: HistoricalRoots
	{
		if len(s7)%32 != 0 {
			return fmt.Errorf("misaligned bytes: c.HistoricalRoots length is %d, which is not a multiple of 32", len(s7))
		}
		numElem := len(s7) / 32
		if numElem > 16777216 {
			return fmt.Errorf("ssz-max exceeded: c.HistoricalRoots has %d elements, ssz-max is 16777216", numElem)
		}
		c.HistoricalRoots = make([][]byte, numElem)
		for i := 0; i < numElem; i++ {
			var tmp []byte

			tmpSlice := s7[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.HistoricalRoots[i] = tmp
		}
	}

	// Field 8: Eth1Data
	c.Eth1Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
	if err = c.Eth1Data.XXUnmarshalSSZ(s8); err != nil {
		return err
	}

	// Field 9: Eth1DataVotes
	{
		if len(s9)%72 != 0 {
			return fmt.Errorf("misaligned bytes: c.Eth1DataVotes length is %d, which is not a multiple of 72", len(s9))
		}
		numElem := len(s9) / 72
		if numElem > 2048 {
			return fmt.Errorf("ssz-max exceeded: c.Eth1DataVotes has %d elements, ssz-max is 2048", numElem)
		}
		c.Eth1DataVotes = make([]*prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data
			tmp = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
			tmpSlice := s9[i*72 : (1+i)*72]
			if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Eth1DataVotes[i] = tmp
		}
	}

	// Field 10: Eth1DepositIndex
	c.Eth1DepositIndex = ssz.UnmarshallUint64(s10)

	// Field 11: Validators
	{
		if len(s11)%121 != 0 {
			return fmt.Errorf("misaligned bytes: c.Validators length is %d, which is not a multiple of 121", len(s11))
		}
		numElem := len(s11) / 121
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Validators has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Validators = make([]*prysmaticlabs_prysm_proto_eth_v1alpha1.Validator, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *prysmaticlabs_prysm_proto_eth_v1alpha1.Validator
			tmp = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Validator)
			tmpSlice := s11[i*121 : (1+i)*121]
			if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Validators[i] = tmp
		}
	}

	// Field 12: Balances
	{
		if len(s12)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.Balances length is %d, which is not a multiple of 8", len(s12))
		}
		numElem := len(s12) / 8
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Balances has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Balances = make([]uint64, numElem)
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := s12[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Balances[i] = tmp
		}
	}

	// Field 13: RandaoMixes
	{
		var tmp []byte
		c.RandaoMixes = make([][]byte, 65536)
		for i := 0; i < 65536; i++ {
			tmpSlice := s13[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.RandaoMixes[i] = tmp
		}
	}

	// Field 14: Slashings
	{
		var tmp uint64
		c.Slashings = make([]uint64, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s14[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Slashings[i] = tmp
		}
	}

	// Field 15: PreviousEpochAttestations
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(s15) > 3 {
			firstOffset := ssz.ReadOffset(s15[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.PreviousEpochAttestations, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 4096 {
				return fmt.Errorf("ssz-max exceeded: c.PreviousEpochAttestations has %d elements, ssz-max is 4096", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(s15[i*4 : (i+1)*4])
			}
			c.PreviousEpochAttestations = make([]*PendingAttestation, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *PendingAttestation
				tmp = new(PendingAttestation)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = s15[listOffsets[i]:]
				} else {
					tmpSlice = s15[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.PreviousEpochAttestations[i] = tmp
			}
		}
	}

	// Field 16: CurrentEpochAttestations
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(s16) > 3 {
			firstOffset := ssz.ReadOffset(s16[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.CurrentEpochAttestations, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 4096 {
				return fmt.Errorf("ssz-max exceeded: c.CurrentEpochAttestations has %d elements, ssz-max is 4096", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(s16[i*4 : (i+1)*4])
			}
			c.CurrentEpochAttestations = make([]*PendingAttestation, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *PendingAttestation
				tmp = new(PendingAttestation)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = s16[listOffsets[i]:]
				} else {
					tmpSlice = s16[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.CurrentEpochAttestations[i] = tmp
			}
		}
	}

	// Field 17: JustificationBits
	c.JustificationBits = make([]byte, 0, 1)
	c.JustificationBits = append(c.JustificationBits, prysmaticlabs_go_bitfield.Bitvector4(s17)...)

	// Field 18: PreviousJustifiedCheckpoint
	c.PreviousJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.PreviousJustifiedCheckpoint.XXUnmarshalSSZ(s18); err != nil {
		return err
	}

	// Field 19: CurrentJustifiedCheckpoint
	c.CurrentJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.CurrentJustifiedCheckpoint.XXUnmarshalSSZ(s19); err != nil {
		return err
	}

	// Field 20: FinalizedCheckpoint
	c.FinalizedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.FinalizedCheckpoint.XXUnmarshalSSZ(s20); err != nil {
		return err
	}
	return err
}
func (c *BeaconState) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconState) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: GenesisTime
	hh.PutUint64(c.GenesisTime)
	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.GenesisValidatorsRoot)
	// Field 2: Slot
	hh.PutUint64(uint64(c.Slot))
	// Field 3: Fork
	if err := c.Fork.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 4: LatestBlockHeader
	if err := c.LatestBlockHeader.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 5: BlockRoots
	{
		if len(c.BlockRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.BlockRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 6: StateRoots
	{
		if len(c.StateRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.StateRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 7: HistoricalRoots
	{
		if len(c.HistoricalRoots) > 16777216 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.HistoricalRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		numItems := uint64(len(c.HistoricalRoots))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(16777216, numItems, 32))
	}
	// Field 8: Eth1Data
	if err := c.Eth1Data.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 9: Eth1DataVotes
	{
		if len(c.Eth1DataVotes) > 2048 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Eth1DataVotes {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Eth1DataVotes)), 2048)
	}
	// Field 10: Eth1DepositIndex
	hh.PutUint64(c.Eth1DepositIndex)
	// Field 11: Validators
	{
		if len(c.Validators) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Validators {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Validators)), 1099511627776)
	}
	// Field 12: Balances
	{
		if len(c.Balances) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Balances {
			hh.AppendUint64(o)
		}
		hh.FillUpTo32()
		numItems := uint64(len(c.Balances))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
	}
	// Field 13: RandaoMixes
	{
		if len(c.RandaoMixes) != 65536 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.RandaoMixes {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 14: Slashings
	{
		if len(c.Slashings) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.Slashings {
			hh.AppendUint64(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 15: PreviousEpochAttestations
	{
		if len(c.PreviousEpochAttestations) > 4096 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.PreviousEpochAttestations {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.PreviousEpochAttestations)), 4096)
	}
	// Field 16: CurrentEpochAttestations
	{
		if len(c.CurrentEpochAttestations) > 4096 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.CurrentEpochAttestations {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.CurrentEpochAttestations)), 4096)
	}
	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.JustificationBits))
	// Field 18: PreviousJustifiedCheckpoint
	if err := c.PreviousJustifiedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 19: CurrentJustifiedCheckpoint
	if err := c.CurrentJustifiedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 20: FinalizedCheckpoint
	if err := c.FinalizedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}
func (c *BeaconStateAltair) XXSizeSSZ() int {
	size := 2736629
	size += len(c.HistoricalRoots) * 32
	size += len(c.Eth1DataVotes) * 72
	size += len(c.Validators) * 121
	size += len(c.Balances) * 8
	size += len(c.PreviousEpochParticipation) * 1
	size += len(c.CurrentEpochParticipation) * 1
	size += len(c.InactivityScores) * 8
	return size
}
func (c *BeaconStateAltair) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *BeaconStateAltair) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 2736629

	// Field 0: GenesisTime
	dst = ssz.MarshalUint64(dst, c.GenesisTime)

	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.GenesisValidatorsRoot...)

	// Field 2: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 3: Fork
	if c.Fork == nil {
		c.Fork = new(Fork)
	}
	if dst, err = c.Fork.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 4: LatestBlockHeader
	if c.LatestBlockHeader == nil {
		c.LatestBlockHeader = new(prysmaticlabs_prysm_proto_eth_v1alpha1.BeaconBlockHeader)
	}
	if dst, err = c.LatestBlockHeader.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 5: BlockRoots
	if len(c.BlockRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.BlockRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 6: StateRoots
	if len(c.StateRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.StateRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 7: HistoricalRoots
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.HistoricalRoots) * 32

	// Field 8: Eth1Data
	if c.Eth1Data == nil {
		c.Eth1Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
	}
	if dst, err = c.Eth1Data.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 9: Eth1DataVotes
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Eth1DataVotes) * 72

	// Field 10: Eth1DepositIndex
	dst = ssz.MarshalUint64(dst, c.Eth1DepositIndex)

	// Field 11: Validators
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Validators) * 121

	// Field 12: Balances
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Balances) * 8

	// Field 13: RandaoMixes
	if len(c.RandaoMixes) != 65536 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.RandaoMixes {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 14: Slashings
	if len(c.Slashings) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Slashings {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochParticipation
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.PreviousEpochParticipation) * 1

	// Field 16: CurrentEpochParticipation
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.CurrentEpochParticipation) * 1

	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.JustificationBits)...)

	// Field 18: PreviousJustifiedCheckpoint
	if c.PreviousJustifiedCheckpoint == nil {
		c.PreviousJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.PreviousJustifiedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 19: CurrentJustifiedCheckpoint
	if c.CurrentJustifiedCheckpoint == nil {
		c.CurrentJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.CurrentJustifiedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 20: FinalizedCheckpoint
	if c.FinalizedCheckpoint == nil {
		c.FinalizedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	}
	if dst, err = c.FinalizedCheckpoint.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 21: InactivityScores
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.InactivityScores) * 8

	// Field 22: CurrentSyncCommittee
	if c.CurrentSyncCommittee == nil {
		c.CurrentSyncCommittee = new(SyncCommittee)
	}
	if dst, err = c.CurrentSyncCommittee.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 23: NextSyncCommittee
	if c.NextSyncCommittee == nil {
		c.NextSyncCommittee = new(SyncCommittee)
	}
	if dst, err = c.NextSyncCommittee.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 7: HistoricalRoots
	if len(c.HistoricalRoots) > 16777216 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.HistoricalRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 9: Eth1DataVotes
	if len(c.Eth1DataVotes) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Eth1DataVotes {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 11: Validators
	if len(c.Validators) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Validators {
		if dst, err = o.XXMarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 12: Balances
	if len(c.Balances) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Balances {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochParticipation
	if len(c.PreviousEpochParticipation) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	dst = append(dst, c.PreviousEpochParticipation...)

	// Field 16: CurrentEpochParticipation
	if len(c.CurrentEpochParticipation) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	dst = append(dst, c.CurrentEpochParticipation...)

	// Field 21: InactivityScores
	if len(c.InactivityScores) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.InactivityScores {
		dst = ssz.MarshalUint64(dst, o)
	}
	return dst, err
}
func (c *BeaconStateAltair) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 2736629 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]              // c.GenesisTime
	s1 := buf[8:40]             // c.GenesisValidatorsRoot
	s2 := buf[40:48]            // c.Slot
	s3 := buf[48:64]            // c.Fork
	s4 := buf[64:176]           // c.LatestBlockHeader
	s5 := buf[176:262320]       // c.BlockRoots
	s6 := buf[262320:524464]    // c.StateRoots
	s8 := buf[524468:524540]    // c.Eth1Data
	s10 := buf[524544:524552]   // c.Eth1DepositIndex
	s13 := buf[524560:2621712]  // c.RandaoMixes
	s14 := buf[2621712:2687248] // c.Slashings
	s17 := buf[2687256:2687257] // c.JustificationBits
	s18 := buf[2687257:2687297] // c.PreviousJustifiedCheckpoint
	s19 := buf[2687297:2687337] // c.CurrentJustifiedCheckpoint
	s20 := buf[2687337:2687377] // c.FinalizedCheckpoint
	s22 := buf[2687381:2712005] // c.CurrentSyncCommittee
	s23 := buf[2712005:2736629] // c.NextSyncCommittee

	v7 := ssz.ReadOffset(buf[524464:524468]) // c.HistoricalRoots
	if v7 < 2736629 {
		return ssz.ErrInvalidVariableOffset
	}
	if v7 > size {
		return ssz.ErrOffset
	}
	v9 := ssz.ReadOffset(buf[524540:524544]) // c.Eth1DataVotes
	if v9 > size || v9 < v7 {
		return ssz.ErrOffset
	}
	v11 := ssz.ReadOffset(buf[524552:524556]) // c.Validators
	if v11 > size || v11 < v9 {
		return ssz.ErrOffset
	}
	v12 := ssz.ReadOffset(buf[524556:524560]) // c.Balances
	if v12 > size || v12 < v11 {
		return ssz.ErrOffset
	}
	v15 := ssz.ReadOffset(buf[2687248:2687252]) // c.PreviousEpochParticipation
	if v15 > size || v15 < v12 {
		return ssz.ErrOffset
	}
	v16 := ssz.ReadOffset(buf[2687252:2687256]) // c.CurrentEpochParticipation
	if v16 > size || v16 < v15 {
		return ssz.ErrOffset
	}
	v21 := ssz.ReadOffset(buf[2687377:2687381]) // c.InactivityScores
	if v21 > size || v21 < v16 {
		return ssz.ErrOffset
	}
	s7 := buf[v7:v9]    // c.HistoricalRoots
	s9 := buf[v9:v11]   // c.Eth1DataVotes
	s11 := buf[v11:v12] // c.Validators
	s12 := buf[v12:v15] // c.Balances
	s15 := buf[v15:v16] // c.PreviousEpochParticipation
	s16 := buf[v16:v21] // c.CurrentEpochParticipation
	s21 := buf[v21:]    // c.InactivityScores

	// Field 0: GenesisTime
	c.GenesisTime = ssz.UnmarshallUint64(s0)

	// Field 1: GenesisValidatorsRoot
	c.GenesisValidatorsRoot = make([]byte, 0, 32)
	c.GenesisValidatorsRoot = append(c.GenesisValidatorsRoot, s1...)

	// Field 2: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s2))

	// Field 3: Fork
	c.Fork = new(Fork)
	if err = c.Fork.XXUnmarshalSSZ(s3); err != nil {
		return err
	}

	// Field 4: LatestBlockHeader
	c.LatestBlockHeader = new(prysmaticlabs_prysm_proto_eth_v1alpha1.BeaconBlockHeader)
	if err = c.LatestBlockHeader.XXUnmarshalSSZ(s4); err != nil {
		return err
	}

	// Field 5: BlockRoots
	{
		var tmp []byte
		c.BlockRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s5[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.BlockRoots[i] = tmp
		}
	}

	// Field 6: StateRoots
	{
		var tmp []byte
		c.StateRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s6[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.StateRoots[i] = tmp
		}
	}

	// Field 7: HistoricalRoots
	{
		if len(s7)%32 != 0 {
			return fmt.Errorf("misaligned bytes: c.HistoricalRoots length is %d, which is not a multiple of 32", len(s7))
		}
		numElem := len(s7) / 32
		if numElem > 16777216 {
			return fmt.Errorf("ssz-max exceeded: c.HistoricalRoots has %d elements, ssz-max is 16777216", numElem)
		}
		c.HistoricalRoots = make([][]byte, numElem)
		for i := 0; i < numElem; i++ {
			var tmp []byte

			tmpSlice := s7[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.HistoricalRoots[i] = tmp
		}
	}

	// Field 8: Eth1Data
	c.Eth1Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
	if err = c.Eth1Data.XXUnmarshalSSZ(s8); err != nil {
		return err
	}

	// Field 9: Eth1DataVotes
	{
		if len(s9)%72 != 0 {
			return fmt.Errorf("misaligned bytes: c.Eth1DataVotes length is %d, which is not a multiple of 72", len(s9))
		}
		numElem := len(s9) / 72
		if numElem > 2048 {
			return fmt.Errorf("ssz-max exceeded: c.Eth1DataVotes has %d elements, ssz-max is 2048", numElem)
		}
		c.Eth1DataVotes = make([]*prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data
			tmp = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Eth1Data)
			tmpSlice := s9[i*72 : (1+i)*72]
			if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Eth1DataVotes[i] = tmp
		}
	}

	// Field 10: Eth1DepositIndex
	c.Eth1DepositIndex = ssz.UnmarshallUint64(s10)

	// Field 11: Validators
	{
		if len(s11)%121 != 0 {
			return fmt.Errorf("misaligned bytes: c.Validators length is %d, which is not a multiple of 121", len(s11))
		}
		numElem := len(s11) / 121
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Validators has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Validators = make([]*prysmaticlabs_prysm_proto_eth_v1alpha1.Validator, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *prysmaticlabs_prysm_proto_eth_v1alpha1.Validator
			tmp = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Validator)
			tmpSlice := s11[i*121 : (1+i)*121]
			if err = tmp.XXUnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Validators[i] = tmp
		}
	}

	// Field 12: Balances
	{
		if len(s12)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.Balances length is %d, which is not a multiple of 8", len(s12))
		}
		numElem := len(s12) / 8
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Balances has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Balances = make([]uint64, numElem)
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := s12[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Balances[i] = tmp
		}
	}

	// Field 13: RandaoMixes
	{
		var tmp []byte
		c.RandaoMixes = make([][]byte, 65536)
		for i := 0; i < 65536; i++ {
			tmpSlice := s13[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.RandaoMixes[i] = tmp
		}
	}

	// Field 14: Slashings
	{
		var tmp uint64
		c.Slashings = make([]uint64, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := s14[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Slashings[i] = tmp
		}
	}

	// Field 15: PreviousEpochParticipation
	c.PreviousEpochParticipation = append([]byte{}, s15...)

	// Field 16: CurrentEpochParticipation
	c.CurrentEpochParticipation = append([]byte{}, s16...)

	// Field 17: JustificationBits
	c.JustificationBits = make([]byte, 0, 1)
	c.JustificationBits = append(c.JustificationBits, prysmaticlabs_go_bitfield.Bitvector4(s17)...)

	// Field 18: PreviousJustifiedCheckpoint
	c.PreviousJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.PreviousJustifiedCheckpoint.XXUnmarshalSSZ(s18); err != nil {
		return err
	}

	// Field 19: CurrentJustifiedCheckpoint
	c.CurrentJustifiedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.CurrentJustifiedCheckpoint.XXUnmarshalSSZ(s19); err != nil {
		return err
	}

	// Field 20: FinalizedCheckpoint
	c.FinalizedCheckpoint = new(prysmaticlabs_prysm_proto_eth_v1alpha1.Checkpoint)
	if err = c.FinalizedCheckpoint.XXUnmarshalSSZ(s20); err != nil {
		return err
	}

	// Field 21: InactivityScores
	{
		if len(s21)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.InactivityScores length is %d, which is not a multiple of 8", len(s21))
		}
		numElem := len(s21) / 8
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.InactivityScores has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.InactivityScores = make([]uint64, numElem)
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := s21[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.InactivityScores[i] = tmp
		}
	}

	// Field 22: CurrentSyncCommittee
	c.CurrentSyncCommittee = new(SyncCommittee)
	if err = c.CurrentSyncCommittee.XXUnmarshalSSZ(s22); err != nil {
		return err
	}

	// Field 23: NextSyncCommittee
	c.NextSyncCommittee = new(SyncCommittee)
	if err = c.NextSyncCommittee.XXUnmarshalSSZ(s23); err != nil {
		return err
	}
	return err
}
func (c *BeaconStateAltair) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconStateAltair) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: GenesisTime
	hh.PutUint64(c.GenesisTime)
	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.GenesisValidatorsRoot)
	// Field 2: Slot
	hh.PutUint64(uint64(c.Slot))
	// Field 3: Fork
	if err := c.Fork.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 4: LatestBlockHeader
	if err := c.LatestBlockHeader.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 5: BlockRoots
	{
		if len(c.BlockRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.BlockRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 6: StateRoots
	{
		if len(c.StateRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.StateRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 7: HistoricalRoots
	{
		if len(c.HistoricalRoots) > 16777216 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.HistoricalRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		numItems := uint64(len(c.HistoricalRoots))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(16777216, numItems, 32))
	}
	// Field 8: Eth1Data
	if err := c.Eth1Data.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 9: Eth1DataVotes
	{
		if len(c.Eth1DataVotes) > 2048 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Eth1DataVotes {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Eth1DataVotes)), 2048)
	}
	// Field 10: Eth1DepositIndex
	hh.PutUint64(c.Eth1DepositIndex)
	// Field 11: Validators
	{
		if len(c.Validators) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Validators {
			if err := o.XXHashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Validators)), 1099511627776)
	}
	// Field 12: Balances
	{
		if len(c.Balances) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Balances {
			hh.AppendUint64(o)
		}
		hh.FillUpTo32()
		numItems := uint64(len(c.Balances))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
	}
	// Field 13: RandaoMixes
	{
		if len(c.RandaoMixes) != 65536 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.RandaoMixes {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 14: Slashings
	{
		if len(c.Slashings) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.Slashings {
			hh.AppendUint64(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 15: PreviousEpochParticipation
	if len(c.PreviousEpochParticipation) > 1099511627776 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PreviousEpochParticipation)
	// Field 16: CurrentEpochParticipation
	if len(c.CurrentEpochParticipation) > 1099511627776 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentEpochParticipation)
	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.JustificationBits))
	// Field 18: PreviousJustifiedCheckpoint
	if err := c.PreviousJustifiedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 19: CurrentJustifiedCheckpoint
	if err := c.CurrentJustifiedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 20: FinalizedCheckpoint
	if err := c.FinalizedCheckpoint.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 21: InactivityScores
	{
		if len(c.InactivityScores) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.InactivityScores {
			hh.AppendUint64(o)
		}
		hh.FillUpTo32()
		numItems := uint64(len(c.InactivityScores))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
	}
	// Field 22: CurrentSyncCommittee
	if err := c.CurrentSyncCommittee.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 23: NextSyncCommittee
	if err := c.NextSyncCommittee.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}
func (c *PendingAttestation) XXSizeSSZ() int {
	size := 148

	return size
}
func (c *PendingAttestation) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *PendingAttestation) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 148

	// Field 0: AggregationBits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AggregationBits) * 1

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.AttestationData)
	}
	if dst, err = c.Data.XXMarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: InclusionDelay
	dst = ssz.MarshalUint64(dst, uint64(c.InclusionDelay))

	// Field 3: ProposerIndex
	dst = ssz.MarshalUint64(dst, uint64(c.ProposerIndex))

	// Field 0: AggregationBits
	if len(c.AggregationBits) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	dst = append(dst, c.AggregationBits...)
	return dst, err
}
func (c *PendingAttestation) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 148 {
		return ssz.ErrSize
	}

	s1 := buf[4:132]   // c.Data
	s2 := buf[132:140] // c.InclusionDelay
	s3 := buf[140:148] // c.ProposerIndex

	v0 := ssz.ReadOffset(buf[0:4]) // c.AggregationBits
	if v0 < 148 {
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
	c.Data = new(prysmaticlabs_prysm_proto_eth_v1alpha1.AttestationData)
	if err = c.Data.XXUnmarshalSSZ(s1); err != nil {
		return err
	}

	// Field 2: InclusionDelay
	c.InclusionDelay = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s2))

	// Field 3: ProposerIndex
	c.ProposerIndex = prysmaticlabs_eth2_types.ValidatorIndex(ssz.UnmarshallUint64(s3))
	return err
}
func (c *PendingAttestation) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *PendingAttestation) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: AggregationBits
	if len(c.AggregationBits) == 0 {
		return ssz.ErrEmptyBitlist
	}
	hh.PutBitlist(c.AggregationBits, 2048)
	// Field 1: Data
	if err := c.Data.XXHashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: InclusionDelay
	hh.PutUint64(uint64(c.InclusionDelay))
	// Field 3: ProposerIndex
	hh.PutUint64(uint64(c.ProposerIndex))
	hh.Merkleize(indx)
	return nil
}
func (c *SigningData) XXSizeSSZ() int {
	size := 64

	return size
}
func (c *SigningData) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *SigningData) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: ObjectRoot
	if len(c.ObjectRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ObjectRoot...)

	// Field 1: Domain
	if len(c.Domain) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Domain...)

	return dst, err
}
func (c *SigningData) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 64 {
		return ssz.ErrSize
	}

	s0 := buf[0:32]  // c.ObjectRoot
	s1 := buf[32:64] // c.Domain

	// Field 0: ObjectRoot
	c.ObjectRoot = make([]byte, 0, 32)
	c.ObjectRoot = append(c.ObjectRoot, s0...)

	// Field 1: Domain
	c.Domain = make([]byte, 0, 32)
	c.Domain = append(c.Domain, s1...)
	return err
}
func (c *SigningData) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SigningData) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: ObjectRoot
	if len(c.ObjectRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ObjectRoot)
	// Field 1: Domain
	if len(c.Domain) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Domain)
	hh.Merkleize(indx)
	return nil
}
func (c *SyncCommittee) XXSizeSSZ() int {
	size := 24624

	return size
}
func (c *SyncCommittee) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *SyncCommittee) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Pubkeys
	if len(c.Pubkeys) != 512 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Pubkeys {
		if len(o) != 48 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 1: AggregatePubkey
	if len(c.AggregatePubkey) != 48 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.AggregatePubkey...)

	return dst, err
}
func (c *SyncCommittee) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 24624 {
		return ssz.ErrSize
	}

	s0 := buf[0:24576]     // c.Pubkeys
	s1 := buf[24576:24624] // c.AggregatePubkey

	// Field 0: Pubkeys
	{
		var tmp []byte
		c.Pubkeys = make([][]byte, 512)
		for i := 0; i < 512; i++ {
			tmpSlice := s0[i*48 : (1+i)*48]
			tmp = make([]byte, 0, 48)
			tmp = append(tmp, tmpSlice...)
			c.Pubkeys[i] = tmp
		}
	}

	// Field 1: AggregatePubkey
	c.AggregatePubkey = make([]byte, 0, 48)
	c.AggregatePubkey = append(c.AggregatePubkey, s1...)
	return err
}
func (c *SyncCommittee) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SyncCommittee) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Pubkeys
	{
		if len(c.Pubkeys) != 512 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.Pubkeys {
			if len(o) != 48 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 1: AggregatePubkey
	if len(c.AggregatePubkey) != 48 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.AggregatePubkey)
	hh.Merkleize(indx)
	return nil
}
func (c *SyncAggregatorSelectionData) XXSizeSSZ() int {
	size := 16

	return size
}
func (c *SyncAggregatorSelectionData) XXMarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.XXSizeSSZ())
	return c.XXMarshalSSZTo(buf[:0])
}

func (c *SyncAggregatorSelectionData) XXMarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Slot
	dst = ssz.MarshalUint64(dst, uint64(c.Slot))

	// Field 1: SubcommitteeIndex
	dst = ssz.MarshalUint64(dst, c.SubcommitteeIndex)

	return dst, err
}
func (c *SyncAggregatorSelectionData) XXUnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	s0 := buf[0:8]  // c.Slot
	s1 := buf[8:16] // c.SubcommitteeIndex

	// Field 0: Slot
	c.Slot = prysmaticlabs_eth2_types.Slot(ssz.UnmarshallUint64(s0))

	// Field 1: SubcommitteeIndex
	c.SubcommitteeIndex = ssz.UnmarshallUint64(s1)
	return err
}
func (c *SyncAggregatorSelectionData) XXHashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.XXHashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SyncAggregatorSelectionData) XXHashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Slot
	hh.PutUint64(uint64(c.Slot))
	// Field 1: SubcommitteeIndex
	hh.PutUint64(c.SubcommitteeIndex)
	hh.Merkleize(indx)
	return nil
}
