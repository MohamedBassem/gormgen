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

	templateMainStruct = parseTemplateOrPanic(`
		type {{.StructName}} {{.StructText}}

		func (t *{{.StructName}}) Save(db *gorm.DB) error {

		}

		func (t *{{.StructName}}) Delete(db *gorm.DB) error {

		}
	`)
	templateQueryBuilder = parseTemplateOrPanic(`
			type {{.QueryBuilderName}} struct {
				order []string
				where []string
				limit int
				offset int
			}

			func (qb *{{.QueryBuilderName}}) Count(db *gorm.DB) (int, error) {

			}

			func (qb *{{.QueryBuilderName}}) QueryOne(db *gorm.DB) (*{{.StructName}}, error) {

			}

			func (qb *{{.QueryBuilderName}}) QueryAll(db *gorm.DB) ([]{{.StructName}}, error) {

			}

			func (qb *{{.QueryBuilderName}}) Limit(limit int) *{{.QueryBuilderName}} {
				qb.limit = limit
			}

			func (qb *{{.QueryBuilderName}}) Offset(offset int) *{{.QueryBuilderName}} {
				qb.offset = offset
			}
		`)
	templateWhereFunction = parseTemplateOrPanic(`
			func Where{{.FieldName}}(p gormgen.Predict, value {{.FieldType}}) *{{.QueryBuilderName}} {

			}
		`)
	templateOrderByFunction = parseTemplateOrPanic(`
			func OrderBy{{.FieldName}}(asc bool) *{{.QueryBuilderName}} {

			}
		`)
)
