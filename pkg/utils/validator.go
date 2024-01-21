package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func GetValidateErr(obj any, rawErr error) error {
	validationErrs, ok := rawErr.(validator.ValidationErrors)
	if !ok {
		return rawErr
	}
	var errString []string
	for _, validationErr := range validationErrs {
		field, ok := reflect.TypeOf(obj).FieldByName(validationErr.Field())
		if ok {
			if e := field.Tag.Get("msg"); e != "" {
				errString = append(errString, e)
				continue
			}
		}
		errString = append(errString, validationErr.Error())
	}
	return errors.New(strings.Join(errString, "\n"))
}
