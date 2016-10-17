package main

import (
	"fmt"
	"html/template"
)

var templateId = 0

func parseTemplateOrPanic(t string) *template.Template {
	templateIdStr := fmt.Sprintf("template_%v", templateId)
	templateId++
	tpl, err := template.New(templateIdStr).Parse(t)
	if err != nil {
		panic(err)
	}
	return tpl
}

var (
	importStatments = parseTemplateOrPanic(`
			import "github.com/MohamedBassem/gormgen"
			import "github.com/jinzhu/gorm"
		`)
	templateQueryBuilder = parseTemplateOrPanic(`
			type {{.StructName}}QueryBuilder struct {
				order []string 
			}

			func (bd *{{.StructName}}QueryBuilder) 
		`)
	templateWhereFunction = parseTemplateOrPanic(`
			func Where{{.FieldName}}(p gormgen.Predict, value {{.FieldType}}) {
				
			}
		`)
)
