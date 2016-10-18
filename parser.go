package main

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strings"
)

type Parser struct {
	dir         string
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
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
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

func (p *Parser) typeCheck() {
	// Copied from the stringer library
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

func (p *Parser) ParseDir(dir string) {
	p.dir = dir
	p.getFiles()
	p.parseFiles()
	p.typeCheck()
	for _, f := range p.parsedFiles {
		p.parseTypes(f)
	}
}

func (p *Parser) GetTypeByName(name string) *ast.StructType {
	for id, v := range p.types {
		if id.Name == name {
			return v
		}
	}
	return nil
}

func (p *Parser) GetIdentByName(name string) *ast.Ident {
	for id := range p.types {
		if id.Name == name {
			return id
		}
	}
	return nil
}

func (p *Parser) GetFieldTag(structName, fieldName string) *reflect.StructTag {
	t := p.GetTypeByName(structName)
	if t == nil {
		return nil
	}
	for _, f := range t.Fields.List {
		for _, id := range f.Names {
			if id.Name == fieldName && f.Tag != nil {
				stag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
				return &stag
			}
		}
	}
	return nil
}
