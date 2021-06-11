package sszgen

type PackageIndex struct {
	sourcePackage string
	index map[string]PackageParser
	structCache map[[2]string]*TypeSpec
}

func NewPackageIndex() *PackageIndex {
	pi := &PackageIndex{
		index: make(map[string]PackageParser),
		structCache: make(map[[2]string]*TypeSpec),
	}
	return pi
}

func (pi *PackageIndex) getParser(packagePath string) (PackageParser, error) {
	pkg, ok := pi.index[packagePath]
	if ok {
		return pkg, nil
	}
	pkg, err := NewPackageParser(packagePath)
	if err == nil {
		pi.index[packagePath] = pkg
	}
	return pkg, err
}

func (pi *PackageIndex) PackageTypes(packagePath string) ([]*TypeSpec, error) {
	pkg, err := pi.getParser(packagePath)
	if err != nil {
		return nil, err
	}
	return pkg.AllTypes(), nil
}

func (pi *PackageIndex) GetType(packagePath, typeName string) (*TypeSpec, error) {
	cached, ok := pi.structCache[[2]string{packagePath,typeName}]
	if ok {
		return cached, nil
	}
	pkg, err := pi.getParser(packagePath)
	if err != nil {
		return nil, err
	}
	ts, err := pkg.GetType(typeName)
	if err != nil {
		return nil, err
	}
	pi.structCache[[2]string{packagePath,typeName}]	= ts
	return ts, nil
}
