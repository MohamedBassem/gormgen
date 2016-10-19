package gormgen

import (
	"math/rand"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/suite"
)

type BasicTestSuite struct {
	MainTestSuite
}

func TestBasicSuite(t *testing.T) {
	suite.Run(t, new(BasicTestSuite))
}

type BasicModel struct {
	ID   int `gorm:"primary_key"`
	Name string
	Age  int
}

func randomBasicModel() *BasicModel {
	return &BasicModel{
		Name: generateRandomString(10),
		Age:  rand.Intn(50),
	}
}

func (b *BasicTestSuite) TestSaveCreate() {
	model := randomBasicModel()
	err := model.Save(b.getDBConn())
	b.Require().Nil(err, "Saving shouldn't return an error")

	db := b.getDBConn()
	fetchedModel := &BasicModel{}
	notFound := db.Find(fetchedModel, randomBasicModel).RecordNotFound()
	b.Require().False(notFound, "The model should have been created")
	b.Require().Equal(model, fetchedModel, "The fetched model should have been correctly saved")
}

func (b *BasicTestSuite) TestSaveUpdate() {
	// Create and save a model
	model := randomBasicModel()
	b.getDBConn().Create(model)

	// Update and save the model
	model.Age = 0
	err := model.Save(b.getDBConn())
	b.Require().Nil(err, "Saving shouldn't return an error")

	// Fetch the model by its name and assert that its the same
	db := b.getDBConn()
	fetchedModel := []BasicModel{}
	db.Find(&fetchedModel, map[string]interface{}{"Name": model.Name})
	b.Require().Equal(1, len(fetchedModel), "The database shouldn't create a new model but rather update the old one")
	b.Require().Equal(model, &fetchedModel[0], "The fetched model should have been correctly updated")
}
