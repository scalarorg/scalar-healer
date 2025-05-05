package bridge_test

import (
	"os"
	"testing"

	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/postgres"
)

var (
	testServer *server.Server
	db         *postgres.PostgresRepository
)

func setup() {
	config.LoadEnvWithPath("../../.env")
	config.Env.ENV = "test"
	config.Env.POSTGRES_DB = "scalar_healer_test"
	server := server.New()
	repo := server.DB.(*postgres.PostgresRepository)
	_ = repo

	// // Clean up test data
	// _, err := repo.DB.Collection("redeem_requests").DeleteMany(context.Background(), bson.M{})

	// if err != nil {
	// 	panic(err)
	// }

	// // Setup test data
	// gatewayAddr := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
	// _, err = repo.DB.Collection("gateway_addresses").InsertOne(context.Background(), bson.M{
	// 	"chain_id": uint64(1),
	// 	"address":  gatewayAddr.Bytes(),
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// err = repo.SaveTokenInfos(context.Background(), []models.Token{
	// 	{
	// 		Symbol:  "ETH",
	// 		ChainID: uint64(1),
	// 	},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// testServer = server
	// db = repo
}

func cleanupTestDB() {
	// Drop test collections
	// _, err := db.BridgeRequests.DeleteMany(context.Background(), bson.M{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.GatewayAddresses.DeleteMany(context.Background(), bson.M{})

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Tokens.DeleteMany(context.Background(), bson.M{})
	// if err != nil {
	// 	panic(err)
	// }

	// // Close connection
	// db.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	cleanupTestDB()
	os.Exit(code)
}
