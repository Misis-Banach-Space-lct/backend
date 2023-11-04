package model

import "fmt"

type ErrFieldNoValue struct {
	Field string
}

func (e ErrFieldNoValue) Error() string {
	return fmt.Sprintf("no value for field %s found in .env", e.Field)
}

type ErrPostgresCreateFailed struct {
	Message string
}

func (e ErrPostgresCreateFailed) Error() string {
	return fmt.Sprintf("can't create postgres instance: %s", e.Message)
}
