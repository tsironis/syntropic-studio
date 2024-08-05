package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/tsironis/syntropic-studio/db"
)

func NewTodoServices(t Todo, tStore db.Store) *TodoServices {

	return &TodoServices{
		Todo:      t,
		TodoStore: tStore,
	}
}

type Todo struct {
	ID          int       `json:"id"`
	CreatedBy   int       `json:"created_by"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      bool      `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type TodoServices struct {
	Todo      Todo
	TodoStore db.Store
}

func (ts *TodoServices) CreateTodo(t Todo) (Todo, error) {
	result := ts.TodoStore.Db.Create(&t)
	if result.Error != nil {
		return Todo{}, result.Error
	}

	fmt.Printf("Todo created: %d\n", t)

	return ts.Todo, nil
}

func (ts *TodoServices) GetAllTodos(createdBy int) ([]Todo, error) {
	todos := []Todo{}

	result := ts.TodoStore.Db.Where("created_by = ?", createdBy).Order("created_at desc").Find(&todos)
	if result.Error != nil {
		return []Todo{}, result.Error
	}
	return todos, nil
}

func (ts *TodoServices) GetTodoById(t Todo) (Todo, error) {
	result := ts.TodoStore.Db.Where(&t, "created_by", "id").Find(&ts.Todo)
	if result.Error != nil {
		return Todo{}, result.Error
	}

	return ts.Todo, nil
}

func (ts *TodoServices) UpdateTodo(t Todo) (Todo, error) {
	result := ts.TodoStore.Db.Model(&ts.Todo).Where(&t, "created_by", "id").Updates(&t)
	if result.Error != nil {
		return Todo{}, result.Error
	}

	return ts.Todo, nil
}

func (ts *TodoServices) DeleteTodo(t Todo) error {
	result := ts.TodoStore.Db.Delete(&t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}

func ConvertDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}
