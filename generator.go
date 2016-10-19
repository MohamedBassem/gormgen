package gormgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"io/ioutil"
	"reflect"
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
}

type structConfig struct {
	StructName       string
	QueryBuilderName string
	Fields           []fieldConfig
}

type structsConfig struct {
	PkgName string
	Helpers structHelpers
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
	g.config = structsConfig{
		PkgName: parser.pkgName,
		Helpers: structHelpers{
			Titelize: strings.Title,
		},
	}
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
	}
	cnf.Fields = g.buildFieldConfig(parser, structName)
	return cnf
}

func (g *Generator) buildFieldConfig(parser *Parser, structName string) []fieldConfig {
	fields := []fieldConfig{}
	structType := parser.GetTypeByName(structName)
	for i := 0; i < structType.Fields.NumFields(); i++ {
		field := structType.Fields.List[i]
		for _, name := range field.Names {
			if !ast.IsExported(name.Name) {
				continue
			}
			columnName := gorm.ToDBName(name.Name)
			if field.Tag != nil {
				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
				if gormt, ok := tag.Lookup("gorm"); ok {
					if gormt == "-" {
						continue
					}
					parts := strings.Split(gormt, ";")
					for _, part := range parts {
						kv := strings.Split(part, ":")
						if len(kv) > 1 && kv[0] == "column" {
							columnName = kv[1]
						}
					}
				}
			}
			// Get field type
			tmpBuf := &bytes.Buffer{}
			printer.Fprint(tmpBuf, parser.fileSet, field.Type)
			fields = append(fields, fieldConfig{
				FieldName:  name.Name,
				ColumnName: columnName,
				FieldType:  tmpBuf.String(),
			})
		}
	}
	return fields
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

func (g *Generator) Flush() error {
	return ioutil.WriteFile(g.outputFile, g.buf.Bytes(), 0644)
}
