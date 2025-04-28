package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type ServerEnv struct {
	ENV             string `validate:"oneof=development test production"`
	CORS_WHITE_LIST []string

	APP_NAME string `validate:"min=1"`
	API_HOST string `validate:"min=1"`

	PORT string `validate:"number"`

	OPENOBSERVE_ENDPOINT   string `validate:"url"`
	OPENOBSERVE_CREDENTIAL string `validate:"min=1"`

	MONGODB_URI      string `validate:"uri"`
	MONGODB_DATABASE string `validate:"min=1"`

	CLIENTS_CONFIG_PATH string `validate:"min=1"`
	EVM_PRIVATE_KEY     string `validate:"min=1"`
}

var Env *ServerEnv

func LoadEnvWithPath(path string) {
	err := godotenv.Load(os.ExpandEnv(path))
	if err != nil {
		log.Fatalf("LoadEnvWithPath: Error loading %s file: %s", path, err)
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

	Env = &ServerEnv{
		ENV:             env,
		CORS_WHITE_LIST: corsWhiteList,

		APP_NAME: os.Getenv("APP_NAME"),
		API_HOST: os.Getenv("API_HOST"),
		PORT:     port,

		OPENOBSERVE_ENDPOINT:   os.Getenv("OPENOBSERVE_ENDPOINT"),
		OPENOBSERVE_CREDENTIAL: os.Getenv("OPENOBSERVE_CREDENTIAL"),

		MONGODB_URI:      os.Getenv("MONGODB_URI"),
		MONGODB_DATABASE: os.Getenv("MONGODB_DATABASE"),

		CLIENTS_CONFIG_PATH: os.Getenv("CLIENTS_CONFIG_PATH"),
		EVM_PRIVATE_KEY:     os.Getenv("EVM_PRIVATE_KEY"),
	}

	validate := validator.New()
	err := validate.Struct(Env)

	if err != nil {
		panic(err)
	}
}
