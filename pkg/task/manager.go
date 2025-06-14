package task

import (
	"sync"

	"xhyovo.cn/community/pkg/log"
)

var (
	once          sync.Once
	globalManager *Manager
)

// Manager 全局任务管理器
type Manager struct {
	scheduler *TaskScheduler
	mu        sync.RWMutex
}

// GetInstance 获取全局任务管理器实例（单例）
func GetInstance() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			scheduler: NewTaskScheduler(),
		}
	})
	return globalManager
}

// Initialize 初始化任务管理器
func (m *Manager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Info("开始初始化定时任务管理器")

	// 注册默认任务
	if err := m.registerDefaultTasks(); err != nil {
		return err
	}

	log.Info("定时任务管理器初始化完成")
	return nil
}

// registerDefaultTasks 注册默认任务
func (m *Manager) registerDefaultTasks() error {
	// 注册AIBase爬虫任务 - 每1小时执行一次
	aibaseTask := NewAIBaseCrawlerTask("0 0 */1 * * ?")
	if err := m.scheduler.RegisterTask(aibaseTask); err != nil {
		return err
	}

	log.Info("默认任务注册完成")
	return nil
}

// Start 启动任务管理器
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.scheduler.Start(); err != nil {
		return err
	}

	log.Info("任务管理器已启动")
	return nil
}

// Stop 停止任务管理器
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.scheduler.Stop()
	log.Info("任务管理器已停止")
}

// RegisterTask 注册新任务
func (m *Manager) RegisterTask(task Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.scheduler.RegisterTask(task)
}
