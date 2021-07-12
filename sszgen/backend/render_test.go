package backend

import (
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/sszgen/types"
)

var generator_generateFixture = `package derp

import (
	"fmt"
	derp "github.com/prysmaticlabs/derp/derp"
)

func main() {
	fmt.printf("hello world")
}
`

func TestGenerator_Generate(t *testing.T) {
	gc := &generatedCode{
		blocks:  []string{"func main() {\n\tfmt.printf(\"hello world\")\n}"},
		imports: map[string]string{
			"github.com/prysmaticlabs/derp/derp": "derp",
			"fmt": "",
		},
	}
	g := &Generator{packageName: "github.com/prysmaticlabs/derp"}
	g.gc = append(g.gc, gc)
	rendered, err := g.Render()
	require.NoError(t, err)
	require.Equal(t, generator_generateFixture, string(rendered))
}

var generator_generateBeaconStateFixture = `package derp

import (
	ssz "github.com/ferranbt/fastssz"
)

func (c *BeaconState) SizeSSZ() (size int) {
	size := 2687377
	size += len(c.HistoricalRoots) * 32
	size += len(c.Eth1DataVotes) * 72
	size += len(c.Validators) * 121
	size += len(c.Balances) * 8
	size += func() int {
		s := 0
		for _, o := range c.PreviousEpochAttestations {
			s += 4
			s += c.PreviousEpochAttestations.SizeSSZ()
		}
		return s
	}()
	size += func() int {
		s := 0
		for _, o := range c.CurrentEpochAttestations {
			s += 4
			s += c.CurrentEpochAttestations.SizeSSZ()
		}
		return s
	}()
	return size
}
`

func TestGenerator_GenerateBeaconState(t *testing.T) {
	g := &Generator{packageName: "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"}
	g.Generate(testFixBeaconState)
	rendered, err := g.Render()
	require.NoError(t, err)
	require.Equal(t, generator_generateBeaconStateFixture, string(rendered))
}

func TestImportAlias(t *testing.T) {
	cases := []struct{
		packageName string
		alias string
	}{
		{
			packageName: "github.com/derp/derp",
			alias: "derp_derp",
		},
		{
			packageName: "text/template",
			alias: "text_template",
		},
		{
			packageName: "fmt",
			alias: "fmt",
		},
	}
	for _, c := range cases {
		require.Equal(t, importAlias(c.packageName), c.alias)
	}
}

func TestExtractImportsFromContainerFields(t *testing.T) {
	vc := testFixBeaconState.(*types.ValueContainer)
	imports := extractImportsFromContainerFields(vc.Contents)
	require.Equal(t, 4, len(imports))
	require.Equal(t, "prysmaticlabs_eth2_types", imports["github.com/prysmaticlabs/eth2-types"])
	require.Equal(t, "prysmaticlabs_prysm_proto_beacon_p2p_v1", imports["github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"])
	require.Equal(t, "prysmaticlabs_prysm_proto_eth_v1alpha1", imports["github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"])
	require.Equal(t, "prysmaticlabs_go_bitfield", imports["github.com/prysmaticlabs/go-bitfield"])
}