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
			
			return db.Save(t).Error
		}

		func (t *{{.StructName}}) Delete(db *gorm.DB) error {
			return db.Delete(t).Error
		}
	`)
	templateQueryBuilder = parseTemplateOrPanic(`
			type {{.QueryBuilderName}} struct {
				order []string
				where []struct{
					prefix string
					value interface{}
				}
				limit int
				offset int
			}

			func (qb *{{.QueryBuilderName}}) _buildQuery(db *gorm.DB) *gorm.DB {
				for _, where := range qb.where {
					ret = ret.Where(where)
				}
				for _, order := range qb.order {
					ret = ret.Order(order)
				}
				ret = ret.Limit(qb.limit).Offset(qb.offset)
				return ret
			}

			func (qb *{{.QueryBuilderName}}) Count(db *gorm.DB) (int, error) {
				var c int
				res := qb._buildQuery(db).Model(&{{.StructName}}{}).Count(&c)
				if res.RecordNotFound() {
					c = 0
				}
				return c, res.Error
			}

			func (qb *{{.QueryBuilderName}}) First(db *gorm.DB) (*{{.StructName}}, error) {
				ret := &{{.StructName}}{}
				res := qb._buildQuery(db).First(ret)
				if res.RecordNotFound() {
					ret = nil
				}
				return ret, res.Error
			}

			func (qb *{{.QueryBuilderName}}) QueryOne(db *gorm.DB) (*{{.StructName}}, error) {
				qb.offset = 1
				ret, err := qb.QueryAll(db)
				if len(ret) > 0 {
					return &ret[0], err
				}else{
					return nil, err
				}
			}

			func (qb *{{.QueryBuilderName}}) QueryAll(db *gorm.DB) ([]{{.StructName}}, error) {
				ret := []{{.StructName}}{}
				err := qb._buildQuery(db).Find(&ret).Error
				return ret, err
			}

			func (qb *{{.QueryBuilderName}}) Limit(limit int) *{{.QueryBuilderName}} {
				qb.limit = limit
				return qb
			}

			func (qb *{{.QueryBuilderName}}) Offset(offset int) *{{.QueryBuilderName}} {
				qb.offset = offset
				return qb
			}
		`)
	templateWhereFunction = parseTemplateOrPanic(`
			func (qb *{{.QueryBuilderName}}) Where{{call .Titelize .FieldName}}(p gormgen.Predict, value {{.FieldType}}) *{{.QueryBuilderName}} {
				 qb.where = append(qb.where, struct {
					prefix string
					value interface{}
				}{
					{{index .FieldToColumn .FieldName}} + p.String(),
					value,
				})
				return qb
			}
		`)
	templateOrderByFunction = parseTemplateOrPanic(`
			func (qb *{{.QueryBuilderName}}) OrderBy{{call .Titelize .FieldName}}(asc bool) *{{.QueryBuilderName}} {
				order := "DESC"
				if asc {
					order = "ASC"
				}

				qb.order = append(qb.order, {{index .FieldToColumn .FieldName}} + order)
				return qb
			}
		`)
)
