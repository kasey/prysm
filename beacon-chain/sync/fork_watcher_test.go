package sync

import (
	"context"
	"testing"
	"time"

	"github.com/prysmaticlabs/prysm/v5/async/abool"
	mockChain "github.com/prysmaticlabs/prysm/v5/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/p2p"
	p2ptest "github.com/prysmaticlabs/prysm/v5/beacon-chain/p2p/testing"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/startup"
	mockSync "github.com/prysmaticlabs/prysm/v5/beacon-chain/sync/initial-sync/testing"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/network/forks"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
	"github.com/prysmaticlabs/prysm/v5/testing/assert"
)

func TestService_CheckForNextEpochFork(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	tests := []struct {
		name         string
		svcCreator   func(t *testing.T) *Service
		currEpoch    primitives.Epoch
		wantErr      bool
		postSvcCheck func(t *testing.T, s *Service)
	}{
		{
			name: "no fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				gt := time.Now().Add(time.Duration(-params.BeaconConfig().SlotsPerEpoch.Mul(uint64(params.BeaconConfig().SlotsPerEpoch))) * time.Second)
				vr := [32]byte{'A'}
				chainService := &mockChain.ChainService{
					Genesis:        gt,
					ValidatorsRoot: vr,
				}
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(gt, vr),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 10,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {

			},
		},
		{
			name: "altair fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				gt := time.Now().Add(-4 * oneEpoch())
				vr := [32]byte{'A'}
				chainService := &mockChain.ChainService{
					Genesis:        gt,
					ValidatorsRoot: vr,
				}
				bCfg := params.BeaconConfig().Copy()
				bCfg.AltairForkEpoch = 5
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(gt, vr),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(5, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
				rpcMap := make(map[string]bool)
				for _, p := range s.cfg.p2p.Host().Mux().Protocols() {
					rpcMap[string(p)] = true
				}
				assert.Equal(t, true, rpcMap[p2p.RPCBlocksByRangeTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
				assert.Equal(t, true, rpcMap[p2p.RPCBlocksByRootTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
				assert.Equal(t, true, rpcMap[p2p.RPCMetaDataTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
			},
		},
		{
			name: "bellatrix fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				chainService := &mockChain.ChainService{
					Genesis:        time.Now().Add(-4 * oneEpoch()),
					ValidatorsRoot: [32]byte{'A'},
				}
				bCfg := params.BeaconConfig().Copy()
				bCfg.AltairForkEpoch = 3
				bCfg.BellatrixForkEpoch = 5
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(chainService.Genesis, chainService.ValidatorsRoot),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(5, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
				rpcMap := make(map[string]bool)
				for _, p := range s.cfg.p2p.Host().Mux().Protocols() {
					rpcMap[string(p)] = true
				}
			},
		},
		{
			name: "deneb fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				gt := time.Now().Add(-4 * oneEpoch())
				vr := [32]byte{'A'}
				chainService := &mockChain.ChainService{
					Genesis:        gt,
					ValidatorsRoot: vr,
				}
				bCfg := params.BeaconConfig().Copy()
				bCfg.DenebForkEpoch = 5
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(gt, vr),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(5, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
				rpcMap := make(map[string]bool)
				for _, p := range s.cfg.p2p.Host().Mux().Protocols() {
					rpcMap[string(p)] = true
				}
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRangeTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRootTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
			},
		},
		{
			name: "electra fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				gt := time.Now().Add(-4 * oneEpoch())
				vr := [32]byte{'A'}
				chainService := &mockChain.ChainService{
					Genesis:        gt,
					ValidatorsRoot: vr,
				}
				bCfg := params.BeaconConfig().Copy()
				bCfg.ElectraForkEpoch = 5
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(gt, vr),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(5, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
				rpcMap := make(map[string]bool)
				for _, p := range s.cfg.p2p.Host().Mux().Protocols() {
					rpcMap[string(p)] = true
				}
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRangeTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRootTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
			},
		},
		{
			name: "fulu fork in the next epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				gt := time.Now().Add(-4 * oneEpoch())
				vr := [32]byte{'A'}
				chainService := &mockChain.ChainService{
					Genesis:        gt,
					ValidatorsRoot: vr,
				}
				bCfg := params.BeaconConfig().Copy()
				bCfg.FuluForkEpoch = 5
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       startup.NewClock(gt, vr),
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(5, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
				rpcMap := make(map[string]bool)
				for _, p := range s.cfg.p2p.Host().Mux().Protocols() {
					rpcMap[string(p)] = true
				}
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRangeTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
				assert.Equal(t, true, rpcMap[p2p.RPCBlobSidecarsByRootTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()], "topic doesn't exist")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.svcCreator(t)
			if err := s.registerForUpcomingFork(tt.currEpoch); (err != nil) != tt.wantErr {
				t.Errorf("registerForUpcomingFork() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.postSvcCheck(t, s)
		})
	}
}

func TestService_CheckForPreviousEpochFork(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	tests := []struct {
		name         string
		svcCreator   func(t *testing.T) *Service
		currEpoch    primitives.Epoch
		wantErr      bool
		postSvcCheck func(t *testing.T, s *Service)
	}{
		{
			name: "no fork in the previous epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				chainService := &mockChain.ChainService{
					Genesis:        time.Now().Add(-oneEpoch()),
					ValidatorsRoot: [32]byte{'A'},
				}
				clock := startup.NewClock(chainService.Genesis, chainService.ValidatorsRoot)
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       clock,
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				err := r.registerRPCHandlers()
				assert.NoError(t, err)
				return r
			},
			currEpoch: 10,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				ptcls := s.cfg.p2p.Host().Mux().Protocols()
				pMap := make(map[string]bool)
				for _, p := range ptcls {
					pMap[string(p)] = true
				}
				assert.Equal(t, true, pMap[p2p.RPCGoodByeTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCStatusTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCPingTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCMetaDataTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCBlocksByRangeTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCBlocksByRootTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
			},
		},
		{
			name: "altair fork in the previous epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				chainService := &mockChain.ChainService{
					Genesis:        time.Now().Add(-4 * oneEpoch()),
					ValidatorsRoot: [32]byte{'A'},
				}
				clock := startup.NewClock(chainService.Genesis, chainService.ValidatorsRoot)
				bCfg := params.BeaconConfig().Copy()
				bCfg.AltairForkEpoch = 3
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       clock,
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				prevGenesis := chainService.Genesis
				// To allow registration of v1 handlers
				chainService.Genesis = time.Now().Add(-1 * oneEpoch())
				err := r.registerRPCHandlers()
				assert.NoError(t, err)

				chainService.Genesis = prevGenesis
				previous, err := r.rpcHandlerByTopicFromFork(version.Phase0)
				assert.NoError(t, err)

				next, err := r.rpcHandlerByTopicFromFork(version.Altair)
				assert.NoError(t, err)

				handlerByTopic := addedRPCHandlerByTopic(previous, next)

				for topic, handler := range handlerByTopic {
					r.registerRPC(topic, handler)
				}

				genRoot := r.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(0, genRoot[:])
				assert.NoError(t, err)
				r.registerSubscribers(0, digest)
				assert.Equal(t, true, r.subHandler.digestExists(digest))

				digest, err = forks.ForkDigestFromEpoch(3, genRoot[:])
				assert.NoError(t, err)
				r.registerSubscribers(3, digest)
				assert.Equal(t, true, r.subHandler.digestExists(digest))

				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(0, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, false, s.subHandler.digestExists(digest))
				digest, err = forks.ForkDigestFromEpoch(3, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))

				ptcls := s.cfg.p2p.Host().Mux().Protocols()
				pMap := make(map[string]bool)
				for _, p := range ptcls {
					pMap[string(p)] = true
				}
				assert.Equal(t, true, pMap[p2p.RPCGoodByeTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCStatusTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCPingTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCMetaDataTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCBlocksByRangeTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, true, pMap[p2p.RPCBlocksByRootTopicV2+s.cfg.p2p.Encoding().ProtocolSuffix()])

				assert.Equal(t, false, pMap[p2p.RPCMetaDataTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, false, pMap[p2p.RPCBlocksByRangeTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
				assert.Equal(t, false, pMap[p2p.RPCBlocksByRootTopicV1+s.cfg.p2p.Encoding().ProtocolSuffix()])
			},
		},
		{
			name: "bellatrix fork in the previous epoch",
			svcCreator: func(t *testing.T) *Service {
				peer2peer := p2ptest.NewTestP2P(t)
				chainService := &mockChain.ChainService{
					Genesis:        time.Now().Add(-4 * oneEpoch()),
					ValidatorsRoot: [32]byte{'A'},
				}
				clock := startup.NewClock(chainService.Genesis, chainService.ValidatorsRoot)
				bCfg := params.BeaconConfig().Copy()
				bCfg.AltairForkEpoch = 1
				bCfg.BellatrixForkEpoch = 3
				params.OverrideBeaconConfig(bCfg)
				params.BeaconConfig().InitializeForkSchedule()
				ctx, cancel := context.WithCancel(context.Background())
				r := &Service{
					ctx:    ctx,
					cancel: cancel,
					cfg: &config{
						p2p:         peer2peer,
						chain:       chainService,
						clock:       clock,
						initialSync: &mockSync.Sync{IsSyncing: false},
					},
					chainStarted: abool.New(),
					subHandler:   newSubTopicHandler(),
				}
				genRoot := r.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(1, genRoot[:])
				assert.NoError(t, err)
				r.registerSubscribers(1, digest)
				assert.Equal(t, true, r.subHandler.digestExists(digest))

				digest, err = forks.ForkDigestFromEpoch(3, genRoot[:])
				assert.NoError(t, err)
				r.registerSubscribers(3, digest)
				assert.Equal(t, true, r.subHandler.digestExists(digest))

				return r
			},
			currEpoch: 4,
			wantErr:   false,
			postSvcCheck: func(t *testing.T, s *Service) {
				genRoot := s.cfg.clock.GenesisValidatorsRoot()
				digest, err := forks.ForkDigestFromEpoch(1, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, false, s.subHandler.digestExists(digest))
				digest, err = forks.ForkDigestFromEpoch(3, genRoot[:])
				assert.NoError(t, err)
				assert.Equal(t, true, s.subHandler.digestExists(digest))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.svcCreator(t)
			if err := s.deregisterFromPastFork(tt.currEpoch); (err != nil) != tt.wantErr {
				t.Errorf("registerForUpcomingFork() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.postSvcCheck(t, s)
		})
	}
}

// oneEpoch returns the duration of one epoch.
func oneEpoch() time.Duration {
	return time.Duration(params.BeaconConfig().SlotsPerEpoch.Mul(params.BeaconConfig().SecondsPerSlot)) * time.Second
}
