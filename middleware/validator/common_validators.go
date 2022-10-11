package validator

import (
	"github.com/go-playground/validator/v10"
	"math/big"
)

var ValidBN validator.Func = func(fl validator.FieldLevel) bool {
	_, ok := new(big.Int).SetString(fl.Field().Interface().(string), 10)
	return ok
}
