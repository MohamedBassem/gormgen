package gormgen

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

//go:generate gormgen -structs BasicModel,ComplexModel -output structs_test.go
var modelRegistry = []interface{}{
	&BasicModel{},
	&ComplexModel{},
}

type MainTestSuite struct {
	suite.Suite
	DB             *gorm.DB
	testDBFileName string
}

func (m *MainTestSuite) migrateModels() {
	db := m.getDBConn()
	for _, m := range modelRegistry {
		err := db.AutoMigrate(m)
		if err != nil {
			panic(err)
		}
	}
}

func (m *MainTestSuite) SetupTest() {
	f, err := ioutil.TempFile("", "test.db")
	if err != nil {
		log.Fatalf("Couldn't create temp sqllite database: %v", err.Error())
	}
	m.testDBFileName = f.Name()
	m.migrateModels()
}

func (m *MainTestSuite) TearDownTest() {
	m.DB = nil
	err := os.Remove(m.testDBFileName)
	if err != nil {
		log.Fatalf("Couldn't delete temp sqllite database: %v", err.Error())
	}
}

func (m *MainTestSuite) getDBConn() *gorm.DB {
	if m.DB == nil {
		db, err := gorm.Open(sqlite.Open(m.testDBFileName), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		m.DB = db
	}
	return m.DB
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}
