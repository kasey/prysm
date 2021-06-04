package ethereum_beacon_p2p_v1

type NoImports struct {
	state         int
	sizeCache     int
	unknownFields int

	GenesisTime                 uint64                                          `protobuf:"varint,1001,opt,name=genesis_time,json=genesisTime,proto3" json:"genesis_time,omitempty"`
	GenesisValidatorsRoot       []byte                                          `protobuf:"bytes,1002,opt,name=genesis_validators_root,json=genesisValidatorsRoot,proto3" json:"genesis_validators_root,omitempty" ssz-size:"32"`
	BlockRoots                  [][]byte                                        `protobuf:"bytes,2002,rep,name=block_roots,json=blockRoots,proto3" json:"block_roots,omitempty" ssz-size:"8192,32"`
	StateRoots                  [][]byte                                        `protobuf:"bytes,2003,rep,name=state_roots,json=stateRoots,proto3" json:"state_roots,omitempty" ssz-size:"8192,32"`
	HistoricalRoots             [][]byte                                        `protobuf:"bytes,2004,rep,name=historical_roots,json=historicalRoots,proto3" json:"historical_roots,omitempty" ssz-max:"16777216" ssz-size:"?,32"`
	Eth1DepositIndex            uint64                                          `protobuf:"varint,3003,opt,name=eth1_deposit_index,json=eth1DepositIndex,proto3" json:"eth1_deposit_index,omitempty"`
	Balances                    []uint64                                        `protobuf:"varint,4002,rep,packed,name=balances,proto3" json:"balances,omitempty" ssz-max:"1099511627776"`
	RandaoMixes                 [][]byte                                        `protobuf:"bytes,5001,rep,name=randao_mixes,json=randaoMixes,proto3" json:"randao_mixes,omitempty" ssz-size:"65536,32"`
	Slashings                   []uint64                                        `protobuf:"varint,6001,rep,packed,name=slashings,proto3" json:"slashings,omitempty" ssz-size:"8192"`
}