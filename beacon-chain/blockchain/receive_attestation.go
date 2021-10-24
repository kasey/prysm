package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/async/event"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed"
	statefeed "github.com/prysmaticlabs/prysm/beacon-chain/core/feed/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/attestation"
	"github.com/prysmaticlabs/prysm/time/slots"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"time"
)

// AttestationStateFetcher allows for retrieving a beacon state corresponding to the block
// root of an attestation's target checkpoint.
type AttestationStateFetcher interface {
	AttestationTargetState(ctx context.Context, target *ethpb.Checkpoint) (state.BeaconState, error)
}

// AttestationReceiver interface defines the methods of chain service receive and processing new attestations.
type AttestationReceiver interface {
	AttestationStateFetcher
	ReceiveAttestationNoPubsub(ctx context.Context, att *ethpb.Attestation) error
	VerifyLmdFfgConsistency(ctx context.Context, att *ethpb.Attestation) error
	VerifyFinalizedConsistency(ctx context.Context, root []byte) error
}

// ReceiveAttestationNoPubsub implements `on_attestation` from the beacon chain spec.
// The operations consist of:
//  1. Validate attestation, update validator's latest vote
//  2. Apply fork choice to the processed attestation
//  3. Save latest head info
//
// on_attestation is called whenever an attestation is received, verifies the attestation is valid and saves
// it to the DB. As a stateless function, this does not hold nor delay attestation based on the spec descriptions.
// The delay is handled by the caller in `processAttestations`.
//
// Spec pseudocode definition:
//   def on_attestation(store: Store, attestation: Attestation) -> None:
//    """
//    Run ``on_attestation`` upon receiving a new ``attestation`` from either within a block or directly on the wire.
//
//    An ``attestation`` that is asserted as invalid may be valid at a later time,
//    consider scheduling it for later processing in such case.
//    """
//    validate_on_attestation(store, attestation)
//    store_target_checkpoint_state(store, attestation.data.target)
//
//    # Get state at the `target` to fully validate attestation
//    target_state = store.checkpoint_states[attestation.data.target]
//    indexed_attestation = get_indexed_attestation(target_state, attestation)
//    assert is_valid_indexed_attestation(target_state, indexed_attestation)
//
//    # Update latest messages for attesting indices
//    update_latest_messages(store, indexed_attestation.attesting_indices, attestation)
func (s *Service) ReceiveAttestationNoPubsub(ctx context.Context, att *ethpb.Attestation) error {
	ctx, span := trace.StartSpan(ctx, "beacon-chain.blockchain.ReceiveAttestationNoPubsub")
	defer span.End()

	if err := helpers.ValidateNilAttestation(att); err != nil {
		return errors.Wrap(err, "could not process attestation")
	}
	if err := helpers.ValidateSlotTargetEpoch(att.Data); err != nil {
		return errors.Wrap(err, "could not process attestation")
	}
	tgt := ethpb.CopyCheckpoint(att.Data.Target)

	// Note that target root check is ignored here because it was performed in sync's validation pipeline:
	// validate_aggregate_proof.go and validate_beacon_attestation.go
	// If missing target root were to fail in this method, it would have just failed in `getAttPreState`.

	// Retrieve attestation's data beacon block pre state. Advance pre state to latest epoch if necessary and
	// save it to the cache.
	baseState, err := s.getAttPreState(ctx, tgt)
	if err != nil {
		return errors.Wrap(err, "could not process attestation")
	}

	genesisTime := baseState.GenesisTime()

	// Verify attestation target is from current epoch or previous epoch.
	if err := s.verifyAttTargetEpoch(ctx, genesisTime, uint64(time.Now().Unix()), tgt); err != nil {
		return errors.Wrap(err, "could not process attestation")
	}

	// Verify attestation beacon block is known and not from the future.
	if err := s.verifyBeaconBlock(ctx, att.Data); err != nil {
		err = errors.Wrap(err, "could not verify attestation beacon block")
		return errors.Wrap(err, "could not process attestation")
	}

	// Note that LMG GHOST and FFG consistency check is ignored because it was performed in sync's validation pipeline:
	// validate_aggregate_proof.go and validate_beacon_attestation.go

	// Verify attestations can only affect the fork choice of subsequent slots.
	if err := slots.VerifyTime(genesisTime, att.Data.Slot+1, params.BeaconNetworkConfig().MaximumGossipClockDisparity); err != nil {
		return errors.Wrap(err, "could not process attestation")
	}

	// Use the target state to verify attesting indices are valid.
	committee, err := helpers.BeaconCommitteeFromState(ctx, baseState, att.Data.Slot, att.Data.CommitteeIndex)
	if err != nil {
		return errors.Wrap(err, "could not process attestation")
	}
	indexedAtt, err := attestation.ConvertToIndexed(ctx, att, committee)
	if err != nil {
		return errors.Wrap(err, "could not process attestation")
	}
	if err := attestation.IsValidAttestationIndices(ctx, indexedAtt); err != nil {
		return errors.Wrap(err, "could not process attestation")
	}

	// Note that signature verification is ignored here because it was performed in sync's validation pipeline:
	// validate_aggregate_proof.go and validate_beacon_attestation.go
	// We assume trusted attestation in this function has verified signature.

	// Update forkchoice store with the new attestation for updating weight.
	s.cfg.ForkChoiceStore.ProcessAttestation(ctx, indexedAtt.AttestingIndices, bytesutil.ToBytes32(att.Data.BeaconBlockRoot), att.Data.Target.Epoch)

	return nil
}

