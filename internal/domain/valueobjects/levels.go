package valueobjects

import (
	"daily-tracker/pkg/errors"
	"fmt"
	"strings"
)

// Value Objects в DDD - неизменяемые объекты без идентичности
// В Go используем типы с валидацией для обеспечения инвариантов

// StressLevel представляет уровень стресса от 0 до 10
type StressLevel int

const (
	StressLevelMin = 0
	StressLevelMax = 10
)

// NewStressLevel конструктор с валидацией
func NewStressLevel(level int) (StressLevel, error) {
	if level < StressLevelMin || level > StressLevelMax {
		return 0, errors.NewDomainError("stress level must be between 0 and 10")
	}
	return StressLevel(level), nil
}

// Int возвращает значение как int
func (sl StressLevel) Int() int {
	return int(sl)
}

// String реализует интерфейс fmt.Stringer (аналог __toString() в PHP)
func (sl StressLevel) String() string {
	return fmt.Sprintf("%d", sl)
}

// IsHigh проверяет, является ли уровень стресса высоким
func (sl StressLevel) IsHigh() bool {
	return sl >= 7
}

// EnergyLevel представляет уровень энергии от 0 до 10
type EnergyLevel int

func NewEnergyLevel(level int) (EnergyLevel, error) {
	if level < 0 || level > 10 {
		return 0, errors.NewDomainError("energy level must be between 0 and 10")
	}
	return EnergyLevel(level), nil
}

func (el EnergyLevel) Int() int {
	return int(el)
}

func (el EnergyLevel) String() string {
	return fmt.Sprintf("%d", el)
}

func (el EnergyLevel) IsLow() bool {
	return el <= 3
}

// MoodLevel представляет уровень настроения от 0 до 10
type MoodLevel int

func NewMoodLevel(level int) (MoodLevel, error) {
	if level < 0 || level > 10 {
		return 0, errors.NewDomainError("mood level must be between 0 and 10")
	}
	return MoodLevel(level), nil
}

func (ml MoodLevel) Int() int {
	return int(ml)
}

func (ml MoodLevel) String() string {
	return fmt.Sprintf("%d", ml)
}

func (ml MoodLevel) IsPositive() bool {
	return ml >= 6
}

// TaskCategory представляет категорию задачи
type TaskCategory string

// Предопределенные категории (enum в стиле Go)
const (
	TaskCategoryWork     TaskCategory = "работа"
	TaskCategoryStudy    TaskCategory = "учеба"
	TaskCategoryPersonal TaskCategory = "личное"
	TaskCategoryHealth   TaskCategory = "здоровье"
	TaskCategoryHobbies  TaskCategory = "хобби"
	TaskCategoryOther    TaskCategory = "другое"
)

// AllTaskCategories возвращает список всех доступных категорий
func AllTaskCategories() []TaskCategory {
	return []TaskCategory{
		TaskCategoryWork,
		TaskCategoryStudy,
		TaskCategoryPersonal,
		TaskCategoryHealth,
		TaskCategoryHobbies,
		TaskCategoryOther,
	}
}

// NewTaskCategory конструктор с валидацией
func NewTaskCategory(category string) (TaskCategory, error) {
	// Приводим к нижнему регистру для сравнения
	category = strings.ToLower(strings.TrimSpace(category))

	for _, validCategory := range AllTaskCategories() {
		if strings.ToLower(string(validCategory)) == category {
			return validCategory, nil
		}
	}

	return "", errors.NewDomainError("invalid task category: " + category)
}

func (tc TaskCategory) String() string {
	return string(tc)
}

func (tc TaskCategory) IsValid() bool {
	for _, validCategory := range AllTaskCategories() {
		if tc == validCategory {
			return true
		}
	}
	return false
}

// SleepQuality представляет качество сна от 0 до 10
type SleepQuality int

func NewSleepQuality(quality int) (SleepQuality, error) {
	if quality < 0 || quality > 10 {
		return 0, errors.NewDomainError("sleep quality must be between 0 and 10")
	}
	return SleepQuality(quality), nil
}

func (sq SleepQuality) Int() int {
	return int(sq)
}

func (sq SleepQuality) String() string {
	return fmt.Sprintf("%d", sq)
}

func (sq SleepQuality) IsGood() bool {
	return sq >= 7
}

// DaytimeSleepiness представляет дневную сонливость от 0 до 10
type DaytimeSleepiness int

func NewDaytimeSleepiness(sleepiness int) (DaytimeSleepiness, error) {
	if sleepiness < 0 || sleepiness > 10 {
		return 0, errors.NewDomainError("daytime sleepiness must be between 0 and 10")
	}
	return DaytimeSleepiness(sleepiness), nil
}

func (ds DaytimeSleepiness) Int() int {
	return int(ds)
}

func (ds DaytimeSleepiness) String() string {
	return fmt.Sprintf("%d", ds)
}

func (ds DaytimeSleepiness) IsHigh() bool {
	return ds >= 7
}
