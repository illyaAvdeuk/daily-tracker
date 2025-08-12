package repositories

import (
	"context"
	"daily-tracker/internal/domain/entities"
	"time"
)

// TaskRepository определяет контракт для работы с задачами
// В Go интерфейсы маленькие и сфокусированные (Interface Segregation)
type TaskRepository interface {
	// Save сохраняет или обновляет запись задачи
	// context.Context - стандартный способ передачи метаданных в Go
	Save(ctx context.Context, task *entities.TaskEntry) error

	// FindByID находит задачу по ID
	// Возвращает указатель и ошибку (идиома Go)
	FindByID(ctx context.Context, id entities.TaskEntryID) (*entities.TaskEntry, error)

	// FindByDate находит все задачи за определенную дату
	// []* - слайс указателей на TaskEntry
	FindByDate(ctx context.Context, date time.Time) ([]*entities.TaskEntry, error)

	// FindByDateRange находит задачи в диапазоне дат
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.TaskEntry, error)

	// Delete удаляет задачу
	Delete(ctx context.Context, id entities.TaskEntryID) error

	// Exists проверяет существование записи
	Exists(ctx context.Context, id entities.TaskEntryID) (bool, error)
}

// Дополнительный интерфейс для расширенных операций
// Принцип: много маленьких интерфейсов лучше одного большого
type TaskStatisticsRepository interface {
	// GetTaskCountByCategory возвращает количество задач по категориям
	GetTaskCountByCategory(ctx context.Context, startDate, endDate time.Time) (map[string]int, error)

	// GetAverageStressReduction вычисляет среднее снижение стресса
	GetAverageStressReduction(ctx context.Context, startDate, endDate time.Time) (float64, error)
}

// Композиция интерфейсов - уникальная особенность Go
// Можно "встраивать" интерфейсы друг в друга
type FullTaskRepository interface {
	TaskRepository           // Встроенный интерфейс
	TaskStatisticsRepository // Еще один встроенный интерфейс

	// Дополнительные методы только для полной реализации
	Backup(ctx context.Context, filePath string) error
	Restore(ctx context.Context, filePath string) error
}

// Пример интерфейса для кеширования
type TaskCache interface {
	Get(key string) (*entities.TaskEntry, bool)
	Set(key string, task *entities.TaskEntry, ttl time.Duration)
	Delete(key string)
	Clear()
}

// Reader и Writer интерфейсы (паттерн разделения чтения/записи)
type TaskReader interface {
	FindByID(ctx context.Context, id entities.TaskEntryID) (*entities.TaskEntry, error)
	FindByDate(ctx context.Context, date time.Time) ([]*entities.TaskEntry, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.TaskEntry, error)
	Exists(ctx context.Context, id entities.TaskEntryID) (bool, error)
}

type TaskWriter interface {
	Save(ctx context.Context, task *entities.TaskEntry) error
	Delete(ctx context.Context, id entities.TaskEntryID) error
}

// Полный репозиторий через композицию интерфейсов
type TaskReadWriter interface {
	TaskReader
	TaskWriter
}
