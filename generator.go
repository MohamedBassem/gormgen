package main

import (
	"bytes"
	"fmt"
	"go/printer"
	"go/types"
	"strings"
)

type structConfig struct {
	StructName       string
	StructText       string
	QueryBuilderName string
}

type Generator struct {
	buf    *bytes.Buffer
	name   string
	parser *Parser
	config structConfig
}

func NewGenerator() *Generator {
	return &Generator{
		buf: new(bytes.Buffer),
	}
}

func (g *Generator) init() {
	g.generateImports()
}

func (g *Generator) buildConfig() {
	structName := strings.TrimSuffix(g.name, "Schema")
	structTextBuf := new(bytes.Buffer)
	printer.Fprint(structTextBuf, g.parser.fileSet, g.parser.GetTypeByName(g.name))
	g.config = structConfig{
		StructName:       structName,
		StructText:       structTextBuf.String(),
		QueryBuilderName: fmt.Sprintf("%sQueryBuilder", structName),
	}
}

func (g *Generator) generateImports() {
	importStatments.Execute(g.buf, nil)
}

func (g *Generator) generateMainStruct() {
	templateMainStruct.Execute(g.buf, g.config)
}

func (g *Generator) generateQueryBuilder() {
	templateQueryBuilder.Execute(g.buf, g.config)
}

func (g *Generator) generateFieldSpecificTemplates() {

	obj := g.parser.defs[g.parser.GetIdentByName(g.name)]
	if obj == nil {
		panic("SHOULDN'T HAPPEN")
	}

	fieldToColumn := make(map[string]string)
	structType := obj.Type().Underlying().(*types.Struct)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldToColumn[field.Name()] = field.Name()
	}

	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		fieldCnf := struct {
			FieldToColumn    map[string]string
			QueryBuilderName string
			FieldName        string
			FieldType        string
			Titelize         func(string) string
		}{
			FieldName:        f.Name(),
			FieldType:        f.Type().String(),
			FieldToColumn:    fieldToColumn,
			QueryBuilderName: g.config.QueryBuilderName,
			Titelize:         strings.Title,
		}
		templateWhereFunction.Execute(g.buf, fieldCnf)
		templateOrderByFunction.Execute(g.buf, fieldCnf)
	}
}

func (g *Generator) Generate(parser *Parser, name string) {
	g.name = name
	g.parser = parser
	g.buildConfig()
	g.generateMainStruct()
	g.generateQueryBuilder()
	g.generateFieldSpecificTemplates()
}

func (g *Generator) Flush() {
	fmt.Println(g.buf.String())
}
