package entities

import (
	"daily-tracker/internal/domain/valueobjects"
	"daily-tracker/pkg/errors"
	"time"
)

// SleepEntry представляет запись о сне
type SleepEntry struct {
	id                 SleepEntryID                   // Уникальный идентификатор
	date               time.Time                      // Дата записи
	bedtime            time.Time                      // Время отхода ко сну
	wakeTime           time.Time                      // Время пробуждения
	sleepLatency       time.Duration                  // Время засыпания в минутах
	nightAwakenings    int                            // Количество пробуждений за ночь
	totalSleepHours    float64                        // Общее время сна в часах
	sleepQuality       valueobjects.SleepQuality      // Качество сна (0-10)
	daytimeSleepiness  valueobjects.DaytimeSleepiness // Дневная сонливость (0-10)
	caffeineAfterNoon  bool                           // Употребление кофеина после полудня
	screenUseBeforeBed time.Duration                  // Время использования экранов перед сном
	eveningFreeTime    time.Duration                  // Время отдыха вечером
	notes              string                         // Заметки

	// DDD: Domain Events
	domainEvents []DomainEvent
}

// SleepEntryID - строго типизированный ID
type SleepEntryID string

// Конструктор для создания новой записи сна
func NewSleepEntry(
	id SleepEntryID,
	date time.Time,
	bedtime, wakeTime time.Time,
	sleepQuality valueobjects.SleepQuality,
) (*SleepEntry, error) {
	// Валидация на уровне домена
	if wakeTime.Before(bedtime) {
		// Учитываем случай, когда просыпаемся на следующий день
		nextDay := bedtime.AddDate(0, 0, 1)
		if wakeTime.Before(time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, wakeTime.Location())) {
			return nil, errors.NewDomainError("wake time cannot be before bedtime on the same day")
		}
	}

	sleepEntry := &SleepEntry{
		id:           id,
		date:         date,
		bedtime:      bedtime,
		wakeTime:     wakeTime,
		sleepQuality: sleepQuality,
		domainEvents: make([]DomainEvent, 0),
	}

	// Автоматически вычисляем общее время сна
	sleepEntry.calculateTotalSleepHours()

	// Генерируем событие создания записи сна
	sleepEntry.addDomainEvent(&SleepEntryCreatedEvent{
		sleepEntryID: id,
		date:         date,
		totalHours:   sleepEntry.totalSleepHours,
		quality:      sleepQuality,
		occurredOn:   time.Now(),
	})

	return sleepEntry, nil
}

// Геттеры
func (se *SleepEntry) ID() SleepEntryID {
	return se.id
}

func (se *SleepEntry) Date() time.Time {
	return se.date
}

func (se *SleepEntry) Bedtime() time.Time {
	return se.bedtime
}

func (se *SleepEntry) WakeTime() time.Time {
	return se.wakeTime
}

func (se *SleepEntry) TotalSleepHours() float64 {
	return se.totalSleepHours
}

func (se *SleepEntry) SleepQuality() valueobjects.SleepQuality {
	return se.sleepQuality
}

func (se *SleepEntry) DaytimeSleepiness() valueobjects.DaytimeSleepiness {
	return se.daytimeSleepiness
}

func (se *SleepEntry) CaffeineAfterNoon() bool {
	return se.caffeineAfterNoon
}

func (se *SleepEntry) ScreenUseBeforeBed() time.Duration {
	return se.screenUseBeforeBed
}

func (se *SleepEntry) EveningFreeTime() time.Duration {
	return se.eveningFreeTime
}

func (se *SleepEntry) Notes() string {
	return se.notes
}

// Доменные методы с бизнес-логикой

// SetSleepLatency устанавливает время засыпания
func (se *SleepEntry) SetSleepLatency(latency time.Duration) error {
	if latency < 0 {
		return errors.NewDomainError("sleep latency cannot be negative")
	}

	if latency > 2*time.Hour {
		return errors.NewDomainError("sleep latency seems too long (over 2 hours)")
	}

	oldLatency := se.sleepLatency
	se.sleepLatency = latency

	// Генерируем событие об изменении времени засыпания
	if oldLatency != latency {
		se.addDomainEvent(&SleepLatencyChangedEvent{
			sleepEntryID: se.id,
			oldLatency:   oldLatency,
			newLatency:   latency,
			occurredOn:   time.Now(),
		})
	}

	return nil
}

// RecordNightAwakening записывает пробуждение ночью
func (se *SleepEntry) RecordNightAwakening() {
	se.nightAwakenings++

	// Генерируем событие о пробуждении
	se.addDomainEvent(&NightAwakeningRecordedEvent{
		sleepEntryID:    se.id,
		awakeningNumber: se.nightAwakenings,
		occurredOn:      time.Now(),
	})

	// Если пробуждений стало много, генерируем событие плохого сна
	if se.nightAwakenings >= 3 {
		se.addDomainEvent(&PoorSleepQualityDetectedEvent{
			sleepEntryID: se.id,
			reason:       "multiple night awakenings",
			awakenings:   se.nightAwakenings,
			occurredOn:   time.Now(),
		})
	}
}

// SetDaytimeSleepiness устанавливает дневную сонливость
func (se *SleepEntry) SetDaytimeSleepiness(sleepiness valueobjects.DaytimeSleepiness) {
	oldSleepiness := se.daytimeSleepiness
	se.daytimeSleepiness = sleepiness

	// Если сонливость изменилась значительно, генерируем событие
	if abs(int(sleepiness)-int(oldSleepiness)) >= 3 {
		se.addDomainEvent(&DaytimeSleepinessChangedEvent{
			sleepEntryID:  se.id,
			oldSleepiness: oldSleepiness,
			newSleepiness: sleepiness,
			occurredOn:    time.Now(),
		})
	}
}

