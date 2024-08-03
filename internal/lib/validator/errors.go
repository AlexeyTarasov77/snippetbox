package validator

import (
	"fmt"
	"reflect"
	go_validator "github.com/go-playground/validator/v10"
)

// func ProcessErrorMsgs(obj any, tags []string) (processedErrors map[string]string) {
// 	t := reflect.TypeOf(obj)
// 	processedErrors = make(map[string]string)
// 	for i := 0; i < t.NumField(); i++ {
// 		fieldName := t.Field(i).Name
// 		errorMsg := GetErrorMsgForField(obj, fieldName, tags[i])
// 		processedErrors[fieldName] = errorMsg
// 	}
// 	return processedErrors
// }

func GetErrorMsgForField(obj any, err go_validator.FieldError) (errorMsg string) {
	t := reflect.TypeOf(obj)
	field, found := t.FieldByName(err.StructField())
	if !found {
		panic(fmt.Sprintf("Field %s not found in type %s", err.StructField(), t.Name()))
	}
	errorMsg = field.Tag.Get("errorMsg")
	if errorMsg == "" {
		switch err.Tag() {
		case "required":
			errorMsg = "This field is required"
		case "max":
			errorMsg = fmt.Sprintf("The maximum value is %s", err.Param())
		case "min":
			errorMsg = fmt.Sprintf("The minimum value is %s", err.Param())
		case "gte":
			errorMsg = fmt.Sprintf("Value should be greater than or equal to %s", err.Param())
		case "lte":
			errorMsg = fmt.Sprintf("Value should be less than or equal to %s", err.Param())
		default:
			errorMsg = "This field is invalid"
		}
	}
	return
}

// func getValidationParam(param string)