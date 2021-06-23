package backend

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
	"strings"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type generatedCode struct {
	blocks []string
	// key=package path, value=alias
	imports map[string]string
}

func (gc *generatedCode) renderImportPairs() string {
	pairs := make([]string, 0)
	for k, v := range gc.imports {
		pairs = append(pairs, fmt.Sprintf("%s \"%s\"", v, k))
	}
	return strings.Join(pairs, "\n")
}

func (gc *generatedCode) renderBlocks() string {
	return strings.Join(gc.blocks, "\n")
}

func (gc *generatedCode) merge(right *generatedCode) {
	gc.blocks = append(gc.blocks, right.blocks...)
	if right.imports == nil {
		return
	}
	for k, v := range right.imports {
		// deduplicate imports and detect collisions
		// we should prevent collisions by normalizing import naming in a preprocessing pass
		if _, ok := gc.imports[k]; ok {
			continue
		}
		gc.imports[k] = v
	}
}

// Generator needs to be initialized with the package name,
// so use the new NewGenerator func for proper setup.
type Generator struct {
	gc []*generatedCode
	packageName string
}

func (g *Generator) Generate(vr types.ValRep) {
	mg := newMethodGenerator(vr)
	sizeSSZ := mg.GenerateSizeSSZ()
	if sizeSSZ != nil {
		g.gc = append(g.gc, sizeSSZ)
	}
}

var fileTemplate = `package {{.Package}}

{{ if .Imports -}}
import (
	{{.Imports}}
)
{{- end }}

{{.Blocks}}`

func (g *Generator) Render() ([]byte, error) {
	ft := template.New("generated.ssz.go")
	tmpl, err := ft.Parse(fileTemplate)
	if err != nil {
		return nil, err
	}
	final := &generatedCode{
		imports: make(map[string]string),
	}
	for _, gc := range g.gc {
		final.merge(gc)
	}
	pparts := strings.Split(g.packageName, "/")
	p := pparts[len(pparts)-1]
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, struct {
		Package string
		Imports string
		Blocks  string
	}{
		Package: p,
		Imports: final.renderImportPairs(),
		Blocks: final.renderBlocks(),
	})
	return format.Source(buf.Bytes())
}

type methodGenerator interface {
	GenerateSizeSSZ() *generatedCode
	variableSizeSSZ(fieldname string) string
	//GenerateMarshalSSZ() jen.Code
	//GenerateUnmarshalSSZ() jen.Code
	//GenerateHashTreeRoot() jen.Code
}

func newMethodGenerator(vr types.ValRep) methodGenerator {
	switch ty := vr.(type) {
	case *types.ValueBool:
		return &generateBool{ty}
	case *types.ValueByte:
		return &generateByte{ty}
	case *types.ValueContainer:
		return &generateContainer{ty}
	case *types.ValueList:
		return &generateList{ty}
	case *types.ValueOverlay:
		return &generateOverlay{ty}
	case *types.ValuePointer:
		return &generatePointer{ty}
	case *types.ValueUint:
		return &generateUint{ty}
	case *types.ValueUnion:
		return &generateUnion{ty}
	case *types.ValueVector:
		return &generateVector{ty}
	}
	panic(fmt.Sprintf("Cannot manage generation for unrecognized ValRep implementation %v", vr))
}
