package usecase

import "github.com/IgorAndrade/analytics-twitter/server/internal/model"

type Poster interface {
	Save(p model.Post)
}

func (t Twitter) Save(p model.Post) {
	t.repository.Post(p)
}
