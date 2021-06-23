package backend

import (
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/require"
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
	g := &Generator{packageName: "github.com/prysmaticlabs/derp"}
	g.Generate(testFixBeaconState)
	rendered, err := g.Render()
	require.NoError(t, err)
	require.Equal(t, generator_generateBeaconStateFixture, string(rendered))
}
