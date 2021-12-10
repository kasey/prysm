package get

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/api/client/openapi"
	"github.com/prysmaticlabs/prysm/proto/sniff"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"time"
)

var getStateFlags = struct {
	BeaconNodeHost string
	Timeout        string
	StateHex       string
	StateSavePath  string
}{}

var getStateCmd = &cli.Command{
	Name:   "state",
	Usage:  "Download a state identified by slot or epoch",
	Action: cliActionGetState,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "beacon-node-host",
			Usage:       "host:port for beacon node connection",
			Destination: &getStateFlags.BeaconNodeHost,
			Value:       "localhost:3500",
		},
		&cli.StringFlag{
			Name:        "http-timeout",
			Usage:       "timeout for http requests made to beacon-node-url (uses duration format, ex: 2m31s). default: 2m",
			Destination: &getStateFlags.Timeout,
			Value:       "2m",
		},
		&cli.StringFlag{
			Name:        "state-root",
			Usage:       "instead of epoch, state root (in 0x hex string format) can be used to retrieve from the beacon-node and save locally.",
			Destination: &getStateFlags.StateHex,
		},
	},
}

func saveStateByRoot(client *openapi.Client, root string) error {
	ctx := context.Background()
	fs, err := client.GetForkSchedule()
	if err != nil {
		return err
	}
	stateReader, err := client.GetStateByRoot(root)
	if err != nil {
		return err
	}
	stateBytes, err := io.ReadAll(stateReader)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read response body for state w/ root=%s from api", root))
	}
	epoch, err := sniff.EpochFromState(stateBytes)
	if err != nil {
		return errors.Wrap(err, "failed to detect the state epoch by inspecting the bytes")
	}
	version, err := fs.VersionForEpoch(epoch)
	if err != nil {
		return errors.Wrap(err, "failed to find a fork version spanning the epoch detected in state")
	}
	cf, err := sniff.ConfigForkForState(stateBytes)
	if err != nil {
		return errors.Wrap(err, "failed to detect the state version by inspecting the bytes")
	}
	if version != cf.Version {
		extra := fmt.Sprintf("version expected for state in epoch %d does not match detected version, detected=%#x, expected=%#x", epoch, cf.Version, version)
		return errors.Wrap(err, extra)
	}
	state, err := sniff.BeaconStateForConfigFork(stateBytes, cf)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal state with the detected type")
	}
	stateRoot, err := state.HashTreeRoot(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to compute HashTreeRoot for unmarshaled state value")
	}
	slot, err := cf.Config.SlotsPerEpoch.SafeMul(uint64(epoch))
	if err != nil {
		extra := fmt.Sprintf("overflow computing first slot of epoch, epoch=%d, slots_per_epoch=%d", epoch, cf.Config.SlotsPerEpoch)
		return errors.Wrap(err, extra)
	}
	statePath := fname("state", cf, uint64(slot), stateRoot)
	outBytes, err := state.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "error when marshaling state")
	}
	return os.WriteFile(statePath, outBytes, 0644)
}

func fname(prefix string, cf *sniff.ConfigFork, slot uint64, root [32]byte) string {
	return fmt.Sprintf("%s_%s_%s_%d-%#x.ssz", prefix, cf.ConfigName.String(), cf.Fork.String(), slot, root)
}

func cliActionGetState(c *cli.Context) error {
	f := getStateFlags
	opts := make([]openapi.ClientOpt, 0)
	log.Printf("--beacon-node-url=%s", f.BeaconNodeHost)
	timeout, err := time.ParseDuration(f.Timeout)
	if err != nil {
		return err
	}
	opts = append(opts, openapi.WithTimeout(timeout))
	client, err := openapi.NewClient(f.BeaconNodeHost, opts...)
	if err != nil {
		return err
	}
	return saveStateByRoot(client, f.StateHex)
}
