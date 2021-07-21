package ssz_static

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"testing"

	fssz "github.com/ferranbt/fastssz"
	"github.com/golang/snappy"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/spectest/utils"
)

// SSZRoots --
type SSZRoots struct {
	Root        string `json:"root"`
	SigningRoot string `json:"signing_root"`
}

// RunSSZStaticTests executes "ssz_static" tests.
func RunSSZStaticTests(t *testing.T, config string) {
	require.NoError(t, utils.SetConfig(t, config))

	testFolders, _ := utils.TestFolders(t, config, "phase0", "ssz_static")
	for _, folder := range testFolders {
		innerPath := path.Join("ssz_static", folder.Name(), "ssz_random")
		innerTestFolders, innerTestsFolderPath := utils.TestFolders(t, config, "phase0", innerPath)

		for _, innerFolder := range innerTestFolders {
			t.Run(path.Join(folder.Name(), innerFolder.Name()), func(t *testing.T) {
				if folder.Name() == "Eth1Block" {
					t.Skip("Unused type")
				}
				serializedBytes, err := testutil.BazelFileBytes(innerTestsFolderPath, innerFolder.Name(), "serialized.ssz_snappy")
				require.NoError(t, err)
				serializedSSZ, err := snappy.Decode(nil /* dst */, serializedBytes)
				require.NoError(t, err, "Failed to decompress")
				object, err := UnmarshalledSSZ(serializedSSZ, folder.Name())
				require.NoError(t, err, "Could not unmarshall serialized SSZ")

				rootsYamlFile, err := testutil.BazelFileBytes(innerTestsFolderPath, innerFolder.Name(), "roots.yaml")
				require.NoError(t, err)
				rootsYaml := &SSZRoots{}
				require.NoError(t, utils.UnmarshalYaml(rootsYamlFile, rootsYaml), "Failed to Unmarshal")

				// Custom hash tree root for beacon state.
				var htr func(interface{}) ([32]byte, error)
				if _, ok := object.(*pb.BeaconState); ok {
					htr = func(s interface{}) ([32]byte, error) {
						beaconState, err := v1.InitializeFromProto(s.(*pb.BeaconState))
						require.NoError(t, err)
						return beaconState.HashTreeRoot(context.Background())
					}
				} else {
					htr = func(s interface{}) ([32]byte, error) {
						sszObj, ok := s.(fssz.HashRoot)
						if !ok {
							return [32]byte{}, errors.New("could not get hash root, not compatible object")
						}
						return sszObj.HashTreeRoot()
					}
				}

				root, err := htr(object)
				require.NoError(t, err)
				rootBytes, err := hex.DecodeString(rootsYaml.Root[2:])
				require.NoError(t, err)
				require.DeepEqual(t, rootBytes, root[:], "Did not receive expected hash tree root")

				if rootsYaml.SigningRoot == "" {
					return
				}

				var signingRoot [32]byte
				if v, ok := object.(fssz.HashRoot); ok {
					signingRoot, err = v.HashTreeRoot()
				} else {
					t.Fatal("object does not meet fssz.HashRoot")
				}

				require.NoError(t, err)
				signingRootBytes, err := hex.DecodeString(rootsYaml.SigningRoot[2:])
				require.NoError(t, err)
				require.DeepEqual(t, signingRootBytes, signingRoot[:], "Did not receive expected signing root")
			})
		}
	}
}

type ExperimentalSSZ interface {
	XXUnmarshalSSZ(buf []byte) error
	XXMarshalSSZ() ([]byte, error)
	fssz.Unmarshaler
}

func GetInstanceByName(name string) (interface{}, error) {
	switch name {
	case "Attestation":
		return &ethpb.Attestation{}, nil
	case "AttestationData":
		return &ethpb.AttestationData{}, nil
	case "AttesterSlashing":
		return &ethpb.AttesterSlashing{}, nil
	case "AggregateAndProof":
		return &ethpb.AggregateAttestationAndProof{}, nil
	case "BeaconBlock":
		return &ethpb.BeaconBlock{}, nil
	case "BeaconBlockBody":
		return &ethpb.BeaconBlockBody{}, nil
	case "BeaconBlockHeader":
		return &ethpb.BeaconBlockHeader{}, nil
	case "BeaconState":
		return &pb.BeaconState{}, nil
	case "Checkpoint":
		return &ethpb.Checkpoint{}, nil
	case "Deposit":
		return &ethpb.Deposit{}, nil
	case "DepositMessage":
		return &pb.DepositMessage{}, nil
	case "DepositData":
		return &ethpb.Deposit_Data{}, nil
	case "Eth1Data":
		return &ethpb.Eth1Data{}, nil
	case "Fork":
		return &pb.Fork{}, nil
	case "ForkData":
		return &pb.ForkData{}, nil
	case "HistoricalBatch":
		return &pb.HistoricalBatch{}, nil
	case "IndexedAttestation":
		return &ethpb.IndexedAttestation{}, nil
	case "PendingAttestation":
		return &pb.PendingAttestation{}, nil
	case "ProposerSlashing":
		return &ethpb.ProposerSlashing{}, nil
	case "SignedAggregateAndProof":
		return &ethpb.SignedAggregateAttestationAndProof{}, nil
	case "SignedBeaconBlock":
		return &ethpb.SignedBeaconBlock{}, nil
	case "SignedBeaconBlockHeader":
		return &ethpb.SignedBeaconBlockHeader{}, nil
	case "SignedVoluntaryExit":
		return &ethpb.SignedVoluntaryExit{}, nil
	case "SigningData":
		return &pb.SigningData{}, nil
	case "Validator":
		return &ethpb.Validator{}, nil
	case "VoluntaryExit":
		return &ethpb.VoluntaryExit{}, nil
	default:
		return nil, errors.New("type not found")
	}
}

// UnmarshalledSSZ unmarshalls serialized input.
func UnmarshalledSSZ(serializedBytes []byte, folderName string) (interface{}, error) {
	obj, err := GetInstanceByName(folderName)
	if err != nil {
		return nil, err
	}
	o, ok := obj.(ExperimentalSSZ)
	if !ok {
		return nil, fmt.Errorf("%s fails ExperimentalSSZ interface check", folderName)
	}
	err = o.XXUnmarshalSSZ(serializedBytes)
	if err != nil {
			  return nil, err
			  }
	marshalled, err := o.XXMarshalSSZ()
	if err != nil {
			  return nil, err
			  }

	// get a new instance for unmarshaling
	obj, err = GetInstanceByName(folderName)
	if err != nil {
			  return nil, err
			  }
	oo, ok := obj.(fssz.Unmarshaler)
	if !ok {
	   return nil, fmt.Errorf("%s implements ExperimentalSSZ but not fssz.Unmarshaler?", folderName)
	   }
	// make sure we can unmarshal with fastssz code
	err = oo.UnmarshalSSZ(marshalled)
	return oo, err
}
