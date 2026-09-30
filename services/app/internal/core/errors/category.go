package core_errors

type Category uint8

const (
	CategoryInvalidArgument Category = iota
	CategoryNotFound
	CategoryConflict
	CategoryBadRequest
	CategoryDatabase
)

type CodedError struct {
	Err      error
	Category Category
}

func (e *CodedError) Error() string {
	return e.Err.Error()
}

func (e *CodedError) Unwrap() error {
	return e.Err
}

func NewError(err error, category Category) error {
	return &CodedError{
		Err:      err,
		Category: category,
	}
}
