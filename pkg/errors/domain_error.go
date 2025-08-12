package errors

import "fmt"

// DomainError представляет ошибку на уровне домена
// В Go ошибки - это значения, а не исключения как в PHP
type DomainError struct {
	message string
	code    string
}

// Error реализует интерфейс error (встроенный в Go)
func (de *DomainError) Error() string {
	return de.message
}

// Code возвращает код ошибки
func (de *DomainError) Code() string {
	return de.code
}

// Message возвращает сообщение ошибки
func (de *DomainError) Message() string {
	return de.message
}

// NewDomainError создает новую доменную ошибку
func NewDomainError(message string) *DomainError {
	return &DomainError{
		message: message,
		code:    "DOMAIN_ERROR",
	}
}

// NewDomainErrorWithCode создает доменную ошибку с кодом
func NewDomainErrorWithCode(message, code string) *DomainError {
	return &DomainError{
		message: message,
		code:    code,
	}
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	field   string
	message string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", ve.field, ve.message)
}

func (ve *ValidationError) Field() string {
	return ve.field
}

func (ve *ValidationError) Message() string {
	return ve.message
}

// NewValidationError создает ошибку валидации
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		field:   field,
		message: message,
	}
}

// NotFoundError представляет ошибку "не найдено"
type NotFoundError struct {
	resource string
	id       string
}

func (nfe *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", nfe.resource, nfe.id)
}

func (nfe *NotFoundError) Resource() string {
	return nfe.resource
}

func (nfe *NotFoundError) ID() string {
	return nfe.id
}

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		resource: resource,
		id:       id,
	}
}

// IsDomainError проверяет, является ли ошибка доменной
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// IsValidationError проверяет, является ли ошибка валидационной
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsNotFoundError проверяет, является ли ошибка "не найдено"
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
