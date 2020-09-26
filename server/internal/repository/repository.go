package repository

import (
	"context"

	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
)

const TODO_LIST = "tudoListRepo"
const ELASTICSEARCH = "Elasticsearch"

type TodoList interface {
	Create(context.Context, *model.TodoList) error
	GetAll(context.Context) ([]model.TodoList, error)
}

type Elasticsearch interface {
	Post(model.Post) error
	Find(ctx context.Context, query map[string]string) ([]model.Post, error)
}
