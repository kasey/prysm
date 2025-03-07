package blocks

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
)

var (
	// ErrUnsupportedSignedBeaconBlock is returned when the struct type is not a supported signed
	// beacon block type.
	ErrUnsupportedSignedBeaconBlock = errors.New("unsupported signed beacon block")
	// errUnsupportedBeaconBlock is returned when the struct type is not a supported beacon block
	// type.
	errUnsupportedBeaconBlock = errors.New("unsupported beacon block")
	// errUnsupportedBeaconBlockBody is returned when the struct type is not a supported beacon block body
	// type.
	errUnsupportedBeaconBlockBody = errors.New("unsupported beacon block body")
	// ErrNilObject is returned in a constructor when the underlying object is nil.
	ErrNilObject = errors.New("received nil object")
	// ErrNilSignedBeaconBlock is returned when a nil signed beacon block is received.
	ErrNilSignedBeaconBlock = errors.New("signed beacon block can't be nil")
	// ErrNilBeaconBlock is returned when a nil beacon block is received.
	ErrNilBeaconBlock              = errors.New("beacon block can't be nil")
	errNonBlindedSignedBeaconBlock = errors.New("can only build signed beacon block from blinded format")
)

// NewSignedBeaconBlock creates a signed beacon block from a protobuf signed beacon block.
func NewSignedBeaconBlock(i interface{}) (interfaces.SignedBeaconBlock, error) {
	switch b := i.(type) {
	case nil:
		return nil, ErrNilObject
	case *eth.GenericSignedBeaconBlock_Phase0:
		return initSignedBlockFromProtoPhase0(b.Phase0)
	case *eth.SignedBeaconBlock:
		return initSignedBlockFromProtoPhase0(b)
	case *eth.GenericSignedBeaconBlock_Altair:
		return initSignedBlockFromProtoAltair(b.Altair)
	case *eth.SignedBeaconBlockAltair:
		return initSignedBlockFromProtoAltair(b)
	case *eth.GenericSignedBeaconBlock_Bellatrix:
		return initSignedBlockFromProtoBellatrix(b.Bellatrix)
	case *eth.SignedBeaconBlockBellatrix:
		return initSignedBlockFromProtoBellatrix(b)
	case *eth.GenericSignedBeaconBlock_BlindedBellatrix:
		return initBlindedSignedBlockFromProtoBellatrix(b.BlindedBellatrix)
	case *eth.SignedBlindedBeaconBlockBellatrix:
		return initBlindedSignedBlockFromProtoBellatrix(b)
	case *eth.GenericSignedBeaconBlock_Capella:
		return initSignedBlockFromProtoCapella(b.Capella)
	case *eth.SignedBeaconBlockCapella:
		return initSignedBlockFromProtoCapella(b)
	case *eth.GenericSignedBeaconBlock_BlindedCapella:
		return initBlindedSignedBlockFromProtoCapella(b.BlindedCapella)
	case *eth.SignedBlindedBeaconBlockCapella:
		return initBlindedSignedBlockFromProtoCapella(b)
	case *eth.GenericSignedBeaconBlock_Deneb:
		return initSignedBlockFromProtoDeneb(b.Deneb.Block)
	case *eth.SignedBeaconBlockDeneb:
		return initSignedBlockFromProtoDeneb(b)
	case *eth.SignedBlindedBeaconBlockDeneb:
		return initBlindedSignedBlockFromProtoDeneb(b)
	case *eth.GenericSignedBeaconBlock_BlindedDeneb:
		return initBlindedSignedBlockFromProtoDeneb(b.BlindedDeneb)
	case *eth.GenericSignedBeaconBlock_Electra:
		return initSignedBlockFromProtoElectra(b.Electra.Block)
	case *eth.SignedBeaconBlockElectra:
		return initSignedBlockFromProtoElectra(b)
	case *eth.SignedBlindedBeaconBlockElectra:
		return initBlindedSignedBlockFromProtoElectra(b)
	case *eth.GenericSignedBeaconBlock_BlindedElectra:
		return initBlindedSignedBlockFromProtoElectra(b.BlindedElectra)
	case *eth.GenericSignedBeaconBlock_Fulu:
		return initSignedBlockFromProtoFulu(b.Fulu.Block)
	case *eth.SignedBeaconBlockFulu:
		return initSignedBlockFromProtoFulu(b)
	case *eth.SignedBlindedBeaconBlockFulu:
		return initBlindedSignedBlockFromProtoFulu(b)
	case *eth.GenericSignedBeaconBlock_BlindedFulu:
		return initBlindedSignedBlockFromProtoFulu(b.BlindedFulu)
	default:
		return nil, errors.Wrapf(ErrUnsupportedSignedBeaconBlock, "unable to create block from type %T", i)
	}
}

