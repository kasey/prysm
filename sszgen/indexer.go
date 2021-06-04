package sszgen

import (
	"fmt"
	"strconv"
)

type PackageIndex struct {
	sourcePackage string
	index map[string]PackageParser
	structCache map[[2]string]*TypeSpec
}

func buildIndex(packagePath string, pi *PackageIndex) error {
	pparser, err := NewPackageParser(packagePath)
	if err != nil {
		return err
	}
	pi.addToIndex(packagePath, pparser)

	imports, err := pparser.Imports()
	if err != nil {
		return err
	}
	for _, pkg := range imports {
		pkgStr, err := strconv.Unquote(pkg.Path.Value)
		if err != nil {
			return err
		}
		if _, indexed := pi.index[pkgStr]; !indexed {
			err = buildIndex(pkgStr, pi)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func BuildPackageIndex(sourcePackage string) (*PackageIndex, error) {
	pi := &PackageIndex{
		sourcePackage: sourcePackage,
		index: make(map[string]PackageParser),
		structCache: make(map[[2]string]*TypeSpec),
	}
	err := buildIndex(sourcePackage, pi)
	return pi, err
}

func (pi *PackageIndex) addToIndex(packagePath string, pp PackageParser) {
	pi.index[packagePath] = pp
}

func (pi *PackageIndex) getParser(packagePath string) (PackageParser, error) {
	pkg, ok := pi.index[packagePath]
	if !ok {
		return nil, fmt.Errorf("Could not find package named %s", packagePath)
	}
	return pkg, nil
}

func (pi *PackageIndex) PackageTypes(packagePath string) ([]*TypeSpec, error) {
	pkg, err := pi.getParser(packagePath)
	if err != nil {
		return nil, err
	}
	return pkg.AllStructs(), nil
}

func (pi *PackageIndex) GetStruct(packagePath, typeName string) (*TypeSpec, error) {
	cached, ok := pi.structCache[[2]string{packagePath,typeName}]
	if ok {
		return cached, nil
	}
	pkg, err := pi.getParser(packagePath)
	if err != nil {
		return nil, err
	}
	ts, err := pkg.GetStruct(typeName)
	if err != nil {
		return nil, err
	}
	pi.structCache[[2]string{packagePath,typeName}]	= ts
	return ts, nil
}
