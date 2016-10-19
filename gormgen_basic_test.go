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
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"column:name2"`
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
	db.Find(&fetchedModel, map[string]interface{}{"name2": model.Name})
	b.Require().Equal(1, len(fetchedModel), "The database shouldn't create a new model but rather update the old one")
	b.Require().Equal(model, &fetchedModel[0], "The fetched model should have been correctly updated")
}

func (b *BasicTestSuite) TestQueryCount() {
	// Create 10 models
	for i := 0; i < 10; i++ {
		b.getDBConn().Create(randomBasicModel())
	}

	// Query the whole table for the number of entries
	c, err := (&BasicModelQueryBuilder{}).Count(b.getDBConn())
	b.Require().Nil(err, "The count function shouldn't return an error")
	b.Require().Equal(10, c, "The count function should report 10 entries")
}

func (b *BasicTestSuite) TestQueryAll() {
	models := []BasicModel{}
	for i := 0; i < 10; i++ {
		models = append(models, *randomBasicModel())
		b.getDBConn().Create(&models[i])
	}

	// Query the whole table for the entries
	fetched, err := (&BasicModelQueryBuilder{}).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The QueryAll function shouldn't return an error")
	b.Require().Equal(models, fetched, "The query all function should return all 10 entries")
}

func (b *BasicTestSuite) TestQueryLimit() {
	models := []BasicModel{}
	for i := 0; i < 10; i++ {
		models = append(models, *randomBasicModel())
		b.getDBConn().Create(&models[i])
	}

	fetched, err := (&BasicModelQueryBuilder{}).Limit(1).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query shouldn't return an error")
	b.Require().Equal(1, len(fetched), "The query should return only one item")
	b.Require().Equal(models[0], fetched[0], "The query one function should return only the first element")
}

func (b *BasicTestSuite) TestQueryOffset() {
	models := []BasicModel{}
	for i := 0; i < 10; i++ {
		models = append(models, *randomBasicModel())
		b.getDBConn().Create(&models[i])
	}

	// NOTE, I added LIMIT a limit because SQLite doesn't understand an offset without a limit and gorm
	// doesn't support negative limits. https://github.com/jinzhu/gorm/issues/1045
	fetched, err := (&BasicModelQueryBuilder{}).Limit(10).Offset(1).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query shouldn't return an error")
	b.Require().Equal(9, len(fetched), "The query should return all but the first item")
	b.Require().Equal(models[1:], fetched, "The query one function should return only the first element")
}

func (b *BasicTestSuite) TestQueryOne() {
	models := []*BasicModel{}
	for i := 0; i < 10; i++ {
		models = append(models, randomBasicModel())
		b.getDBConn().Create(models[i])
	}

	// Query the whole table for the number of entries
	fetched, err := (&BasicModelQueryBuilder{}).QueryOne(b.getDBConn())
	b.Require().Nil(err, "The QueryOne function shouldn't return an error")
	b.Require().Equal(models[0], fetched, "The query one function should return only the first element")
}

func (b *BasicTestSuite) TestQueryWhere() {
	models := []BasicModel{
		{
			Name: "Test1",
			Age:  10,
		},
		{
			Name: "Test2",
			Age:  20,
		},
		{
			Name: "Test3",
			Age:  30,
		},
		{
			Name: "Test4",
			Age:  50,
		},
	}
	for i := 0; i < len(models); i++ {
		b.getDBConn().Create(&models[i])
	}

	fetched, err := (&BasicModelQueryBuilder{}).WhereID(EqualPredict, 1).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(1, len(fetched), "The query function return a single element")
	b.Require().Equal(models[0], fetched[0], "The query should return only the first element")

	fetched, err = (&BasicModelQueryBuilder{}).WhereName(EqualPredict, "Test2").QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(1, len(fetched), "The query function return a single element")
	b.Require().Equal(models[1], fetched[0], "The query should return only the first element")

	fetched, err = (&BasicModelQueryBuilder{}).WhereAge(GreaterThanOrEqualPredict, 30).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(2, len(fetched), "The query function return two element")
	b.Require().Equal(models[2:], fetched, "The query should return only elements with age creater than 30")

	fetched, err = (&BasicModelQueryBuilder{}).
		WhereName(SmallerThanPredict, "Test4").
		WhereAge(GreaterThanOrEqualPredict, 30).
		QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(1, len(fetched), "The query function return one element")
	b.Require().Equal(models[2], fetched[0], "The query should return only the third element")

	fetched, err = (&BasicModelQueryBuilder{}).WhereAge(EqualPredict, -100).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(0, len(fetched), "The query function return two element")
	b.Require().Equal([]BasicModel{}, fetched, "The query should return only elements with age creater than 30")
}

func (b *BasicTestSuite) TestQueryOrder() {
	models := []BasicModel{
		{
			Name: "Test1",
			Age:  10,
		},
		{
			Name: "Test2",
			Age:  10,
		},
		{
			Name: "Test3",
			Age:  30,
		},
		{
			Name: "Test4",
			Age:  50,
		},
	}
	for i := 0; i < len(models); i++ {
		b.getDBConn().Create(&models[i])
	}

	fetched, err := (&BasicModelQueryBuilder{}).OrderByID(true).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(len(models), len(fetched), "The query function return all the items")
	b.Require().Equal(models, fetched, "The query should return the items sorted by ID")

	fetched, err = (&BasicModelQueryBuilder{}).OrderByID(false).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(len(models), len(fetched), "The query function return all the items")
	b.Require().Equal(reverseBasicModelSlice(models), fetched, "The query should return the items sorted by ID")

	fetched, err = (&BasicModelQueryBuilder{}).OrderByAge(false).OrderByName(true).QueryAll(b.getDBConn())
	b.Require().Nil(err, "The query function shouldn't return an error")
	b.Require().Equal(len(models), len(fetched), "The query function return all the items")
	b.Require().Equal([]BasicModel{models[3], models[2], models[0], models[1]}, fetched, "The query should return the items sorted by ID")
}