// NewBeaconBlock creates a beacon block from a protobuf beacon block.
func NewBeaconBlock(i interface{}) (interfaces.ReadOnlyBeaconBlock, error) {
	switch b := i.(type) {
	case nil:
		return nil, ErrNilObject
	case *eth.GenericBeaconBlock_Phase0:
		return initBlockFromProtoPhase0(b.Phase0)
	case *eth.BeaconBlock:
		return initBlockFromProtoPhase0(b)
	case *eth.GenericBeaconBlock_Altair:
		return initBlockFromProtoAltair(b.Altair)
	case *eth.BeaconBlockAltair:
		return initBlockFromProtoAltair(b)
	case *eth.GenericBeaconBlock_Bellatrix:
		return initBlockFromProtoBellatrix(b.Bellatrix)
	case *eth.BeaconBlockBellatrix:
		return initBlockFromProtoBellatrix(b)
	case *eth.GenericBeaconBlock_BlindedBellatrix:
		return initBlindedBlockFromProtoBellatrix(b.BlindedBellatrix)
	case *eth.BlindedBeaconBlockBellatrix:
		return initBlindedBlockFromProtoBellatrix(b)
	case *eth.GenericBeaconBlock_Capella:
		return initBlockFromProtoCapella(b.Capella)
	case *eth.BeaconBlockCapella:
		return initBlockFromProtoCapella(b)
	case *eth.GenericBeaconBlock_BlindedCapella:
		return initBlindedBlockFromProtoCapella(b.BlindedCapella)
	case *eth.BlindedBeaconBlockCapella:
		return initBlindedBlockFromProtoCapella(b)
	case *eth.GenericBeaconBlock_Deneb:
		return initBlockFromProtoDeneb(b.Deneb.Block)
	case *eth.BeaconBlockDeneb:
		return initBlockFromProtoDeneb(b)
	case *eth.BlindedBeaconBlockDeneb:
		return initBlindedBlockFromProtoDeneb(b)
	case *eth.GenericBeaconBlock_BlindedDeneb:
		return initBlindedBlockFromProtoDeneb(b.BlindedDeneb)
	case *eth.GenericBeaconBlock_Electra:
		return initBlockFromProtoElectra(b.Electra.Block)
	case *eth.BeaconBlockElectra:
		return initBlockFromProtoElectra(b)
	case *eth.BlindedBeaconBlockElectra:
		return initBlindedBlockFromProtoElectra(b)
	case *eth.GenericBeaconBlock_BlindedElectra:
		return initBlindedBlockFromProtoElectra(b.BlindedElectra)
	case *eth.GenericBeaconBlock_Fulu:
		return initBlockFromProtoFulu(b.Fulu.Block)
	case *eth.BeaconBlockFulu:
		return initBlockFromProtoFulu(b)
	case *eth.BlindedBeaconBlockFulu:
		return initBlindedBlockFromProtoFulu(b)
	case *eth.GenericBeaconBlock_BlindedFulu:
		return initBlindedBlockFromProtoFulu(b.BlindedFulu)
	default:
		return nil, errors.Wrapf(errUnsupportedBeaconBlock, "unable to create block from type %T", i)
	}
}

