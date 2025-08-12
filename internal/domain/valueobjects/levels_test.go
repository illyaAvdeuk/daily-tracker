package valueobjects

import (
	"fmt"
	"testing"
)

// Тестируем StressLevel
func TestNewStressLevel_Valid(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"minimum level", 0},
		{"medium level", 5},
		{"maximum level", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl, err := NewStressLevel(tt.level)
			if err != nil {
				t.Errorf("Expected no error for level %d, got: %v", tt.level, err)
			}

			if sl.Int() != tt.level {
				t.Errorf("Expected level %d, got %d", tt.level, sl.Int())
			}
		})
	}
}

func TestNewStressLevel_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"below minimum", -1},
		{"above maximum", 11},
		{"way too high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStressLevel(tt.level)
			if err == nil {
				t.Errorf("Expected error for invalid level %d, got nil", tt.level)
			}
		})
	}
}

func TestStressLevel_IsHigh(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected bool
	}{
		{"low stress", 3, false},
		{"medium stress", 6, false},
		{"high stress boundary", 7, true},
		{"very high stress", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl, _ := NewStressLevel(tt.level)
			if sl.IsHigh() != tt.expected {
				t.Errorf("Expected IsHigh() = %v for level %d, got %v",
					tt.expected, tt.level, sl.IsHigh())
			}
		})
	}
}

// Тестируем TaskCategory
func TestNewTaskCategory_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TaskCategory
	}{
		{"work category", "работа", TaskCategoryWork},
		{"work category uppercase", "РАБОТА", TaskCategoryWork},
		{"work category with spaces", " работа ", TaskCategoryWork},
		{"study category", "учеба", TaskCategoryStudy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := NewTaskCategory(tt.input)
			if err != nil {
				t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
			}

			if category != tt.expected {
				t.Errorf("Expected category %s, got %s", tt.expected, category)
			}
		})
	}
}

func TestNewTaskCategory_Invalid(t *testing.T) {
	invalidInputs := []string{
		"invalid_category",
		"спорт", // не входит в предопределенные категории
		"",
	}

	for _, input := range invalidInputs {
		t.Run("invalid: "+input, func(t *testing.T) {
			_, err := NewTaskCategory(input)
			if err == nil {
				t.Errorf("Expected error for invalid input '%s', got nil", input)
			}
		})
	}
}

func TestTaskCategory_IsValid(t *testing.T) {
	// Тестируем валидные категории
	validCategories := AllTaskCategories()
	for _, category := range validCategories {
		if !category.IsValid() {
			t.Errorf("Category %s should be valid", category)
		}
	}

	// Тестируем невалидную категорию
	invalidCategory := TaskCategory("несуществующая_категория")
	if invalidCategory.IsValid() {
		t.Error("Invalid category should not be valid")
	}
}

// Пример тестирования с подготовкой и очисткой (setup/teardown)
func TestAllTaskCategories(t *testing.T) {
	categories := AllTaskCategories()

	// Проверяем, что возвращается ожидаемое количество категорий
	expectedCount := 6
	if len(categories) != expectedCount {
		t.Errorf("Expected %d categories, got %d", expectedCount, len(categories))
	}

	// Проверяем, что нет дубликатов
	seen := make(map[TaskCategory]bool)
	for _, category := range categories {
		if seen[category] {
			t.Errorf("Duplicate category found: %s", category)
		}
		seen[category] = true
	}

	// Проверяем, что все категории валидны
	for _, category := range categories {
		if !category.IsValid() {
			t.Errorf("Category %s should be valid", category)
		}
	}
}

// Бенчмарк для проверки производительности создания категории
func BenchmarkNewTaskCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewTaskCategory("работа")
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

// Пример сравнительного бенчмарка
func BenchmarkTaskCategory_IsValid_Valid(b *testing.B) {
	category := TaskCategoryWork
	b.ResetTimer() // Не учитываем время подготовки

	for i := 0; i < b.N; i++ {
		category.IsValid()
	}
}

func BenchmarkTaskCategory_IsValid_Invalid(b *testing.B) {
	category := TaskCategory("invalid")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		category.IsValid()
	}
}

// Пример теста с параллельным выполнением
func TestStressLevel_String_Parallel(t *testing.T) {
	t.Parallel() // Этот тест может выполняться параллельно

	levels := []int{0, 5, 10}

	for _, level := range levels {
		level := level // Важно: копируем переменную для горутин
		t.Run("level_"+string(rune(level+'0')), func(t *testing.T) {
			t.Parallel() // И подтесты тоже параллельно

			sl, err := NewStressLevel(level)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			expected := fmt.Sprintf("%d", level)
			if sl.String() != expected {
				t.Errorf("Expected string %s, got %s", expected, sl.String())
			}
		})
	}
}

// Пример использования testify (если добавим зависимость)
// func TestWithTestify(t *testing.T) {
//     assert := assert.New(t)
//     sl, err := NewStressLevel(5)
//
//     assert.NoError(err)
//     assert.Equal(5, sl.Int())
//     assert.Equal("5", sl.String())
//     assert.False(sl.IsHigh())
// }
