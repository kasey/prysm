package validator

import (
	"context"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/feed"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/feed/operation"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/rpc/core"
	"github.com/prysmaticlabs/prysm/v5/config/features"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls"
	"github.com/prysmaticlabs/prysm/v5/monitoring/tracing/trace"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/time/slots"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAttestationData requests that the beacon node produce an attestation data object,
// which the validator acting as an attester will then sign.
func (vs *Server) GetAttestationData(ctx context.Context, req *ethpb.AttestationDataRequest) (*ethpb.AttestationData, error) {
	ctx, span := trace.StartSpan(ctx, "AttesterServer.RequestAttestation")
	defer span.End()
	span.SetAttributes(
		trace.Int64Attribute("slot", int64(req.Slot)),
		trace.Int64Attribute("committeeIndex", int64(req.CommitteeIndex)),
	)

	if vs.SyncChecker.Syncing() {
		return nil, status.Errorf(codes.Unavailable, "Syncing to latest head, not ready to respond")
	}
	res, err := vs.CoreService.GetAttestationData(ctx, req)
	if err != nil {
		return nil, status.Errorf(core.ErrorReasonToGRPC(err.Reason), "Could not get attestation data: %v", err.Err)
	}
	return res, nil
}

// ProposeAttestation is a function called by an attester to vote
// on a block via an attestation object as defined in the Ethereum specification.
func (vs *Server) ProposeAttestation(ctx context.Context, att *ethpb.Attestation) (*ethpb.AttestResponse, error) {
	ctx, span := trace.StartSpan(ctx, "AttesterServer.ProposeAttestation")
	defer span.End()

	resp, err := vs.proposeAtt(ctx, att, att.GetData().CommitteeIndex)
	if err != nil {
		return nil, err
	}

	if features.Get().EnableExperimentalAttestationPool {
		if err = vs.AttestationCache.Add(att); err != nil {
			log.WithError(err).Error("Could not save attestation")
		}
	} else {
		go func() {
			attCopy := att.Copy()
			if err := vs.AttPool.SaveUnaggregatedAttestation(attCopy); err != nil {
				log.WithError(err).Error("Could not save unaggregated attestation")
				return
			}
		}()
	}

	return resp, nil
}

// ProposeAttestationElectra is a function called by an attester to vote
// on a block via an attestation object as defined in the Ethereum specification.
func (vs *Server) ProposeAttestationElectra(ctx context.Context, att *ethpb.AttestationElectra) (*ethpb.AttestResponse, error) {
	ctx, span := trace.StartSpan(ctx, "AttesterServer.ProposeAttestationElectra")
	defer span.End()

	committeeIndex, err := att.GetCommitteeIndex()
	if err != nil {
		return nil, err
	}

	resp, err := vs.proposeAtt(ctx, att, committeeIndex)
	if err != nil {
		return nil, err
	}

	if features.Get().EnableExperimentalAttestationPool {
		if err = vs.AttestationCache.Add(att); err != nil {
			log.WithError(err).Error("Could not save attestation")
		}
	} else {
		go func() {
			ctx = trace.NewContext(context.Background(), trace.FromContext(ctx))
			attCopy := att.Copy()
			if err := vs.AttPool.SaveUnaggregatedAttestation(attCopy); err != nil {
				log.WithError(err).Error("Could not save unaggregated attestation")
				return
			}
		}()
	}

	return resp, nil
}

// SubscribeCommitteeSubnets subscribes to the committee ID subnet given subscribe request.
func (vs *Server) SubscribeCommitteeSubnets(ctx context.Context, req *ethpb.CommitteeSubnetsSubscribeRequest) (*emptypb.Empty, error) {
	ctx, span := trace.StartSpan(ctx, "AttesterServer.SubscribeCommitteeSubnets")
	defer span.End()

	if len(req.Slots) != len(req.CommitteeIds) || len(req.CommitteeIds) != len(req.IsAggregator) {
		return nil, status.Error(codes.InvalidArgument, "request fields are not the same length")
	}
	if len(req.Slots) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no attester slots provided")
	}

	fetchValsLen := func(slot primitives.Slot) (uint64, error) {
		wantedEpoch := slots.ToEpoch(slot)
		vals, err := vs.HeadFetcher.HeadValidatorsIndices(ctx, wantedEpoch)
		if err != nil {
			return 0, err
		}
		return uint64(len(vals)), nil
	}

	// Request the head validator indices of epoch represented by the first requested
	// slot.
	currValsLen, err := fetchValsLen(req.Slots[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not retrieve head validator length: %v", err)
	}
	currEpoch := slots.ToEpoch(req.Slots[0])

	for i := 0; i < len(req.Slots); i++ {
		// If epoch has changed, re-request active validators length
		if currEpoch != slots.ToEpoch(req.Slots[i]) {
			currValsLen, err = fetchValsLen(req.Slots[i])
			if err != nil {
				return nil, status.Errorf(codes.Internal, "Could not retrieve head validator length: %v", err)
			}
			currEpoch = slots.ToEpoch(req.Slots[i])
		}
		subnet := helpers.ComputeSubnetFromCommitteeAndSlot(currValsLen, req.CommitteeIds[i], req.Slots[i])
		cache.SubnetIDs.AddAttesterSubnetID(req.Slots[i], subnet)
		if req.IsAggregator[i] {
			cache.SubnetIDs.AddAggregatorSubnetID(req.Slots[i], subnet)
		}
	}

	return &emptypb.Empty{}, nil
}

func (vs *Server) proposeAtt(ctx context.Context, att ethpb.Att, committee primitives.CommitteeIndex) (*ethpb.AttestResponse, error) {
	if _, err := bls.SignatureFromBytes(att.GetSignature()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect attestation signature")
	}

	root, err := att.GetData().HashTreeRoot()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not tree hash attestation: %v", err)
	}

	// Broadcast the unaggregated attestation on a feed to notify other services in the beacon node
	// of a received unaggregated attestation.
	vs.OperationNotifier.OperationFeed().Send(&feed.Event{
		Type: operation.UnaggregatedAttReceived,
		Data: &operation.UnAggregatedAttReceivedData{
			Attestation: att,
		},
	})

	// Determine subnet to broadcast attestation to
	wantedEpoch := slots.ToEpoch(att.GetData().Slot)
	vals, err := vs.HeadFetcher.HeadValidatorsIndices(ctx, wantedEpoch)
	if err != nil {
		return nil, err
	}
	subnet := helpers.ComputeSubnetFromCommitteeAndSlot(uint64(len(vals)), committee, att.GetData().Slot)

	// Broadcast the new attestation to the network.
	if err := vs.P2P.BroadcastAttestation(ctx, subnet, att); err != nil {
		return nil, status.Errorf(codes.Internal, "Could not broadcast attestation: %v", err)
	}

	return &ethpb.AttestResponse{
		AttestationDataRoot: root[:],
	}, nil
}