// NewBeaconBlockBody creates a beacon block body from a protobuf beacon block body.
func NewBeaconBlockBody(i interface{}) (interfaces.ReadOnlyBeaconBlockBody, error) {
	switch b := i.(type) {
	case nil:
		return nil, ErrNilObject
	case *eth.BeaconBlockBody:
		return initBlockBodyFromProtoPhase0(b)
	case *eth.BeaconBlockBodyAltair:
		return initBlockBodyFromProtoAltair(b)
	case *eth.BeaconBlockBodyBellatrix:
		return initBlockBodyFromProtoBellatrix(b)
	case *eth.BlindedBeaconBlockBodyBellatrix:
		return initBlindedBlockBodyFromProtoBellatrix(b)
	case *eth.BeaconBlockBodyCapella:
		return initBlockBodyFromProtoCapella(b)
	case *eth.BlindedBeaconBlockBodyCapella:
		return initBlindedBlockBodyFromProtoCapella(b)
	case *eth.BeaconBlockBodyDeneb:
		return initBlockBodyFromProtoDeneb(b)
	case *eth.BlindedBeaconBlockBodyDeneb:
		return initBlindedBlockBodyFromProtoDeneb(b)
	case *eth.BeaconBlockBodyElectra:
		return initBlockBodyFromProtoElectra(b)
	case *eth.BlindedBeaconBlockBodyElectra:
		return initBlindedBlockBodyFromProtoElectra(b)
	case *eth.BeaconBlockBodyFulu:
		return initBlockBodyFromProtoFulu(b)
	case *eth.BlindedBeaconBlockBodyFulu:
		return initBlindedBlockBodyFromProtoFulu(b)
	default:
		return nil, errors.Wrapf(errUnsupportedBeaconBlockBody, "unable to create block body from type %T", i)
	}
}

// BuildSignedBeaconBlock assembles a block.ReadOnlySignedBeaconBlock interface compatible struct from a
// given beacon block and the appropriate signature. This method may be used to easily create a
// signed beacon block.
func BuildSignedBeaconBlock(blk interfaces.ReadOnlyBeaconBlock, signature []byte) (interfaces.SignedBeaconBlock, error) {
	pb, err := blk.Proto()
	if err != nil {
		return nil, err
	}

	switch blk.Version() {
	case version.Phase0:
		pb, ok := pb.(*eth.BeaconBlock)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlock{Block: pb, Signature: signature})
	case version.Altair:
		pb, ok := pb.(*eth.BeaconBlockAltair)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockAltair{Block: pb, Signature: signature})
	case version.Bellatrix:
		if blk.IsBlinded() {
			pb, ok := pb.(*eth.BlindedBeaconBlockBellatrix)
			if !ok {
				return nil, errIncorrectBlockVersion
			}
			return NewSignedBeaconBlock(&eth.SignedBlindedBeaconBlockBellatrix{Block: pb, Signature: signature})
		}
		pb, ok := pb.(*eth.BeaconBlockBellatrix)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockBellatrix{Block: pb, Signature: signature})
	case version.Capella:
		if blk.IsBlinded() {
			pb, ok := pb.(*eth.BlindedBeaconBlockCapella)
			if !ok {
				return nil, errIncorrectBlockVersion
			}
			return NewSignedBeaconBlock(&eth.SignedBlindedBeaconBlockCapella{Block: pb, Signature: signature})
		}
		pb, ok := pb.(*eth.BeaconBlockCapella)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockCapella{Block: pb, Signature: signature})
	case version.Deneb:
		if blk.IsBlinded() {
			pb, ok := pb.(*eth.BlindedBeaconBlockDeneb)
			if !ok {
				return nil, errIncorrectBlockVersion
			}
			return NewSignedBeaconBlock(&eth.SignedBlindedBeaconBlockDeneb{Message: pb, Signature: signature})
		}
		pb, ok := pb.(*eth.BeaconBlockDeneb)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockDeneb{Block: pb, Signature: signature})
	case version.Electra:
		if blk.IsBlinded() {
			pb, ok := pb.(*eth.BlindedBeaconBlockElectra)
			if !ok {
				return nil, errIncorrectBlockVersion
			}
			return NewSignedBeaconBlock(&eth.SignedBlindedBeaconBlockElectra{Message: pb, Signature: signature})
		}
		pb, ok := pb.(*eth.BeaconBlockElectra)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockElectra{Block: pb, Signature: signature})
	case version.Fulu:
		if blk.IsBlinded() {
			pb, ok := pb.(*eth.BlindedBeaconBlockFulu)
			if !ok {
				return nil, errIncorrectBlockVersion
			}
			return NewSignedBeaconBlock(&eth.SignedBlindedBeaconBlockFulu{Message: pb, Signature: signature})
		}
		pb, ok := pb.(*eth.BeaconBlockFulu)
		if !ok {
			return nil, errIncorrectBlockVersion
		}
		return NewSignedBeaconBlock(&eth.SignedBeaconBlockFulu{Block: pb, Signature: signature})
	default:
		return nil, errUnsupportedBeaconBlock
	}
}

