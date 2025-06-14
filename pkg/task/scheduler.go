package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"xhyovo.cn/community/pkg/log"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	cron  *cron.Cron
	tasks map[string]Task
	mu    sync.RWMutex
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		cron:  cron.New(cron.WithSeconds()),
		tasks: make(map[string]Task),
	}
}

// RegisterTask 注册任务
func (s *TaskScheduler) RegisterTask(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !task.IsEnabled() {
		log.Infof("任务 %s 已禁用，跳过注册", task.GetName())
		return nil
	}

	// 保存任务
	s.tasks[task.GetName()] = task

	// 添加到cron调度器
	_, err := s.cron.AddFunc(task.GetCronExpr(), func() {
		s.executeTaskWrapper(task)
	})

	if err != nil {
		return fmt.Errorf("注册任务失败: %v", err)
	}

	log.Infof("任务 %s 注册成功，cron表达式: %s", task.GetName(), task.GetCronExpr())
	return nil
}

// executeTaskWrapper 执行任务的包装器
func (s *TaskScheduler) executeTaskWrapper(task Task) {
	taskName := task.GetName()

	startTime := time.Now()
	ctx := context.Background()

	err := task.Execute(ctx)
	duration := time.Since(startTime)

	if err != nil {
		log.Errorf("任务 %s 执行失败，耗时: %v，错误: %v", taskName, duration, err)
	}
}

// Start 启动调度器
func (s *TaskScheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cron.Start()
	log.Info("任务调度器已启动")
	return nil
}

// Stop 停止调度器
func (s *TaskScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Info("任务调度器已停止")
}

// GetTaskInfo 获取任务信息
func (s *TaskScheduler) GetTaskInfo(taskName string) (*TaskInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskName]
	if !exists {
		return nil, fmt.Errorf("任务 %s 不存在", taskName)
	}

	return &TaskInfo{
		Name:     task.GetName(),
		CronExpr: task.GetCronExpr(),
		Enabled:  task.IsEnabled(),
	}, nil
}

// GetAllTasks 获取所有任务信息
func (s *TaskScheduler) GetAllTasks() ([]*TaskInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*TaskInfo
	for _, task := range s.tasks {
		tasks = append(tasks, &TaskInfo{
			Name:     task.GetName(),
			CronExpr: task.GetCronExpr(),
			Enabled:  task.IsEnabled(),
		})
	}

	return tasks, nil
}

// ExecuteTaskManually 手动执行任务
func (s *TaskScheduler) ExecuteTaskManually(taskName string) error {
	s.mu.RLock()
	task, exists := s.tasks[taskName]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("任务 %s 不存在", taskName)
	}

	// 在goroutine中异步执行
	go s.executeTaskWrapper(task)
	return nil
}
