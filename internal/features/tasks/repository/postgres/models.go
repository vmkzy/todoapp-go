package tasks_postgres_repository

import (
	"time"

	"github.com/vmkzy/todoapp-go/internal/core/domain"
	core_postgres_pool "github.com/vmkzy/todoapp-go/internal/core/repository/postgres/pool"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func (m *TaskModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Title,
		&m.Description,
		&m.Completed,
		&m.CreatedAt,
		&m.CompletedAt,
		&m.AuthorUserID,
	)
}

func modelToDomain(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}

func modelsToDomains(models []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
