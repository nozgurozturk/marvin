package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	configs *configurations
)

type configurations struct {
	HTTP  *httpConfig
	SMTP  *smtpConfig
	Mongo *mongoConfig
}

type httpConfig struct {
	MainHost  string
	MainPort  string
	Host      string
	Port      string
	SubSecret string
	SubExpire int64
}

type smtpConfig struct {
	Port     string
	Host     string
	From     string
	Password string
}

type mongoConfig struct {
	Username string
	Password string
	Host     string
	DBName   string
	Query    string
}

func Set() *configurations {

	// load .env file
	env := os.Getenv("GO_ENV")
	err := godotenv.Load(".env")

	if err != nil && env == "dev" {
		log.Fatal("Error loading .env file")
	}

	subExpire, err := strconv.ParseInt(os.Getenv("SUB_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	cnf := new(configurations)
	// notifier client config
	cnf.HTTP = &httpConfig{
		MainHost:  os.Getenv("MAIN_HOST"),
		MainPort:  ":" + os.Getenv("MAIN_PORT"),
		Host:      os.Getenv("HOST"),
		SubSecret: os.Getenv("SUB_SECRET"),
		SubExpire: subExpire,
		Port:      ":" + os.Getenv("PORT"),
	}

	// email client config
	cnf.SMTP = &smtpConfig{
		Port:     os.Getenv("EMAIL_PORT"),
		Host:     os.Getenv("EMAIL_HOST"),
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

	// mongo_db config
	cnf.Mongo = &mongoConfig{
		Username: os.Getenv("MONGO_DB_USERNAME"),
		Password: os.Getenv("MONGO_DB_PASSWORD"),
		Host:     os.Getenv("MONGO_DB_HOST"),
		DBName:   os.Getenv("MONGO_DB_NAME"),
		Query:    os.Getenv("MONGO_DB_QUERY"),
	}

	configs = cnf
	return configs
}

func Get() *configurations {
	return configs
}