func getWrappedPayload(payload interface{}) (wrappedPayload interfaces.ExecutionData, wrapErr error) {
	switch p := payload.(type) {
	case *enginev1.ExecutionPayload:
		wrappedPayload, wrapErr = WrappedExecutionPayload(p)
	case *enginev1.ExecutionPayloadCapella:
		wrappedPayload, wrapErr = WrappedExecutionPayloadCapella(p)
	case *enginev1.ExecutionPayloadDeneb:
		wrappedPayload, wrapErr = WrappedExecutionPayloadDeneb(p)
	default:
		wrappedPayload, wrapErr = nil, fmt.Errorf("%T is not a type of execution payload", p)
	}
	return wrappedPayload, wrapErr
}

func checkPayloadAgainstHeader(wrappedPayload, payloadHeader interfaces.ExecutionData) error {
	empty, err := IsEmptyExecutionData(wrappedPayload)
	if err != nil {
		return err
	}
	if empty {
		return nil
	}
	payloadRoot, err := wrappedPayload.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not hash tree root execution payload")
	}
	payloadHeaderRoot, err := payloadHeader.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not hash tree root payload header")
	}
	if payloadRoot != payloadHeaderRoot {
		return fmt.Errorf(
			"payload %#x and header %#x roots do not match",
			payloadRoot,
			payloadHeaderRoot,
		)
	}
	return nil
}

