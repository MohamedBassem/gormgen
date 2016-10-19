package gormgen

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type ComplexTestSuite struct {
	MainTestSuite
}

func TestComplexSuite(t *testing.T) {
	suite.Run(t, new(ComplexTestSuite))
}

type ComplexModel struct {
	gorm.Model
	Name string
}
