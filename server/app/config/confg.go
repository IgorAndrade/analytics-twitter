package config

import (
	"fmt"
	"os"

	"github.com/sarulabs/di"
)

type Config struct {
	Rest          Rest
	Mongo         Mongo
	Twitter       Twitter
	Elasticsearch Elasticsearch
}

type Rest struct {
	Port string
}

type Mongo struct {
	Address  string
	User     string
	Password string
}

type Elasticsearch struct {
	Address  string
	Username string
	Password string
}

type Twitter struct {
	Token          string
	TokenSecret    string
	ConsumerKey    string
	ConsumerSecret string
}

var NAME = "config"
var c *Config

func Define(b *di.Builder) {
	c = &Config{
		Rest: Rest{
			Port: fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		},
		Mongo: Mongo{
			Address:  os.Getenv("MONGO_URL"),
			User:     os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASSWORD"),
		},
		Twitter: Twitter{
			Token:          os.Getenv("TWITTER_TOKEN"),
			TokenSecret:    os.Getenv("TWITTER_TOKEN_SECRET"),
			ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
			ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		},
		Elasticsearch: Elasticsearch{
			Address:  os.Getenv("ELASTICSEARCH_ADDRESS"),
			Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
			Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		},
	}
	b.Add(di.Def{
		Name:  NAME,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return c, nil
		},
	})
}