// BuildSignedBeaconBlockFromExecutionPayload takes a signed, blinded beacon block and converts into
// a full, signed beacon block by specifying an execution payload.
// nolint:gocognit
func BuildSignedBeaconBlockFromExecutionPayload(blk interfaces.ReadOnlySignedBeaconBlock, payload interface{}) (interfaces.SignedBeaconBlock, error) {
	if err := BeaconBlockIsNil(blk); err != nil {
		return nil, err
	}
	if !blk.IsBlinded() {
		return nil, errNonBlindedSignedBeaconBlock
	}
	b := blk.Block()
	payloadHeader, err := b.Body().Execution()
	if err != nil {
		return nil, errors.Wrap(err, "could not get execution payload header")
	}

	wrappedPayload, err := getWrappedPayload(payload)
	if err != nil {
		return nil, err
	}
	if err := checkPayloadAgainstHeader(wrappedPayload, payloadHeader); err != nil {
		return nil, err
	}
	syncAgg, err := b.Body().SyncAggregate()
	if err != nil {
		return nil, errors.Wrap(err, "could not get sync aggregate from block body")
	}
	parentRoot := b.ParentRoot()
	stateRoot := b.StateRoot()
	randaoReveal := b.Body().RandaoReveal()
	graffiti := b.Body().Graffiti()
	sig := blk.Signature()

	var fullBlock interface{}
	switch blk.Version() {
	case version.Bellatrix:
		p, ok := payload.(*enginev1.ExecutionPayload)
		if !ok {
			return nil, errors.New("payload not of Bellatrix type")
		}
		var atts []*eth.Attestation
		if b.Body().Attestations() != nil {
			atts = make([]*eth.Attestation, len(b.Body().Attestations()))
			for i, att := range b.Body().Attestations() {
				a, ok := att.(*eth.Attestation)
				if !ok {
					return nil, fmt.Errorf("attestation has wrong type (expected %T, got %T)", &eth.Attestation{}, att)
				}
				atts[i] = a
			}
		}
		var attSlashings []*eth.AttesterSlashing
		if b.Body().AttesterSlashings() != nil {
			attSlashings = make([]*eth.AttesterSlashing, len(b.Body().AttesterSlashings()))
			for i, slashing := range b.Body().AttesterSlashings() {
				s, ok := slashing.(*eth.AttesterSlashing)
				if !ok {
					return nil, fmt.Errorf("attester slashing has wrong type (expected %T, got %T)", &eth.AttesterSlashing{}, slashing)
				}
				attSlashings[i] = s
			}
		}
		fullBlock = &eth.SignedBeaconBlockBellatrix{
			Block: &eth.BeaconBlockBellatrix{
				Slot:          b.Slot(),
				ProposerIndex: b.ProposerIndex(),
				ParentRoot:    parentRoot[:],
				StateRoot:     stateRoot[:],
				Body: &eth.BeaconBlockBodyBellatrix{
					RandaoReveal:      randaoReveal[:],
					Eth1Data:          b.Body().Eth1Data(),
					Graffiti:          graffiti[:],
					ProposerSlashings: b.Body().ProposerSlashings(),
					AttesterSlashings: attSlashings,
					Attestations:      atts,
					Deposits:          b.Body().Deposits(),
					VoluntaryExits:    b.Body().VoluntaryExits(),
					SyncAggregate:     syncAgg,
					ExecutionPayload:  p,
				},
			},
			Signature: sig[:],
		}
	case version.Capella:
		p, ok := payload.(*enginev1.ExecutionPayloadCapella)
		if !ok {
			return nil, errors.New("payload not of Capella type")
		}
		blsToExecutionChanges, err := b.Body().BLSToExecutionChanges()
		if err != nil {
			return nil, err
		}
		var atts []*eth.Attestation
		if b.Body().Attestations() != nil {
			atts = make([]*eth.Attestation, len(b.Body().Attestations()))
			for i, att := range b.Body().Attestations() {
				a, ok := att.(*eth.Attestation)
				if !ok {
					return nil, fmt.Errorf("attestation has wrong type (expected %T, got %T)", &eth.Attestation{}, att)
				}
				atts[i] = a
			}
		}
		var attSlashings []*eth.AttesterSlashing
		if b.Body().AttesterSlashings() != nil {
			attSlashings = make([]*eth.AttesterSlashing, len(b.Body().AttesterSlashings()))
			for i, slashing := range b.Body().AttesterSlashings() {
				s, ok := slashing.(*eth.AttesterSlashing)
				if !ok {
					return nil, fmt.Errorf("attester slashing has wrong type (expected %T, got %T)", &eth.AttesterSlashing{}, slashing)
				}
				attSlashings[i] = s
			}
		}
		fullBlock = &eth.SignedBeaconBlockCapella{
			Block: &eth.BeaconBlockCapella{
				Slot:          b.Slot(),
				ProposerIndex: b.ProposerIndex(),
				ParentRoot:    parentRoot[:],
				StateRoot:     stateRoot[:],
				Body: &eth.BeaconBlockBodyCapella{
					RandaoReveal:          randaoReveal[:],
					Eth1Data:              b.Body().Eth1Data(),
					Graffiti:              graffiti[:],
					ProposerSlashings:     b.Body().ProposerSlashings(),
					AttesterSlashings:     attSlashings,
					Attestations:          atts,
					Deposits:              b.Body().Deposits(),
					VoluntaryExits:        b.Body().VoluntaryExits(),
					SyncAggregate:         syncAgg,
					ExecutionPayload:      p,
					BlsToExecutionChanges: blsToExecutionChanges,
				},
			},
			Signature: sig[:],
		}
	case version.Deneb:
		p, ok := payload.(*enginev1.ExecutionPayloadDeneb)
		if !ok {
			return nil, errors.New("payload not of Deneb type")
		}
		blsToExecutionChanges, err := b.Body().BLSToExecutionChanges()
		if err != nil {
			return nil, err
		}
		commitments, err := b.Body().BlobKzgCommitments()
		if err != nil {
			return nil, err
		}
		var atts []*eth.Attestation
		if b.Body().Attestations() != nil {
			atts = make([]*eth.Attestation, len(b.Body().Attestations()))
			for i, att := range b.Body().Attestations() {
				a, ok := att.(*eth.Attestation)
				if !ok {
					return nil, fmt.Errorf("attestation has wrong type (expected %T, got %T)", &eth.Attestation{}, att)
				}
				atts[i] = a
			}
		}
		var attSlashings []*eth.AttesterSlashing
		if b.Body().AttesterSlashings() != nil {
			attSlashings = make([]*eth.AttesterSlashing, len(b.Body().AttesterSlashings()))
			for i, slashing := range b.Body().AttesterSlashings() {
				s, ok := slashing.(*eth.AttesterSlashing)
				if !ok {
					return nil, fmt.Errorf("attester slashing has wrong type (expected %T, got %T)", &eth.AttesterSlashing{}, slashing)
				}
				attSlashings[i] = s
			}
		}
		fullBlock = &eth.SignedBeaconBlockDeneb{
			Block: &eth.BeaconBlockDeneb{
				Slot:          b.Slot(),
				ProposerIndex: b.ProposerIndex(),
				ParentRoot:    parentRoot[:],
				StateRoot:     stateRoot[:],
				Body: &eth.BeaconBlockBodyDeneb{
					RandaoReveal:          randaoReveal[:],
					Eth1Data:              b.Body().Eth1Data(),
					Graffiti:              graffiti[:],
					ProposerSlashings:     b.Body().ProposerSlashings(),
					AttesterSlashings:     attSlashings,
					Attestations:          atts,
					Deposits:              b.Body().Deposits(),
					VoluntaryExits:        b.Body().VoluntaryExits(),
					SyncAggregate:         syncAgg,
					ExecutionPayload:      p,
					BlsToExecutionChanges: blsToExecutionChanges,
					BlobKzgCommitments:    commitments,
				},
			},
			Signature: sig[:],
		}
	case version.Electra:
		p, ok := payload.(*enginev1.ExecutionPayloadDeneb)
		if !ok {
			return nil, errors.New("payload not of Electra type")
		}
		blsToExecutionChanges, err := b.Body().BLSToExecutionChanges()
		if err != nil {
			return nil, err
		}
		commitments, err := b.Body().BlobKzgCommitments()
		if err != nil {
			return nil, err
		}
		var atts []*eth.AttestationElectra
		if b.Body().Attestations() != nil {
			atts = make([]*eth.AttestationElectra, len(b.Body().Attestations()))
			for i, att := range b.Body().Attestations() {
				a, ok := att.(*eth.AttestationElectra)
				if !ok {
					return nil, fmt.Errorf("attestation has wrong type (expected %T, got %T)", &eth.Attestation{}, att)
				}
				atts[i] = a
			}
		}
		var attSlashings []*eth.AttesterSlashingElectra
		if b.Body().AttesterSlashings() != nil {
			attSlashings = make([]*eth.AttesterSlashingElectra, len(b.Body().AttesterSlashings()))
			for i, slashing := range b.Body().AttesterSlashings() {
				s, ok := slashing.(*eth.AttesterSlashingElectra)
				if !ok {
					return nil, fmt.Errorf("attester slashing has wrong type (expected %T, got %T)", &eth.AttesterSlashing{}, slashing)
				}
				attSlashings[i] = s
			}
		}

		er, err := b.Body().ExecutionRequests()
		if err != nil {
			return nil, err
		}

		fullBlock = &eth.SignedBeaconBlockElectra{
			Block: &eth.BeaconBlockElectra{
				Slot:          b.Slot(),
				ProposerIndex: b.ProposerIndex(),
				ParentRoot:    parentRoot[:],
				StateRoot:     stateRoot[:],
				Body: &eth.BeaconBlockBodyElectra{
					RandaoReveal:          randaoReveal[:],
					Eth1Data:              b.Body().Eth1Data(),
					Graffiti:              graffiti[:],
					ProposerSlashings:     b.Body().ProposerSlashings(),
					AttesterSlashings:     attSlashings,
					Attestations:          atts,
					Deposits:              b.Body().Deposits(),
					VoluntaryExits:        b.Body().VoluntaryExits(),
					SyncAggregate:         syncAgg,
					ExecutionPayload:      p,
					BlsToExecutionChanges: blsToExecutionChanges,
					BlobKzgCommitments:    commitments,
					ExecutionRequests:     er,
				},
			},
			Signature: sig[:],
		}
	case version.Fulu:
		p, ok := payload.(*enginev1.ExecutionPayloadDeneb)
		if !ok {
			return nil, errors.New("payload not of Fulu type")
		}
		blsToExecutionChanges, err := b.Body().BLSToExecutionChanges()
		if err != nil {
			return nil, err
		}
		commitments, err := b.Body().BlobKzgCommitments()
		if err != nil {
			return nil, err
		}
		var atts []*eth.AttestationElectra
		if b.Body().Attestations() != nil {
			atts = make([]*eth.AttestationElectra, len(b.Body().Attestations()))
			for i, att := range b.Body().Attestations() {
				a, ok := att.(*eth.AttestationElectra)
				if !ok {
					return nil, fmt.Errorf("attestation has wrong type (expected %T, got %T)", &eth.Attestation{}, att)
				}
				atts[i] = a
			}
		}
		var attSlashings []*eth.AttesterSlashingElectra
		if b.Body().AttesterSlashings() != nil {
			attSlashings = make([]*eth.AttesterSlashingElectra, len(b.Body().AttesterSlashings()))
			for i, slashing := range b.Body().AttesterSlashings() {
				s, ok := slashing.(*eth.AttesterSlashingElectra)
				if !ok {
					return nil, fmt.Errorf("attester slashing has wrong type (expected %T, got %T)", &eth.AttesterSlashing{}, slashing)
				}
				attSlashings[i] = s
			}
		}

		er, err := b.Body().ExecutionRequests()
		if err != nil {
			return nil, err
		}

		fullBlock = &eth.SignedBeaconBlockFulu{
			Block: &eth.BeaconBlockFulu{
				Slot:          b.Slot(),
				ProposerIndex: b.ProposerIndex(),
				ParentRoot:    parentRoot[:],
				StateRoot:     stateRoot[:],
				Body: &eth.BeaconBlockBodyFulu{
					RandaoReveal:          randaoReveal[:],
					Eth1Data:              b.Body().Eth1Data(),
					Graffiti:              graffiti[:],
					ProposerSlashings:     b.Body().ProposerSlashings(),
					AttesterSlashings:     attSlashings,
					Attestations:          atts,
					Deposits:              b.Body().Deposits(),
					VoluntaryExits:        b.Body().VoluntaryExits(),
					SyncAggregate:         syncAgg,
					ExecutionPayload:      p,
					BlsToExecutionChanges: blsToExecutionChanges,
					BlobKzgCommitments:    commitments,
					ExecutionRequests:     er,
				},
			},
			Signature: sig[:],
		}
	default:
		return nil, errors.New("Block not of known type")
	}

	return NewSignedBeaconBlock(fullBlock)
}

