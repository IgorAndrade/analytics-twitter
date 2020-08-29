package usecase

import "github.com/IgorAndrade/analytics-twitter/server/internal/model"

type Poster interface {
	SavePost(ID int64, p model.Post)
}

func (t Twitter) SavePost(ID int64, p model.Post) {
	t.r.Post(ID, p)
}