// AttestationTargetState returns the pre state of attestation.
func (s *Service) AttestationTargetState(ctx context.Context, target *ethpb.Checkpoint) (state.BeaconState, error) {
	ss, err := slots.EpochStart(target.Epoch)
	if err != nil {
		return nil, err
	}
	if err := slots.ValidateClock(ss, uint64(s.genesisTime.Unix())); err != nil {
		return nil, err
	}
	return s.getAttPreState(ctx, target)
}

// VerifyLmdFfgConsistency verifies that attestation's LMD and FFG votes are consistency to each other.
func (s *Service) VerifyLmdFfgConsistency(ctx context.Context, a *ethpb.Attestation) error {
	targetSlot, err := slots.EpochStart(a.Data.Target.Epoch)
	if err != nil {
		return err
	}
	r, err := s.ancestor(ctx, a.Data.BeaconBlockRoot, targetSlot)
	if err != nil {
		return err
	}
	if !bytes.Equal(a.Data.Target.Root, r) {
		return errors.New("FFG and LMD votes are not consistent")
	}
	return nil
}

// VerifyFinalizedConsistency verifies input root is consistent with finalized store.
// When the input root is not be consistent with finalized store then we know it is not
// on the finalized check point that leads to current canonical chain and should be rejected accordingly.
func (s *Service) VerifyFinalizedConsistency(ctx context.Context, root []byte) error {
	// A canonical root implies the root to has an ancestor that aligns with finalized check point.
	// In this case, we could exit early to save on additional computation.
	if s.cfg.ForkChoiceStore.IsCanonical(bytesutil.ToBytes32(root)) {
		return nil
	}

	f := s.FinalizedCheckpt()
	ss, err := slots.EpochStart(f.Epoch)
	if err != nil {
		return err
	}
	r, err := s.ancestor(ctx, root, ss)
	if err != nil {
		return err
	}
	if !bytes.Equal(f.Root, r) {
		return errors.New("Root and finalized store are not consistent")
	}

	return nil
}

type attestationProcessor interface {
	processAttestations(ctx context.Context, genesisTime time.Time)
}

// This function invokes an attestationProcessor once per slot.
// It blocks until genesis time is initialized and then runs in a non-terminating loop.
// It should run in its own goroutine.
func processAttestationsLoop(ctx context.Context, stateFeed *event.Feed, ap attestationProcessor) {
	genesisTime, err := blockUntilGenesisTimeInitialized(ctx, stateFeed)
	if err != nil {
		log.Error(err)
		return
	}
	st := slots.NewSlotTicker(genesisTime, params.BeaconConfig().SecondsPerSlot)
	for {
		select {
		case <-st.C():
			// Continue when there's no fork choice attestation, there's nothing to process and update head.
			// This covers the condition when the node is still initial syncing to the head of the chain.
			ap.processAttestations(ctx, genesisTime)
		case <-ctx.Done():
			log.Warn(errors.Wrap(ctx.Err(), "processAttestationsLoop received Done signal"))
			return
		}
	}
}

