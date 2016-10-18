package main

import (
	"bytes"
	"fmt"
	"go/printer"
	"go/types"
	"strings"

	"github.com/jinzhu/gorm"
)

type fieldConfig struct {
	FieldToColumn    map[string]string
	QueryBuilderName string
	FieldName        string
	FieldType        string
	Titelize         func(string) string
}

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

func (g *Generator) buildFieldConfig() []fieldConfig {
	obj := g.parser.defs[g.parser.GetIdentByName(g.name)]
	if obj == nil {
		panic("SHOULDN'T HAPPEN")
	}
	fieldToColumn := make(map[string]string)
	fields := []*types.Var{}
	structType := obj.Type().Underlying().(*types.Struct)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}
		tag := g.parser.GetFieldTag(g.name, field.Name())
		fieldToColumn[field.Name()] = gorm.ToDBName(field.Name())
		if tag != nil {
			if gormt, ok := tag.Lookup("gorm"); ok {
				if gormt == "-" {
					continue
				}
				parts := strings.Split(gormt, ";")
				for _, part := range parts {
					kv := strings.Split(part, ":")
					if len(kv) > 1 && kv[0] == "column" {
						fieldToColumn[field.Name()] = kv[1]
					}
				}
			}
		}
		fields = append(fields, field)
	}

	ret := []fieldConfig{}

	for _, f := range fields {
		ret = append(ret, fieldConfig{
			FieldName:        f.Name(),
			FieldType:        types.TypeString(f.Type(), func(p *types.Package) string { return p.Name() }),
			FieldToColumn:    fieldToColumn,
			QueryBuilderName: g.config.QueryBuilderName,
			Titelize:         strings.Title,
		})
	}
	return ret
}

func (g *Generator) generateFieldSpecificTemplates() {
	for _, f := range g.buildFieldConfig() {
		templateWhereFunction.Execute(g.buf, f)
		templateOrderByFunction.Execute(g.buf, f)
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