// BeaconBlockContainerToSignedBeaconBlock converts BeaconBlockContainer (API response) to a SignedBeaconBlock.
// This is particularly useful for using the values from API calls.
func BeaconBlockContainerToSignedBeaconBlock(obj *eth.BeaconBlockContainer) (interfaces.ReadOnlySignedBeaconBlock, error) {
	switch obj.Block.(type) {
	case *eth.BeaconBlockContainer_BlindedDenebBlock:
		return NewSignedBeaconBlock(obj.GetBlindedDenebBlock())
	case *eth.BeaconBlockContainer_DenebBlock:
		return NewSignedBeaconBlock(obj.GetDenebBlock())
	case *eth.BeaconBlockContainer_BlindedCapellaBlock:
		return NewSignedBeaconBlock(obj.GetBlindedCapellaBlock())
	case *eth.BeaconBlockContainer_CapellaBlock:
		return NewSignedBeaconBlock(obj.GetCapellaBlock())
	case *eth.BeaconBlockContainer_BlindedBellatrixBlock:
		return NewSignedBeaconBlock(obj.GetBlindedBellatrixBlock())
	case *eth.BeaconBlockContainer_BellatrixBlock:
		return NewSignedBeaconBlock(obj.GetBellatrixBlock())
	case *eth.BeaconBlockContainer_AltairBlock:
		return NewSignedBeaconBlock(obj.GetAltairBlock())
	case *eth.BeaconBlockContainer_Phase0Block:
		return NewSignedBeaconBlock(obj.GetPhase0Block())
	default:
		return nil, errors.New("container block type not recognized")
	}
}
