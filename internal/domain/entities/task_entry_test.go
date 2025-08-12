package entities

import (
	"daily-tracker/internal/domain/valueobjects"
	"testing"
	"time"
)

// В Go тесты пишутся в том же пакете с суффиксом _test.go
// Функции тестов должны начинаться с Test и принимать *testing.T

func TestNewTaskEntry_Success(t *testing.T) {
	// Arrange (подготовка данных)
	id := TaskEntryID("test-id-123")
	date := time.Now()
	dayNumber := 1
	keyTask := "Написать тесты"
	category, _ := valueobjects.NewTaskCategory("работа")
	stressBefore, _ := valueobjects.NewStressLevel(7)

	// Act (выполнение тестируемого кода)
	taskEntry, err := NewTaskEntry(id, date, dayNumber, keyTask, category, stressBefore)

	// Assert (проверка результатов)
	// В Go нет встроенных assertions, используем простые if
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if taskEntry == nil {
		t.Fatal("Expected taskEntry to be created, got nil")
		// t.Fatal останавливает выполнение теста
	}

	// Проверяем все поля
	if taskEntry.ID() != id {
		t.Errorf("Expected ID %s, got %s", id, taskEntry.ID())
	}

	if taskEntry.KeyTask() != keyTask {
		t.Errorf("Expected key task %s, got %s", keyTask, taskEntry.KeyTask())
	}

	if taskEntry.Started() {
		t.Error("Expected task to not be started initially")
	}

	if taskEntry.DayNumber() != dayNumber {
		t.Errorf("Expected day number %d, got %d", dayNumber, taskEntry.DayNumber())
	}
}

func TestNewTaskEntry_EmptyKeyTask(t *testing.T) {
	// Тестируем валидацию
	id := TaskEntryID("test-id")
	date := time.Now()
	category, _ := valueobjects.NewTaskCategory("работа")
	stressBefore, _ := valueobjects.NewStressLevel(5)

	// Пустая задача должна вызвать ошибку
	taskEntry, err := NewTaskEntry(id, date, 1, "", category, stressBefore)

	if err == nil {
		t.Error("Expected error for empty key task, got nil")
	}

	if taskEntry != nil {
		t.Error("Expected taskEntry to be nil when error occurs")
	}
}

func TestNewTaskEntry_InvalidDayNumber(t *testing.T) {
	id := TaskEntryID("test-id")
	date := time.Now()
	keyTask := "Test task"
	category, _ := valueobjects.NewTaskCategory("работа")
	stressBefore, _ := valueobjects.NewStressLevel(5)

	// Отрицательный номер дня должен вызвать ошибку
	taskEntry, err := NewTaskEntry(id, date, -1, keyTask, category, stressBefore)

	if err == nil {
		t.Error("Expected error for negative day number, got nil")
	}

	if taskEntry != nil {
		t.Error("Expected taskEntry to be nil when error occurs")
	}
}

func TestTaskEntry_StartTask(t *testing.T) {
	// Подготавливаем задачу
	taskEntry := createValidTaskEntry(t)

	// Проверяем начальное состояние
	if taskEntry.Started() {
		t.Error("Task should not be started initially")
	}

	if taskEntry.StartTime() != nil {
		t.Error("Start time should be nil initially")
	}

	// Запускаем задачу
	err := taskEntry.StartTask()
	if err != nil {
		t.Errorf("Expected no error when starting task, got: %v", err)
	}

	// Проверяем состояние после запуска
	if !taskEntry.Started() {
		t.Error("Task should be started after StartTask()")
	}

	if taskEntry.StartTime() == nil {
		t.Error("Start time should not be nil after starting task")
	}

	// Проверяем доменные события
	events := taskEntry.DomainEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 domain event, got %d", len(events))
	}

	if len(events) > 0 {
		event := events[0]
		if event.EventType() != "TaskStarted" {
			t.Errorf("Expected TaskStarted event, got %s", event.EventType())
		}
	}
}

func TestTaskEntry_StartTaskTwice(t *testing.T) {
	taskEntry := createValidTaskEntry(t)

	// Первый запуск должен быть успешным
	err := taskEntry.StartTask()
	if err != nil {
		t.Errorf("First StartTask() should succeed, got: %v", err)
	}

	// Второй запуск должен вернуть ошибку
	err = taskEntry.StartTask()
	if err == nil {
		t.Error("Second StartTask() should return error")
	}
}

