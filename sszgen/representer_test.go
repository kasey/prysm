package sszgen

import (
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestGetSimpleRepresentation(t *testing.T) {
	packageName := "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	sourceFiles := []string{"testdata/simple.go"}
	pp, err:= newTestPackageParser(packageName, sourceFiles)
	require.NoError(t, err)
	pi := newTestIndexer(packageName)
	pi.addToIndex(packageName, pp)
	rep := NewRepresenter(pi)
	noImportsSpec := &TypeSpec{PackagePath: packageName, Name: "NoImports"}
	_, err = rep.GetRepresentation(noImportsSpec)
	require.NoError(t, err)
}
