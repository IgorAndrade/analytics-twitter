package service

import "github.com/IgorAndrade/analytics-twitter/server/internal/model"

type Poster interface {
	Post(model.Post) error
}

type poster struct {
}
