package usecase

import (
	"context"

	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
)

type Search interface {
	Find(context.Context, map[string]string) ([]model.Post, error)
}

func (t Twitter) Find(ctx context.Context, query map[string]string) ([]model.Post, error) {
	return t.repository.Find(ctx, query)
}
