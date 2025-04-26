package mongo

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client           *mongo.Client
	DB               *mongo.Database
	GatewayAddresses *mongo.Collection
	Tokens           *mongo.Collection
	RedeemRequests   *mongo.Collection
}

var _ db.DbAdapter = (*MongoRepository)(nil)

func NewMongoRepository() *MongoRepository {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Env.MONGODB_URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		panic(err)
	}

	DB := client.Database(config.Env.MONGODB_DATABASE)

	m := &MongoRepository{
		client:           client,
		DB:               DB,
		GatewayAddresses: DB.Collection("gateway_addresses"),
		Tokens:           DB.Collection("tokens"),
		RedeemRequests:   DB.Collection("redeem_requests"),
	}

	m.initIndexes()

	log.Info().Msg("Connected to MongoDB!")

	return m
}

func (m *MongoRepository) Close() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	log.Info().Msg("Connection to MongoDB closed.")
}

func (m *MongoRepository) initIndexes() {
	var errs []error

	// TODO: Add indexes here

	if len(errs) > 0 {
		log.Logger.Fatal().Errs("Failed to create indexes", errs).Msg("")
	}
	log.Info().Msg("Indexes created")
}
