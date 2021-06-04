package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prysmaticlabs/prysm/sszgen"
	"github.com/urfave/cli/v2"
)

var sourcePackage, output, typeNames string
var generate = &cli.Command{
	Name:    "generate",
	ArgsUsage: "<input package, eg github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1>",
	Aliases: []string{"gen"},
	Usage:   "generate methodsets for a go struct type to support ssz ser/des",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Value:       "",
			Usage:       "directory to write generated code (same as input by default)",
			Destination: &output,
		},
		&cli.StringFlag{
			Name:        "type-names",
			Value:       "",
			Usage:       "if specified, only generate methods for types specified in this comma-separated list",
			Destination: &typeNames,
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() > 0 {
			sourcePackage = c.Args().Get(0)
		}
		index, err := sszgen.BuildPackageIndex(sourcePackage)
		if err != nil {
			return err
		}
		rep := sszgen.NewRepresenter(index)
		if err != nil {
			return err
		}

		var specs []*sszgen.TypeSpec
		if len(typeNames) > 0 {
			for _, n := range strings.Split(strings.TrimSpace(typeNames), ",") {
				specs = append(specs, &sszgen.TypeSpec{PackagePath: sourcePackage, Name: n})
			}
		} else {
			specs, err = index.PackageTypes(sourcePackage)
			if err != nil {
				return err
			}
		}
		if len(specs) == 0 {
			return fmt.Errorf("Could not find any codegen targets in source package %s", sourcePackage)
		}

		if output == "" {
			output = "generated.ssz.go"
		}
		outFh, err := os.Create(output)
		defer outFh.Close()
		if err != nil {
			return err
		}
		mf := &sszgen.MergedFile{}

		for _, s := range specs {
			typeRep, err := rep.GetRepresentation(s)
			if err != nil {
				return err
			}
			rendered, err := sszgen.Render(typeRep)
			if err != nil {
				return err
			}
			mf.Accumulate(rendered)
		}
		merged, err := mf.Merge()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFh, merged)
		return err
	},
}