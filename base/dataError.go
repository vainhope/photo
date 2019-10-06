package base

import "fmt"

type DataError struct {
	Code    int
	Message string
}

func (data *DataError) Error() string {
	return fmt.Sprintf("code: %d ,error message:%s", data.Code, data.Message)
}

func Err(code int, message string) error {
	return &DataError{Code: code, Message: message}
}
