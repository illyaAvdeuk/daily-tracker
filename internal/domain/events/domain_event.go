package events

import (
	"encoding/json"
	"strconv"
	"time"
)

// DomainEvent базовый интерфейс для всех доменных событий
// Размещен в отдельном пакете для переиспользования
type DomainEvent interface {
	// OccurredOn возвращает время возникновения события
	OccurredOn() time.Time

	// EventType возвращает тип события (для сериализации и маршрутизации)
	EventType() string

	// EventID возвращает уникальный идентификатор события
	EventID() string

	// AggregateID возвращает ID агрегата, который сгенерировал событие
	AggregateID() string

	// EventVersion возвращает версию схемы события (для совместимости)
	EventVersion() int
}

// Serializable интерфейс для событий, которые нужно сериализовать
// Отдельный интерфейс следует принципу Interface Segregation
type Serializable interface {
	// ToJSON сериализует событие в JSON
	ToJSON() ([]byte, error)

	// FromJSON десериализует событие из JSON
	FromJSON(data []byte) error
}

// Publishable интерфейс для событий, которые нужно публиковать
type Publishable interface {
	// RoutingKey возвращает ключ маршрутизации для брокера сообщений
	RoutingKey() string

	// Priority возвращает приоритет события
	Priority() EventPriority
}

// EventPriority приоритет события
type EventPriority int

const (
	PriorityLow EventPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// BaseEvent базовая реализация DomainEvent
// Используется как embedded struct в конкретных событиях
type BaseEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	AggregateId string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	Version     int       `json:"version"`
}

// NewBaseEvent создает новое базовое событие
func NewBaseEvent(eventType, aggregateID string) BaseEvent {
	return BaseEvent{
		ID:          generateEventID(), // Функцию создадим позже
		Type:        eventType,
		AggregateId: aggregateID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

// Реализация интерфейса DomainEvent
func (be BaseEvent) EventID() string {
	return be.ID
}

func (be BaseEvent) EventType() string {
	return be.Type
}

func (be BaseEvent) AggregateID() string {
	return be.AggregateId
}

func (be BaseEvent) OccurredOn() time.Time {
	return be.OccurredAt
}

func (be BaseEvent) EventVersion() int {
	return be.Version
}

// ToJSON базовая сериализация
func (be BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(be)
}

// EventStore интерфейс для хранения событий (Event Sourcing)
type EventStore interface {
	// SaveEvent сохраняет событие
	SaveEvent(event DomainEvent) error

	// GetEvents получает события для агрегата
	GetEvents(aggregateID string) ([]DomainEvent, error)

	// GetEventsByType получает события определенного типа
	GetEventsByType(eventType string, limit int) ([]DomainEvent, error)
}

// EventPublisher интерфейс для публикации событий
type EventPublisher interface {
	// Publish публикует событие
	Publish(event DomainEvent) error

	// PublishBatch публикует несколько событий
	PublishBatch(events []DomainEvent) error
}

// EventHandler интерфейс для обработчиков событий
type EventHandler interface {
	// Handle обрабатывает событие
	Handle(event DomainEvent) error

	// CanHandle проверяет, может ли обработчик обработать событие
	CanHandle(eventType string) bool
}

// EventBus интерфейс для шины событий
type EventBus interface {
	// Subscribe подписывает обработчик на тип события
	Subscribe(eventType string, handler EventHandler) error

	// Unsubscribe отписывает обработчик
	Unsubscribe(eventType string, handler EventHandler) error

	// Publish публикует событие через шину
	Publish(event DomainEvent) error
}

// Временная функция генерации ID (позже заменим на UUID)
func generateEventID() string {
	return "event-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
