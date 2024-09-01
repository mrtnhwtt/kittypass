package kittypass

import "fmt"

type MalformedDataError struct {
	Data string
}

func (e MalformedDataError) Error() string {
	return fmt.Sprintf("data for %s is malformed, could not be processed", e.Data)
}

type IncorrectPasswordError struct {}

func (e IncorrectPasswordError) Error() string {
	return "incorrect password"
}