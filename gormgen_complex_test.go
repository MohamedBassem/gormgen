package gormgen

import (
	"reflect"
	"testing"

	"github.com/MohamedBassem/gormgen/internal/tmp"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ComplexTestSuite struct {
	MainTestSuite
}

func TestComplexSuite(t *testing.T) {
	suite.Run(t, new(ComplexTestSuite))
}

type EmbeddedStruct struct {
	EmbeddedName        string
	notExportedEmbedded string
	GormIgnoredEmbedded string `gorm:"-"`
}

type EmbeddedStruct2 struct {
	MustNotBeThere string
}

type RelationStruct struct{}

type M2MStruct struct {
	gorm.Model
	Name         string
	ComplexModel []*ComplexModel `gorm:"many2many:m2m_rel;"`
}

type ComplexModel struct {
	gorm.Model
	Name        string
	notExported int
	GormIgnored string `gorm:"-"`
	EmbeddedStruct
	EmbeddedStruct2 `gorm:"-"`
	Test            struct {
		NoIdea string
	} `gorm:"embedded"`
	Relation             RelationStruct           `gorm:"-"` // Should be ignored for now
	AnotherPackageStruct tmp.AnotherPackageStruct `gorm:"-"` // Should be ignored for now
	M2MStruct            []*M2MStruct             `gorm:"many2many:m2m_rel;"`
}

func (c *ComplexTestSuite) TestNormalField() {
	_, ok := reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereName")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByName")
	c.Assert().True(ok, "The method should be created")
}

func (c *ComplexTestSuite) TestFieldsOfEmbeddedStruct() {

	_, ok := reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereID")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereCreatedAt")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereUpdatedAt")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereDeletedAt")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByID")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByCreatedAt")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByUpdatedAt")
	c.Assert().True(ok, "The method should be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByDeletedAt")
	c.Assert().True(ok, "The method should be created")
}

func (c *ComplexTestSuite) TestIgnoreGormIgnored() {
	_, ok := reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereGormIgnored")
	c.Assert().False(ok, "The method is ignored by struct tag, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereGormIgnoredEmbedded")
	c.Assert().False(ok, "The method is ignored by struct tag, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByGormIgnored")
	c.Assert().False(ok, "The method is ignored by struct tag, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByGormIgnoredEmbedded")
	c.Assert().False(ok, "The method is ignored by struct tag, it shouldn't be created")
}

func (c *ComplexTestSuite) TestIgnoreUnexportedFields() {
	_, ok := reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereNotExported")
	c.Assert().False(ok, "The method's field is not exported, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereNotExportedEmbedded")
	c.Assert().False(ok, "The method's field is not exported, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByNotExported")
	c.Assert().False(ok, "The method's field is not exported, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByExportedEmbedded")
	c.Assert().False(ok, "The method's field is not exported, it shouldn't be created")
}

func (c *ComplexTestSuite) TestIgnoreStructsByTag() {
	_, ok := reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("WhereMustNotBeThere")
	c.Assert().False(ok, "The method's struct is ignored, it shouldn't be created")

	_, ok = reflect.TypeOf(&ComplexModelQueryBuilder{}).MethodByName("OrderByMustNotBeThere")
	c.Assert().False(ok, "The method's struct is ignored, it shouldn't be created")
}
