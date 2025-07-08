package task

import (
	"context"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/core/task/register"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// executeTask 执行任务
func (s *Scheduler) executeTask(taskData *task.Data) {
	ctx := context.Background()
	taskID := taskData.ID

	// 检查任务是否已在运行
	if _, running := s.runningTasks.Load(taskID); running {
		log.Warnf("任务 %d 已在运行中，跳过本次执行", taskID)
		return
	}

	// 标记任务开始运行
	s.runningTasks.Store(taskID, true)
	defer s.runningTasks.Delete(taskID)

	// 更新任务状态为运行中
	taskData.Status = task.StatusRunning
	if err := s.updateTaskStatus(ctx, taskData); err != nil {
		log.Errorf("Failed to update task %d status to running: %v", taskID, err)
	}

	startTime := time.Now()
	var success bool
	var resultMsg string

	// 执行任务逻辑
	success, resultMsg = s.executeTaskLogic(ctx, taskData)

	// 根据执行结果更新任务状态
	if success {
		taskData.Status = task.StatusCompleted
		taskData.SuccessCount++
		log.Infof("任务 %d (%s) 执行成功", taskID, taskData.Name)
	} else {
		taskData.Status = task.StatusFailed
		taskData.FailedCount++
		log.Errorf("任务 %d (%s) 执行失败: %s", taskID, taskData.Name, resultMsg)
	}

	// 更新任务执行结果
	duration := int(time.Since(startTime).Milliseconds())
	taskData.LastRunTime = &startTime
	taskData.LastRunDuration = &duration
	taskData.LastRunResult = resultMsg

	if err := s.updateTaskStatus(ctx, taskData); err != nil {
		log.Errorf("Failed to update task %d execution result: %v", taskID, err)
	}

	// 写入日志文件
	s.writeTaskLog(taskData, startTime, success, resultMsg)
}

// executeTaskLogic 执行任务的核心逻辑
func (s *Scheduler) executeTaskLogic(ctx context.Context, taskData *task.Data) (bool, string) {
	// 从注册表获取处理器
	handler, exists := register.GetHandler(taskData.Type)
	if !exists {
		return false, fmt.Sprintf("未找到任务类型处理器: %s", taskData.Type)
	}

	// 验证配置
	if err := handler.Validate(taskData.Config); err != nil {
		return false, fmt.Sprintf("任务配置无效: %v", err)
	}

	// 执行任务
	if err := handler.Execute(ctx, taskData.Config); err != nil {
		return false, fmt.Sprintf("执行失败: %v", err)
	}

	return true, "执行成功"
}

// updateTaskStatus 更新任务状态到数据库
func (s *Scheduler) updateTaskStatus(ctx context.Context, taskData *task.Data) error {
	return s.repo.Update(ctx, taskData)
}

// writeTaskLog 写入任务日志到文件
func (s *Scheduler) writeTaskLog(taskData *task.Data, startTime time.Time, success bool, message string) {
	// 生成执行ID
	executionID := fmt.Sprintf("exec_%d_%d", taskData.ID, startTime.Unix())

	// 创建日志记录
	taskLog := TaskLog{
		TaskID:      taskData.ID,
		ExecutionID: executionID,
		Timestamp:   startTime,
		Level:       "STATUS",
		Message:     message,
		Status:      taskData.Status,
		Progress:    100,
	}

	// 写入日志文件
	if err := WriteTaskLog(taskLog, success); err != nil {
		log.Errorf("Failed to write task log for task %d: %v", taskData.ID, err)
	}
}
