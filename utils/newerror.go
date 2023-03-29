package utils

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Message string `json:"message"`
}

func ReturnMess(err string) *ErrorResponse {
	return &ErrorResponse{
		Message: err,
	}
}

func ValidateError(err validator.ValidationErrors) string {
	var message string
	for _, e := range err {
		switch e.Tag() {
		case "required":
			message = e.Field() + " alanı zorunludur"
		case "email":
			message = "Geçerli bir email adresi giriniz"
		}
	}
	return message
}

func ValidateUpdateError(err validator.ValidationErrors) string {
	var message string
	for _, e := range err {
		switch e.Tag() {
		case "email":
			message = "Geçerli bir email adresi giriniz"
		}
	}
	return message
}
