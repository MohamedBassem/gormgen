package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

type structConfig struct {
	StructName       string
	QueryBuilderName string
}

type Generator struct {
	buf       *bytes.Buffer
	name      string
	structAST *ast.StructType
	config    structConfig
}

func NewGenerator() *Generator {
	return &Generator{
		buf: new(bytes.Buffer),
	}
}

func (g *Generator) buildConfig() {
	g.config = structConfig{
		StructName: strings.TrimSuffix(g.name, "Schema"),
	}
}

func (g *Generator) generateImports() {
	importStatments.Execute(g.buf, nil)
}

func (g *Generator) generateQueryBuilder() {
}

func (g *Generator) Generate(name string, ast *ast.StructType) {
	g.name = name
	g.structAST = ast
	g.generateImports()
	fmt.Println(g.buf.String())
}
