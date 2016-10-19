package gormgen

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

type Parser struct {
	dir         string
	pkgName     string
	types       map[*ast.Ident]*ast.StructType
	files       []string
	parsedFiles []*ast.File
	fileSet     *token.FileSet
	defs        map[*ast.Ident]types.Object
}

func NewParser() *Parser {
	return &Parser{
		types: make(map[*ast.Ident]*ast.StructType),
	}
}

func (p *Parser) getFiles() {
	pkg, err := build.Default.ImportDir(p.dir, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", p.dir, err)
	}
	var files []string
	p.pkgName = pkg.Name
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	p.files = files
}

func (p *Parser) parseFiles() {
	var parsedFiles []*ast.File
	fs := token.NewFileSet()
	for _, file := range p.files {
		parsedFile, err := parser.ParseFile(fs, file, nil, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s\n", file, err)
		}
		parsedFiles = append(parsedFiles, parsedFile)
	}
	p.parsedFiles, p.fileSet = parsedFiles, fs
}

func (p *Parser) parseTypes(file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// We only care about struct declaration
			var structType *ast.StructType
			if structType, ok = typeSpec.Type.(*ast.StructType); !ok {
				continue
			}
			p.types[typeSpec.Name] = structType
		}
		return true
	})
}

// Copied from the stringer package. Check the README
func (p *Parser) typeCheck() {
	p.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: p.defs,
	}
	_, err := config.Check(p.dir, p.fileSet, p.parsedFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
}

func (p *Parser) ParseDir(dir string) {
	p.dir = dir
	p.getFiles()
	p.parseFiles()
	for _, f := range p.parsedFiles {
		p.parseTypes(f)
	}
	p.typeCheck()
}

func (p *Parser) GetTypeByName(name string) *types.Struct {
	ident := p.GetIdentByName(name)
	if ident == nil {
		return nil
	}
	def, ok := p.defs[ident]
	if !ok {
		return nil
	}
	structType, ok := def.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	return structType
}

func (p *Parser) GetIdentByName(name string) *ast.Ident {
	for id := range p.types {
		if id.Name == name {
			return id
		}
	}
	return nil
}
