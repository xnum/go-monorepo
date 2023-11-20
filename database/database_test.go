package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-monorepo/database"
	"go-monorepo/logging"
	"go-monorepo/models"
)

func testDatabase(t *testing.T, opt database.ConnectOption) {
	assert := assert.New(t)
	logging.TestingInitialize()
	defer logging.TestingFinalize()

	database.TestingInitialize(database.Default, opt)
	database.AutoMigrate(database.Default, models.Models())
	defer database.TestingFinalize()

	conn := database.GetDB(database.Default)

	type Animal struct {
		ID   int64
		Name string
		Age  int64
	}

	err := conn.Migrator().CreateTable(Animal{})
	assert.NoError(err)

	animal := Animal{ID: 1, Name: "Bear", Age: 33}
	err = conn.Create(&animal).Error
	assert.NoError(err)

	var bear Animal
	err = conn.First(&bear).Error
	assert.NoError(err)

	assert.Equal(animal, bear)
}

func TestDatabase(t *testing.T) {
	testDatabase(t, database.PostgresOpt)
}
