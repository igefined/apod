package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

type Config struct {
	DBHost         string `env:"DB_HOST"`
	DBPort         string `env:"DB_PORT"`
	DBName         string `env:"DB_NAME"`
	DBUser         string `env:"DB_USER"`
	Port           string `env:"PORT" envDefault:"80"`
	Host           string `env:"HOST" envDefault:"0.0.0.0"`
	DBPassword     string `env:"DB_PASSWORD"`
	NASAApiKey     string `env:"NASA_API_KEY"`
	AWSAccessKey   string `env:"AWS_ACCESS_KEY"`
	AWSSecretKey   string `env:"AWS_SECRET_KEY"`
	AWSRegion      string `env:"AWS_REGION"`
	AWSBucketName  string `env:"AWS_BUCKET_NAME"`
	JWTSecretKey   string `env:"JWT_SECRET_KEY"`
	TokenSecretKey string `env:"TOKEN_SECRET_KEY"`
}

var instance *Config
var once sync.Once

func NewConfig() *Config {
	once.Do(configEnv)
	return instance
}

func configEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	instance = &Config{}
	if err := env.Parse(instance); err != nil {
		log.Fatal(err)
	}
}