// UpdateSleepQuality обновляет качество сна
func (se *SleepEntry) UpdateSleepQuality(quality valueobjects.SleepQuality) {
	oldQuality := se.sleepQuality
	se.sleepQuality = quality

	// Генерируем событие об изменении качества сна
	se.addDomainEvent(&SleepQualityUpdatedEvent{
		sleepEntryID: se.id,
		oldQuality:   oldQuality,
		newQuality:   quality,
		occurredOn:   time.Now(),
	})

	// Если качество сна стало очень плохим, генерируем специальное событие
	if quality.Int() <= 3 {
		se.addDomainEvent(&PoorSleepQualityDetectedEvent{
			sleepEntryID: se.id,
			reason:       "low quality rating",
			quality:      &quality,
			occurredOn:   time.Now(),
		})
	}
}

// IsSleepHealthy проверяет, является ли сон здоровым
func (se *SleepEntry) IsSleepHealthy() bool {
	// Бизнес-правила для здорового сна
	return se.totalSleepHours >= 7.0 &&
		se.totalSleepHours <= 9.0 &&
		se.sleepQuality.Int() >= 6 &&
		se.nightAwakenings <= 1
}

// calculateTotalSleepHours вычисляет общее время сна
func (se *SleepEntry) calculateTotalSleepHours() {
	duration := se.wakeTime.Sub(se.bedtime)

	// Если отрицательное время, значит проснулись на следующий день
	if duration < 0 {
		duration = duration + 24*time.Hour
	}

	// Вычитаем время засыпания из общего времени
	actualSleepDuration := duration - se.sleepLatency
	se.totalSleepHours = actualSleepDuration.Hours()
}

// DomainEvents возвращает список доменных событий
func (se *SleepEntry) DomainEvents() []DomainEvent {
	return se.domainEvents
}

// ClearDomainEvents очищает список событий
func (se *SleepEntry) ClearDomainEvents() {
	se.domainEvents = make([]DomainEvent, 0)
}

// Приватный метод для добавления доменных событий
func (se *SleepEntry) addDomainEvent(event DomainEvent) {
	se.domainEvents = append(se.domainEvents, event)
}

// Вспомогательная функция для вычисления модуля числа
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// === ДОМЕННЫЕ СОБЫТИЯ ДЛЯ SleepEntry ===

// SleepEntryCreatedEvent - событие создания записи сна
type SleepEntryCreatedEvent struct {
	sleepEntryID SleepEntryID
	date         time.Time
	totalHours   float64
	quality      valueobjects.SleepQuality
	occurredOn   time.Time
}

func (e *SleepEntryCreatedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *SleepEntryCreatedEvent) EventType() string {
	return "SleepEntryCreated"
}

func (e *SleepEntryCreatedEvent) SleepEntryID() SleepEntryID {
	return e.sleepEntryID
}

func (e *SleepEntryCreatedEvent) TotalHours() float64 {
	return e.totalHours
}

// SleepLatencyChangedEvent - событие изменения времени засыпания
type SleepLatencyChangedEvent struct {
	sleepEntryID SleepEntryID
	oldLatency   time.Duration
	newLatency   time.Duration
	occurredOn   time.Time
}

func (e *SleepLatencyChangedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *SleepLatencyChangedEvent) EventType() string {
	return "SleepLatencyChanged"
}

// NightAwakeningRecordedEvent - событие записи ночного пробуждения
type NightAwakeningRecordedEvent struct {
	sleepEntryID    SleepEntryID
	awakeningNumber int
	occurredOn      time.Time
}

func (e *NightAwakeningRecordedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *NightAwakeningRecordedEvent) EventType() string {
	return "NightAwakeningRecorded"
}

// PoorSleepQualityDetectedEvent - событие обнаружения плохого качества сна
type PoorSleepQualityDetectedEvent struct {
	sleepEntryID SleepEntryID
	reason       string
	awakenings   int                        // опционально
	quality      *valueobjects.SleepQuality // опционально
	occurredOn   time.Time
}

func (e *PoorSleepQualityDetectedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *PoorSleepQualityDetectedEvent) EventType() string {
	return "PoorSleepQualityDetected"
}

func (e *PoorSleepQualityDetectedEvent) Reason() string {
	return e.reason
}

// DaytimeSleepinessChangedEvent - событие изменения дневной сонливости
type DaytimeSleepinessChangedEvent struct {
	sleepEntryID  SleepEntryID
	oldSleepiness valueobjects.DaytimeSleepiness
	newSleepiness valueobjects.DaytimeSleepiness
	occurredOn    time.Time
}

func (e *DaytimeSleepinessChangedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *DaytimeSleepinessChangedEvent) EventType() string {
	return "DaytimeSleepinessChanged"
}

// SleepQualityUpdatedEvent - событие обновления качества сна
type SleepQualityUpdatedEvent struct {
	sleepEntryID SleepEntryID
	oldQuality   valueobjects.SleepQuality
	newQuality   valueobjects.SleepQuality
	occurredOn   time.Time
}

func (e *SleepQualityUpdatedEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *SleepQualityUpdatedEvent) EventType() string {
	return "SleepQualityUpdated"
}
