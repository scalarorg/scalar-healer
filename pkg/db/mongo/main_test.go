package mongo_test

import (
	"os"
	"testing"

	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
)

var (
	repo *mongo.MongoRepository
)

func TestMain(m *testing.M) {
	config.LoadEnvWithPath("../../../.env")

	repo = mongo.NewMongoRepository()
	code := m.Run()
	cleanupTestDB()
	os.Exit(code)
}

func cleanupTestDB() {
	repo.Close()
}
