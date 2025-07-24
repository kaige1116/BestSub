package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/task"
)

type TaskRepository interface {
	Create(ctx context.Context, task *task.Data) error
	Update(ctx context.Context, task *task.Data) error
	Delete(ctx context.Context, id uint16) error
	GetByID(ctx context.Context, id uint16) (*task.Data, error)
	List(ctx context.Context) (*[]task.Data, error)
}
