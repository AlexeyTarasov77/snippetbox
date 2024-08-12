package forms

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	validatorErrs "snippetbox.proj.net/internal/lib/validator"
	"snippetbox.proj.net/internal/utils"
)

type BaseForm struct {
	NonFieldErrors []string `schema:"-" validate:"omitempty"`
	FieldErrors map[string]string `schema:"-" validate:"omitempty"`
}

func (bf *BaseForm) Validate(form any) {
	validate := validator.New()
	bf.FieldErrors = make(map[string]string)
	if err := validate.Struct(form); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := utils.ConvertCamelToSnake(e.StructField())
			bf.FieldErrors[fieldName] = validatorErrs.GetErrorMsgForField(form, e)
		}
	}
}


// Checks whether the form is valid.
// If false was returned, form.FieldErrors will be prepolulated with validation errors.
func (bf *BaseForm) IsValid(form any) bool {
	bf.Validate(form)
	return len(bf.FieldErrors) == 0
}

func IsRequiredField(form any, fieldName string) bool {
	field, found := reflect.TypeOf(form).FieldByName(fieldName)
	if !found {
		panic(fmt.Sprintf("Field %s not found in type %s", fieldName, reflect.TypeOf(form).Name()))
	}
	return strings.Contains(field.Tag.Get("validate"), "required")
}
