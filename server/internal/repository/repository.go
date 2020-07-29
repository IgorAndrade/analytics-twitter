package repository

import (
	"context"

	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
)

const TODO_LIST = "tudoListRepo"

type TodoList interface {
	Create(context.Context, *model.TodoList) error
	GetAll(context.Context) ([]model.TodoList, error)
}
