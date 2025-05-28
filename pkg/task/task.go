package task

import "context"

// Task 定时任务接口
type Task interface {
	// GetName 获取任务名称
	GetName() string
	// GetCronExpr 获取cron表达式
	GetCronExpr() string
	// IsEnabled 是否启用
	IsEnabled() bool
	// Execute 执行任务
	Execute(ctx context.Context) error
}

// TaskInfo 任务信息
type TaskInfo struct {
	Name     string `json:"name"`
	CronExpr string `json:"cron_expr"`
	Enabled  bool   `json:"enabled"`
}
