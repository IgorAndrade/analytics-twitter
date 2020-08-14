package service

import "github.com/IgorAndrade/analytics-twitter/server/internal/model"

type Poster interface {
	Post(int64, model.Post) error
}

type poster struct {
}
