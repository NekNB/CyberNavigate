package errors

type CustomError struct {
	Type error
	Msg  string
}

func (e *CustomError) Error() string {
	return e.Msg
}

func (e *CustomError) Is(target error) bool {
	return target == e.Type
}

func NewTypedError(typeError error, errorMsg string) error {

	if errorMsg == "" {
		errorMsg = typeError.Error()
	}

	return &CustomError{Type: typeError, Msg: errorMsg}
}
