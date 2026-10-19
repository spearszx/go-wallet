package apperrors

import "errors"

type Category int

const (
	CategoryUnknown Category = iota
	CategoryBadRequest
	CategoryNotFound
	CategoryConflict
	CategoryInternal
)

type Categorized interface {
	error
	Category() Category
}

type AppError struct {
	category Category
	message  string
	err      error
}

func New(category Category, message string) *AppError {
	return &AppError{category: category, message: message}
}

func Wrap(category Category, message string, err error) *AppError {
	return &AppError{category: category, message: message, err: err}
}

func (e *AppError) Category() Category { return e.category }
func (e *AppError) Error() string      { return e.message }
func (e *AppError) Unwrap() error      { return e.err }

func IsCategory(err error, c Category) bool {
	var cat Categorized
	return errors.As(err, &cat) && cat.Category() == c
}
