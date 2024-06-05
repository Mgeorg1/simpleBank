package api

import (
	"github.com/Mgeorg1/simpleBank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
