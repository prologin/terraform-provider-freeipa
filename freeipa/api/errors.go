package api

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v (%v): %v", e.Name, e.Code, e.Message)
}
