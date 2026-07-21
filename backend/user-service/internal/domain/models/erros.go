package models

type CustomError struct {
	Msg string
}

func (c *CustomError) Error() string {
	return c.Msg
}
