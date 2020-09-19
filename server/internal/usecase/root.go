package usecase

import (
	"github.com/IgorAndrade/analytics-twitter/server/internal/repository"
	"github.com/sarulabs/di"
)

const TWITTER = "UC-Twitter"

type Usecase interface {
	Poster
}

type Twitter struct {
	repository repository.Elasticsearch
}

func Define(b *di.Builder) {
	b.Add(di.Def{
		Name:  TWITTER,
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			r := ctn.Get(repository.ELASTICSEARCH).(repository.Elasticsearch)
			return new(r), nil
		},
	})
}

func new(r repository.Elasticsearch) Usecase {
	return &Twitter{repository: r}
}
