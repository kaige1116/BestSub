package op

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/cache"
)

var taskRepo interfaces.TaskRepository
var taskCache = cache.New[uint16, task.Data](16)

func TaskRepo() interfaces.TaskRepository {
	if taskRepo == nil {
		taskRepo = repo.Task()
	}
	return taskRepo
}
func GetTaskByID(id uint16) (task.Data, error) {
	if taskCache.Len() == 0 {
		if err := refreshTaskCache(context.Background()); err != nil {
			return task.Data{}, err
		}
	}
	if t, ok := taskCache.Get(id); ok {
		return t, nil
	}
	return task.Data{}, fmt.Errorf("task not found")
}
func CreateTask(ctx context.Context, t *task.Data) error {
	if taskCache.Len() == 0 {
		if err := refreshTaskCache(context.Background()); err != nil {
			return err
		}
	}
	if err := TaskRepo().Create(ctx, t); err != nil {
		return err
	}
	taskCache.Set(t.ID, *t)
	return nil
}
func UpdateTask(ctx context.Context, t *task.Data) error {
	if taskCache.Len() == 0 {
		if err := refreshTaskCache(context.Background()); err != nil {
			return err
		}
	}
	oldTask, ok := taskCache.Get(t.ID)
	if !ok {
		return fmt.Errorf("task not found")
	}
	t.Result = oldTask.Result
	if err := TaskRepo().Update(ctx, t); err != nil {
		return err
	}
	taskCache.Set(t.ID, *t)
	return nil
}
func UpdateTaskResult(ctx context.Context, id uint16, result string) error {
	if taskCache.Len() == 0 {
		if err := refreshTaskCache(ctx); err != nil {
			return err
		}
	}
	oldTask, ok := taskCache.Get(id)
	if !ok {
		return fmt.Errorf("task not found")
	}
	oldTask.Result = result
	if err := TaskRepo().Update(ctx, &oldTask); err != nil {
		return err
	}
	taskCache.Set(id, oldTask)
	return nil
}
func DeleteTask(ctx context.Context, id uint16) error {
	if taskCache.Len() == 0 {
		if err := refreshTaskCache(context.Background()); err != nil {
			return err
		}
	}
	if err := TaskRepo().Delete(ctx, id); err != nil {
		return err
	}
	taskCache.Del(id)
	return nil
}
func GetTaskList(ctx context.Context) ([]task.Data, error) {
	taskList := taskCache.GetAll()
	if len(taskList) == 0 {
		err := refreshTaskCache(context.Background())
		if err != nil {
			return nil, err
		}
		taskList = taskCache.GetAll()
	}
	var result = make([]task.Data, 0, len(taskList))
	for _, v := range taskList {
		result = append(result, v)
	}
	return result, nil
}
func refreshTaskCache(ctx context.Context) error {
	taskCache.Clear()
	tasks, err := TaskRepo().List(ctx)
	if err != nil {
		return err
	}
	for _, task := range *tasks {
		taskCache.Set(task.ID, task)
	}
	return nil
}
