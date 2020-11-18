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
	Redis *redisConfig
}

type httpConfig struct {
	Host          string
	Port          string
	AccessSecret  string
	AccessExpire  int64
	RefreshExpire int64
	RefreshSecret string
	SubSecret     string
	SubExpire     int64
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

type redisConfig struct {
	Address  string
	UserName string
	Password string
	DB       int
}

func Set() *configurations {

	// load .env file
	env := os.Getenv("GO_ENV")
	err := godotenv.Load(".env")

	if err != nil && env == "dev" {
		log.Fatal("Error loading .env file")
	}

	// marvin-server config
	atExpire, err := strconv.ParseInt(os.Getenv("ACCESS_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	rtExpire, err := strconv.ParseInt(os.Getenv("REFRESH_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	subExpire, err := strconv.ParseInt(os.Getenv("SUB_EXPIRE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	cnf := new(configurations)
	// api client config
	cnf.HTTP = &httpConfig{
		Host:          os.Getenv("HOST"),
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		AccessExpire:  atExpire,
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		RefreshExpire: rtExpire,
		SubSecret:     os.Getenv("SUB_SECRET"),
		SubExpire:     subExpire,
		Port:          ":" + os.Getenv("PORT"),
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

	// redis config
	cnf.Redis = &redisConfig{
		Address:  os.Getenv("REDIS_DB_ADDRESS"),
		UserName: os.Getenv("REDIS_DB_USERNAME"),
		Password: os.Getenv("REDIS_DB_PASSWORD"),
		DB:       0,
	}
	configs = cnf
	return configs
}

func Get() *configurations {
	return configs
}
