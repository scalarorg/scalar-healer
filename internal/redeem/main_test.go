package redeem_test

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func setup(t *testing.T) *server.Server {

	os.Setenv("MONGODB_DATABASE", "scalar_test")
	config.LoadEnvWithPath("../../.env")

	testServer := server.New()

	repo := testServer.DB.(*mongo.MongoRepository)

	// Clean up test data
	_, err := repo.DB.Collection("redeem_requests").DeleteMany(context.Background(), bson.M{})

	require.NoError(t, err)

	// Setup test data
	gatewayAddr := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
	_, err = repo.DB.Collection("gateway_addresses").InsertOne(context.Background(), bson.M{
		"chain_id": uint64(1),
		"address":  gatewayAddr.Bytes(),
	})
	require.NoError(t, err)

	_, err = repo.DB.Collection("tokens").InsertOne(context.Background(), bson.M{
		"symbol": "ETH",
		"active": true,
	})
	require.NoError(t, err)

	return testServer
}

func cleanupTestDB(t *testing.T, db *mongo.MongoRepository) {
	// Drop test collections
	_, err := db.DB.Collection("redeem_requests").DeleteMany(context.Background(), bson.M{})

	require.NoError(t, err)
	_, err = db.DB.Collection("gateway_addresses").DeleteMany(context.Background(), bson.M{})

	require.NoError(t, err)
	_, err = db.DB.Collection("tokens").DeleteMany(context.Background(), bson.M{})

	require.NoError(t, err)

	// Close connection
	db.Close()
}
