package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type ServerEnv struct {
	ENV             string `validate:"oneof=development test production"`
	CORS_WHITE_LIST []string

	APP_NAME string `validate:"min=1"`
	API_HOST string `validate:"min=1"`

	PORT string `validate:"number"`

	OPENOBSERVE_ENDPOINT   string `validate:"url"`
	OPENOBSERVE_CREDENTIAL string `validate:"min=1"`

	HEALER_POSTGRES_USER     string `validate:"min=1"`
	HEALER_POSTGRES_PASSWORD string `validate:"min=1"`
	HEALER_POSTGRES_DB       string `validate:"min=1"`
	HEALER_POSTGRES_HOST     string `validate:"min=1"`
	HEALER_POSTGRES_PORT     int    `validate:"min=1"`

	INDEXER_POSTGRES_USER     string `validate:"min=1"`
	INDEXER_POSTGRES_PASSWORD string `validate:"min=1"`
	INDEXER_POSTGRES_DB       string `validate:"min=1"`
	INDEXER_POSTGRES_HOST     string `validate:"min=1"`
	INDEXER_POSTGRES_PORT     int    `validate:"min=1"`

	MIGRATION_URL string `validate:"url"`

	CLIENTS_CONFIG_PATH string `validate:"min=1"`
	EVM_PRIVATE_KEY     string `validate:"min=1"`
	IS_TEST             bool
	IS_DEV              bool

	JWT_SECRET   string        `validate:"min=10"`
	JWT_DURATION time.Duration `validate:"min=0m"`
	AUTH_DOMAIN  string        `validate:"min=1"`
}

var Env *ServerEnv

func LoadEnvWithPath(relativePath string) {
	// Get the current working directory

	root, err := utils.RootPath()
	if err != nil {
		log.Fatalf("LoadEnvWithPath: Error getting root path: %s", err)
	}

	// Join the root path with the provided path
	fullPath := filepath.Join(root, relativePath)

	// Load the environment file
	err = godotenv.Load(os.ExpandEnv(fullPath))
	if err != nil {
		log.Fatalf("LoadEnvWithPath: Error loading %s file: %s", fullPath, err)
	}

	loadEnv()
}

func LoadEnv() {
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			log.Fatalln("LoadEnv: Error loading .env file: ", err)
		}
	} else if os.Getenv("ENV") == "test" {
		err := godotenv.Load(os.ExpandEnv(".env.test"))
		if err != nil {
			log.Fatalln("LoadEnv: Error loading .env.test file: ", err)
		}
	}

	loadEnv()
}

func loadEnv() {
	rawCORSWhiteList := os.Getenv("CORS_WHITE_LIST")
	var corsWhiteList []string
	if rawCORSWhiteList == "" {
		corsWhiteList = []string{
			"http://localhost:3000",
		}
	} else {
		corsWhiteList = strings.Split(rawCORSWhiteList, ",")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "12345"
	}

	env := os.Getenv("ENV")

	healerPgPort, err := strconv.Atoi(os.Getenv("HEALER_POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	indexerPgPort, err := strconv.Atoi(os.Getenv("INDEXER_POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	jwtDuration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		panic(err)
	}

	Env = &ServerEnv{
		ENV:             env,
		CORS_WHITE_LIST: corsWhiteList,

		APP_NAME: os.Getenv("APP_NAME"),
		API_HOST: os.Getenv("API_HOST"),
		PORT:     port,

		OPENOBSERVE_ENDPOINT:   os.Getenv("OPENOBSERVE_ENDPOINT"),
		OPENOBSERVE_CREDENTIAL: os.Getenv("OPENOBSERVE_CREDENTIAL"),

		HEALER_POSTGRES_USER:     os.Getenv("HEALER_POSTGRES_USER"),
		HEALER_POSTGRES_PASSWORD: os.Getenv("HEALER_POSTGRES_PASSWORD"),
		HEALER_POSTGRES_DB:       os.Getenv("HEALER_POSTGRES_DB"),
		HEALER_POSTGRES_HOST:     os.Getenv("HEALER_POSTGRES_HOST"),
		HEALER_POSTGRES_PORT:     healerPgPort,

		INDEXER_POSTGRES_USER:     os.Getenv("INDEXER_POSTGRES_USER"),
		INDEXER_POSTGRES_PASSWORD: os.Getenv("INDEXER_POSTGRES_PASSWORD"),
		INDEXER_POSTGRES_DB:       os.Getenv("INDEXER_POSTGRES_DB"),
		INDEXER_POSTGRES_HOST:     os.Getenv("INDEXER_POSTGRES_HOST"),
		INDEXER_POSTGRES_PORT:     indexerPgPort,

		MIGRATION_URL: os.Getenv("MIGRATION_URL"),

		CLIENTS_CONFIG_PATH: os.Getenv("CLIENTS_CONFIG_PATH"),
		EVM_PRIVATE_KEY:     os.Getenv("EVM_PRIVATE_KEY"),
		JWT_SECRET:          os.Getenv("JWT_SECRET"),
		JWT_DURATION:        jwtDuration,
		AUTH_DOMAIN:         os.Getenv("AUTH_DOMAIN"),

		IS_TEST: env == "test",
		IS_DEV:  env == "development",
	}

	validate := validator.New()
	err = validate.Struct(Env)
	if err != nil {
		panic(err)
	}

}
