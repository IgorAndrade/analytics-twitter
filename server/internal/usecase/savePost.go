package usecase

import "github.com/IgorAndrade/analytics-twitter/server/internal/model"

type Poster interface {
	Save(ID int64, p model.Post)
}

func (t Twitter) Save(ID int64, p model.Post) {
	t.repository.Post(ID, p)
}
