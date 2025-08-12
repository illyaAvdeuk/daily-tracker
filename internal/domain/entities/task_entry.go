package entities

import (
	"daily-tracker/internal/domain/valueobjects"
	"daily-tracker/pkg/errors"
	"time"
)

// TaskEntry представляет запись о выполнении задачи
// В DDD это Entity - объект с уникальной идентичностью
type TaskEntry struct {
	id              TaskEntryID               // Уникальный идентификатор
	date            time.Time                 // Дата выполнения
	dayNumber       int                       // Номер дня в периоде
	keyTask         string                    // Основная задача
	category        valueobjects.TaskCategory // Категория задачи (Value Object)
	stressBefore    valueobjects.StressLevel  // Уровень стресса до (0-10)
	started         bool                      // Начата ли задача
	startTime       *time.Time                // Время начала (может быть nil)
	activeDuration  time.Duration             // Активное время выполнения
	continuedAfter  bool                      // Продолжалась ли после 10 мин
	stressAfter     valueobjects.StressLevel  // Уровень стресса после
	distractions    time.Duration             // Время отвлечений
	blocksCompleted int                       // Количество завершенных блоков
	pomodoroCount   int                       // Количество помидорок
	lightExposure   time.Duration             // Время на свету
	energy          valueobjects.EnergyLevel  // Уровень энергии (0-10)
	mood            valueobjects.MoodLevel    // Уровень настроения (0-10)
	notes           string                    // Заметки

	// DDD: Domain Events для отслеживания изменений
	domainEvents []DomainEvent
}

// TaskEntryID - строго типизированный ID (Go идиома)
// В отличие от PHP, где ID часто int, в Go принято создавать типы
type TaskEntryID string

// DomainEvent интерфейс для доменных событий
type DomainEvent interface {
	OccurredOn() time.Time
	EventType() string
}

// Конструктор для создания новой записи задачи
// В Go нет ключевого слова constructor, используем функции-фабрики
func NewTaskEntry(
	id TaskEntryID,
	date time.Time,
	dayNumber int,
	keyTask string,
	category valueobjects.TaskCategory,
	stressBefore valueobjects.StressLevel,
) (*TaskEntry, error) {
	// Валидация входных данных на уровне домена
	if keyTask == "" {
		return nil, errors.NewDomainError("key task cannot be empty")
	}

	if dayNumber < 1 {
		return nil, errors.NewDomainError("day number must be positive")
	}

	return &TaskEntry{
		id:           id,
		date:         date,
		dayNumber:    dayNumber,
		keyTask:      keyTask,
		category:     category,
		stressBefore: stressBefore,
		started:      false,
		domainEvents: make([]DomainEvent, 0),
	}, nil
}

// Геттеры (в Go принято не использовать префикс Get)
func (te *TaskEntry) ID() TaskEntryID {
	return te.id
}

func (te *TaskEntry) Date() time.Time {
	return te.date
}

func (te *TaskEntry) DayNumber() int {
	return te.dayNumber
}

func (te *TaskEntry) KeyTask() string {
	return te.keyTask
}

func (te *TaskEntry) Category() valueobjects.TaskCategory {
	return te.category
}

func (te *TaskEntry) StressBefore() valueobjects.StressLevel {
	return te.stressBefore
}

func (te *TaskEntry) Started() bool {
	return te.started
}

func (te *TaskEntry) StartTime() *time.Time {
	return te.startTime
}

func (te *TaskEntry) ActiveDuration() time.Duration {
	return te.activeDuration
}

func (te *TaskEntry) ContinuedAfter() bool {
	return te.continuedAfter
}

func (te *TaskEntry) StressAfter() valueobjects.StressLevel {
	return te.stressAfter
}

func (te *TaskEntry) Distractions() time.Duration {
	return te.distractions
}

func (te *TaskEntry) BlocksCompleted() int {
	return te.blocksCompleted
}
func (te *TaskEntry) PomodoroCount() int {
	return te.pomodoroCount
}

func (te *TaskEntry) LightExposure() time.Duration {
	return te.lightExposure
}

func (te *TaskEntry) Energy() valueobjects.EnergyLevel {
	return te.energy
}

func (te *TaskEntry) Mood() valueobjects.MoodLevel {
	return te.mood
}

// Доменные методы - бизнес-логика инкапсулирована в Entity

// StartTask начинает выполнение задачи
func (te *TaskEntry) StartTask() error {
	if te.started {
		return errors.NewDomainError("task already started")
	}

	now := time.Now()
	te.started = true
	te.startTime = &now

	// Генерируем доменное событие
	te.addDomainEvent(&TaskStartedEvent{
		taskEntryID: te.id,
		occurredOn:  now,
	})

	return nil
}

// UpdateDuration обновляет продолжительность активной работы
func (te *TaskEntry) UpdateDuration(duration time.Duration) error {
	if !te.started {
		return errors.NewDomainError("cannot update duration: task not started")
	}

	if duration < 0 {
		return errors.NewDomainError("duration cannot be negative")
	}

	te.activeDuration = duration
	return nil
}

// SetStressAfter устанавливает уровень стресса после выполнения
func (te *TaskEntry) SetStressAfter(stressLevel valueobjects.StressLevel) {
	te.stressAfter = stressLevel

	// Генерируем событие об изменении стресса
	te.addDomainEvent(&StressLevelChangedEvent{
		taskEntryID:  te.id,
		stressBefore: te.stressBefore,
		stressAfter:  stressLevel,
		occurredOn:   time.Now(),
	})
}

// CalculateStressReduction вычисляет снижение стресса
func (te *TaskEntry) CalculateStressReduction() int {
	return int(te.stressBefore) - int(te.stressAfter)
}

// AddNotes добавляет заметки к записи
func (te *TaskEntry) AddNotes(notes string) {
	te.notes = notes
}

// DomainEvents возвращает список доменных событий
func (te *TaskEntry) DomainEvents() []DomainEvent {
	return te.domainEvents
}

// ClearDomainEvents очищает список событий (обычно после публикации)
func (te *TaskEntry) ClearDomainEvents() {
	te.domainEvents = make([]DomainEvent, 0)
}

// Приватный метод для добавления доменных событий
func (te *TaskEntry) addDomainEvent(event DomainEvent) {
	te.domainEvents = append(te.domainEvents, event)
}

// Доменные события

// TaskStartedEvent событие начала задачи
type TaskStartedEvent struct {
	taskEntryID TaskEntryID
	occurredOn  time.Time
}

func (e *TaskStartedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *TaskStartedEvent) EventType() string {
	return "TaskStarted"
}

func (e *TaskStartedEvent) TaskEntryID() TaskEntryID {
	return e.taskEntryID
}

// StressLevelChangedEvent событие изменения уровня стресса
type StressLevelChangedEvent struct {
	taskEntryID  TaskEntryID
	stressBefore valueobjects.StressLevel
	stressAfter  valueobjects.StressLevel
	occurredOn   time.Time
}

func (e *StressLevelChangedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *StressLevelChangedEvent) EventType() string {
	return "StressLevelChanged"
}

func (e *StressLevelChangedEvent) TaskEntryID() TaskEntryID {
	return e.taskEntryID
}

func (e *StressLevelChangedEvent) StressBefore() valueobjects.StressLevel {
	return e.stressBefore
}

func (e *StressLevelChangedEvent) StressAfter() valueobjects.StressLevel {
	return e.stressAfter
}
