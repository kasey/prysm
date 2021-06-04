package sszgen

import (
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func newFixedTestIndexer(sourcePackage string, index map[string]PackageParser) *PackageIndex {
	return &PackageIndex{
		sourcePackage: sourcePackage,
		index: index,
		structCache: make(map[[2]string]*TypeSpec),
	}
}

func newTestIndexer(sourcePackage string) *PackageIndex {
	return &PackageIndex{
		sourcePackage: sourcePackage,
		index: make(map[string]PackageParser),
		structCache: make(map[[2]string]*TypeSpec),
	}
}

func TestAddGet(t *testing.T) {
	packageName := "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	pi := newTestIndexer(packageName)
	sourceFiles := []string{"testdata/simple.go"}
	pp, err := newTestPackageParser(packageName, sourceFiles)
	require.NoError(t, err)
	pi.addToIndex(packageName, pp)
	parser, err := pi.getParser(packageName)
	require.NoError(t, err)
	_, err = parser.GetStruct("NoImports")
	require.NoError(t, err)
	_, err = pi.GetStruct(packageName, "NoImports")
	require.NoError(t, err)
}