func TestTaskEntry_UpdateDuration(t *testing.T) {
	taskEntry := createValidTaskEntry(t)

	// Обновление длительности без запуска должно вернуть ошибку
	err := taskEntry.UpdateDuration(30 * time.Minute)
	if err == nil {
		t.Error("UpdateDuration should fail for unstarted task")
	}

	// Запускаем задачу
	taskEntry.StartTask()

	// Теперь обновление должно работать
	duration := 25 * time.Minute
	err = taskEntry.UpdateDuration(duration)
	if err != nil {
		t.Errorf("UpdateDuration should succeed for started task, got: %v", err)
	}

	if taskEntry.ActiveDuration() != duration {
		t.Errorf("Expected duration %v, got %v", duration, taskEntry.ActiveDuration())
	}
}

func TestTaskEntry_UpdateDuration_NegativeDuration(t *testing.T) {
	taskEntry := createValidTaskEntry(t)
	taskEntry.StartTask()

	// Отрицательная длительность должна вернуть ошибку
	err := taskEntry.UpdateDuration(-10 * time.Minute)
	if err == nil {
		t.Error("UpdateDuration should fail for negative duration")
	}
}

func TestTaskEntry_SetStressAfter(t *testing.T) {
	taskEntry := createValidTaskEntry(t)
	stressAfter, _ := valueobjects.NewStressLevel(3)

	// Устанавливаем стресс после
	taskEntry.SetStressAfter(stressAfter)

	// Проверяем, что значение установилось
	if taskEntry.stressAfter != stressAfter {
		t.Errorf("Expected stress after %d, got %d", stressAfter, taskEntry.stressAfter)
	}

	// Проверяем, что создалось событие
	events := taskEntry.DomainEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 domain event, got %d", len(events))
	}

	if len(events) > 0 {
		if events[0].EventType() != "StressLevelChanged" {
			t.Errorf("Expected StressLevelChanged event, got %s", events[0].EventType())
		}
	}
}

func TestTaskEntry_CalculateStressReduction(t *testing.T) {
	taskEntry := createValidTaskEntry(t)
	stressAfter, _ := valueobjects.NewStressLevel(3)
	taskEntry.SetStressAfter(stressAfter)

	// stressBefore = 7 (из createValidTaskEntry), stressAfter = 3
	expectedReduction := 7 - 3
	actualReduction := taskEntry.CalculateStressReduction()

	if actualReduction != expectedReduction {
		t.Errorf("Expected stress reduction %d, got %d", expectedReduction, actualReduction)
	}
}

// Вспомогательная функция для создания валидной записи задачи
// В Go принято выносить общую логику в helper-функции
func createValidTaskEntry(t *testing.T) *TaskEntry {
	id := TaskEntryID("test-id-123")
	date := time.Now()
	dayNumber := 1
	keyTask := "Test task"
	category, err := valueobjects.NewTaskCategory("работа")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	stressBefore, err := valueobjects.NewStressLevel(7)
	if err != nil {
		t.Fatalf("Failed to create stress level: %v", err)
	}

	taskEntry, err := NewTaskEntry(id, date, dayNumber, keyTask, category, stressBefore)
	if err != nil {
		t.Fatalf("Failed to create task entry: %v", err)
	}

	return taskEntry
}

// Бенчмарки - измерение производительности
// Функции должны начинаться с Benchmark и принимать *testing.B
func BenchmarkNewTaskEntry(b *testing.B) {
	// Подготавливаем данные один раз
	id := TaskEntryID("test-id-123")
	date := time.Now()
	keyTask := "Benchmark task"
	category, _ := valueobjects.NewTaskCategory("работа")
	stressBefore, _ := valueobjects.NewStressLevel(7)

	// b.N - количество итераций, определяется автоматически
	for i := 0; i < b.N; i++ {
		_, err := NewTaskEntry(id, date, i+1, keyTask, category, stressBefore)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

// Table-driven tests - популярный паттерн в Go
func TestTaskEntry_StartTask_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		taskStarted bool
		expectError bool
	}{
		{
			name:        "start fresh task",
			taskStarted: false,
			expectError: false,
		},
		{
			name:        "start already started task",
			taskStarted: true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		// t.Run создает подтест с именем
		t.Run(tt.name, func(t *testing.T) {
			taskEntry := createValidTaskEntry(t)

			if tt.taskStarted {
				// Предварительно запускаем задачу
				taskEntry.StartTask()
			}

			err := taskEntry.StartTask()

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
