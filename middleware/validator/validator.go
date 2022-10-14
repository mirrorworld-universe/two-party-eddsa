package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidators() {
	// custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validbn", ValidBN)
	}
}
