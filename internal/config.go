package internal

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token           string `required:"true"`
	InitialUserID   int64
	InitialUserName string
	Database        string `default:"data.db" required:"true"`
}

func ReadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("linkbot", &cfg)
	log.Printf("userid: %d", cfg.InitialUserID)
	return &cfg, err
}
