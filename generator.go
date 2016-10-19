package gormgen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"io/ioutil"
	"strings"

	"github.com/jinzhu/gorm"
)

type structHelpers struct {
	Titelize func(string) string
}

type fieldConfig struct {
	FieldName  string
	ColumnName string
	FieldType  string
	Titelize   func(string) string
}

type structConfig struct {
	StructName       string
	QueryBuilderName string
	Fields           []fieldConfig
	Helpers          structHelpers
}

type structsConfig struct {
	PkgName string
	Structs []structConfig
}

type Generator struct {
	buf        *bytes.Buffer
	outputFile string
	config     structsConfig
}

func NewGenerator(outputFile string) *Generator {
	return &Generator{
		buf:        new(bytes.Buffer),
		outputFile: outputFile,
	}
}

func (g *Generator) Init(parser *Parser, structs []string) error {
	if err := g.validateStructs(parser, structs); err != nil {
		return err
	}
	g.config.PkgName = parser.pkgName
	for _, st := range structs {
		g.config.Structs = append(g.config.Structs, *g.buildConfig(parser, st))
	}
	return nil
}

func (g *Generator) validateStructs(parser *Parser, structs []string) error {
	for _, st := range structs {
		if parser.GetTypeByName(st) == nil {
			return fmt.Errorf("Type %v is not found", st)
		}
	}
	return nil
}

func (g *Generator) buildConfig(parser *Parser, structName string) *structConfig {
	cnf := &structConfig{
		StructName:       structName,
		QueryBuilderName: fmt.Sprintf("%sQueryBuilder", structName),
		Helpers: structHelpers{
			Titelize: strings.Title,
		},
	}
	cnf.Fields = g.buildFieldConfig(parser, structName)
	return cnf
}

func (g *Generator) buildFieldConfig(parser *Parser, structName string) []fieldConfig {
	obj := parser.defs[parser.GetIdentByName(structName)]
	fieldToColumn := make(map[string]string)
	fields := []*types.Var{}
	structType := obj.Type().Underlying().(*types.Struct)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}
		tag := parser.GetFieldTag(structName, field.Name())
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
			FieldName:  f.Name(),
			FieldType:  types.TypeString(f.Type(), func(p *types.Package) string { return p.Name() }),
			ColumnName: fieldToColumn[f.Name()],
		})
	}
	return ret
}

func (g *Generator) Generate() error {
	return tpl.Execute(g.buf, g.config)
}

func (g *Generator) Format() error {
	formatedOutput, err := format.Source(g.buf.Bytes())
	if err != nil {
		return err
	}
	g.buf = bytes.NewBuffer(formatedOutput)
	return nil
}

func (g *Generator) Flush() {
	ioutil.WriteFile(g.outputFile, g.buf.Bytes(), 0644)
}
