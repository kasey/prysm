package checkpoint

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"net"
	"net/url"
	"time"

	"github.com/prysmaticlabs/prysm/api/client/openapi"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var latestFlags = struct{
	BeaconNodeHost string
	Timeout        string
}{}

var latestCmd = &cli.Command{
	Name: "latest",
	Usage: "Connect to a beacon-node server and print the block_root:epoch for the latest checkpoint.",
	Action: cliActionLatest,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "beacon-node-host",
			Usage: "host:port for beacon node to query",
			Destination: &latestFlags.BeaconNodeHost,
		},
		&cli.StringFlag{
			Name: "http-timeout",
			Usage: "timeout for http requests made to beacon-node-url (uses duration format, ex: 2m31s). default: 2m",
			Destination: &latestFlags.Timeout,
			Value: "2m",
		},
	},
}

func cliActionLatest(c *cli.Context) error {
	f := latestFlags
	opts := make([]openapi.ClientOpt, 0)
	log.Printf("beacon-node-url=%s", f.BeaconNodeHost)
	timeout, err := time.ParseDuration(f.Timeout)
	if err != nil {
		return err
	}
	opts = append(opts, openapi.WithTimeout(timeout))
	validatedHost, err := validHostname(latestFlags.BeaconNodeHost)
	if err !=  nil {
		return err
	}
	log.Printf("host-port=%s", validatedHost)
	client, err := openapi.NewClient(validatedHost, opts...)
	if err !=  nil {
		return err
	}
	wsc, err := client.GetWeakSubjectivityCheckpoint()
	if err != nil {
		return err
	}
	log.Print("writing weak subjectivity results to stdout")
	fmt.Printf("epoch: %d\nblock_root: %s\nstate_root: %s\n", int(wsc.Epoch), hexutil.Encode(wsc.BlockRoot), hexutil.Encode(wsc.StateRoot))

	return nil
}

func validHostname(h string) (string, error){
	// try to parse as url (being permissive)
	u, err := url.Parse(h)
	if err == nil && u.Host != "" {
		return u.Host, nil
	}
	// try to parse as host:port
	host, port, err := net.SplitHostPort(h)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}

/*
+       "encoding/hex"
        "flag"
        "fmt"
+       "github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"

        types "github.com/prysmaticlabs/eth2-types"
        "github.com/prysmaticlabs/prysm/beacon-chain/core/transition/interop"
@@ -37,7 +39,12 @@ func main() {
        if len(roots) != 1 {
                fmt.Printf("Expected 1 block root for slot %d, got %d roots", *state, len(roots))
        }
-       s, err := d.State(ctx, roots[0])
+       var blockRoot [32]byte = roots[0]
+       dst := make([]byte, hex.EncodedLen(len(blockRoot)))
+       hex.Encode(dst, blockRoot[:])
+       fmt.Printf("root for slot %d = %s\n", slot, string(dst))
+       stateGen := stategen.New(d)
+       s, err := stateGen.StateByRoot(ctx, blockRoot)

 */