func blockUntilGenesisTimeInitialized(ctx context.Context, stateFeed *event.Feed) (time.Time, error) {
	var genesisTime time.Time
	sleeper := time.NewTicker(1 * time.Second)
	stateChannel := make(chan *feed.Event, 1)
	stateSub := stateFeed.Subscribe(stateChannel)
	defer sleeper.Stop()
	defer stateSub.Unsubscribe()
	defer close(stateChannel)
	for {
		select {
		case ev := <-stateChannel:
			if ev.Type == statefeed.Initialized {
				evData := ev.Data.(statefeed.InitializedData)
				if evData.StartTime.IsZero() {
					log.Error("blockUntilGenesisTimeInitialized received InitializedData with a zero-valued StartTime")
					continue // can't proceed into the attestation routine until valid genesis time is received
				} else {
					genesisTime = evData.StartTime // breaks loop invariant so we can proceed
					log.Warn("Genesis time received, now available to process attestations")
					return genesisTime, nil
				}
			}
		case <-sleeper.C:
			log.Warn("ProcessAttestations routine waiting for genesis time (1 second check interval)")
		case <-ctx.Done():
			return genesisTime, errors.Wrap(ctx.Err(), "blockUntilGenesisTimeInitialized received Done signal")
		}
	}
}

// This processes fork choice attestations from the pool to account for validator votes and fork choice.
func (s *Service) processAttestations(ctx context.Context, genesisTime time.Time) {
	if s.cfg.AttPool.ForkchoiceAttestationCount() == 0 {
		return
	}
	atts := s.cfg.AttPool.ForkchoiceAttestations()
	for _, a := range atts {
		// Based on the spec, don't process the attestation until the subsequent slot.
		// This delays consideration in the fork choice until their slot is in the past.
		// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/fork-choice.md#validate_on_attestation
		nextSlot := a.Data.Slot + 1
		if err := slots.VerifyTime(uint64(genesisTime.Unix()), nextSlot, params.BeaconNetworkConfig().MaximumGossipClockDisparity); err != nil {
			continue
		}

		hasState := s.cfg.BeaconDB.HasStateSummary(ctx, bytesutil.ToBytes32(a.Data.BeaconBlockRoot))
		blockRoot := bytesutil.ToBytes32(a.Data.BeaconBlockRoot)
		hasBlock := s.cfg.ForkChoiceStore.HasNode(blockRoot) || s.cfg.BeaconDB.HasBlock(ctx, blockRoot)
		if !(hasState && hasBlock) {
			continue
		}

		if err := s.cfg.AttPool.DeleteForkchoiceAttestation(a); err != nil {
			log.WithError(err).Error("Could not delete fork choice attestation in pool")
		}

		if !helpers.VerifyCheckpointEpoch(a.Data.Target, genesisTime) {
			continue
		}

		if err := s.ReceiveAttestationNoPubsub(ctx, a); err != nil {
			log.WithFields(logrus.Fields{
				"slot":             a.Data.Slot,
				"committeeIndex":   a.Data.CommitteeIndex,
				"beaconBlockRoot":  fmt.Sprintf("%#x", bytesutil.Trunc(a.Data.BeaconBlockRoot)),
				"targetRoot":       fmt.Sprintf("%#x", bytesutil.Trunc(a.Data.Target.Root)),
				"aggregationCount": a.AggregationBits.Count(),
			}).WithError(err).Warn("Could not process attestation for fork choice")
		}
	}
	if err := s.updateHead(ctx, s.getJustifiedBalances()); err != nil {
		log.Warnf("Resolving fork due to new attestation: %v", err)
	}
}
