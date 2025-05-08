package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	POSTGRES_USER     string `validate:"min=1"`
	POSTGRES_PASSWORD string `validate:"min=1"`
	POSTGRES_DB       string `validate:"min=1"`
	POSTGRES_HOST     string `validate:"min=1"`
	POSTGRES_PORT     int    `validate:"min=1"`

	MIGRATION_URL string `validate:"url"`

	CLIENTS_CONFIG_PATH string `validate:"min=1"`
	EVM_PRIVATE_KEY     string `validate:"min=1"`
	IS_TEST             bool
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

	postgresPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
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

		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     postgresPort,
		MIGRATION_URL:     os.Getenv("MIGRATION_URL"),

		CLIENTS_CONFIG_PATH: os.Getenv("CLIENTS_CONFIG_PATH"),
		EVM_PRIVATE_KEY:     os.Getenv("EVM_PRIVATE_KEY"),
		IS_TEST:             env == "test",
	}

	validate := validator.New()
	err = validate.Struct(Env)
	if err != nil {
		panic(err)
	}

}
