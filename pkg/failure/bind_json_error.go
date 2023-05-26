package failure

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BindJSONErr struct {
	Code        ErrCode
	OriginalErr error
	Model       reflect.Type
}

func (e BindJSONErr) Error() string {
	var ve validator.ValidationErrors
	if !errors.As(e.OriginalErr, &ve) {
		return "Thông tin không hợp lệ"
	}

	var errs []string
	for _, fe := range ve {
		errs = append(errs, e.translateFieldError(fe))
	}
	return strings.Join(errs, "\n")
}

func (e BindJSONErr) translateFieldError(fe validator.FieldError) string {
	fieldName := getFieldName(e.Model, fe.Field())
	var err string
	switch fe.Tag() {
	case "required":
		err = fmt.Sprintf("%s bị thiếu", fieldName)
	case "numeric":
		err = fmt.Sprintf("%s phải là định dạng số", fieldName)
	case "email":
		err = fmt.Sprintf("%s phải là định dạng email", fieldName)
	case "max":
		err = fmt.Sprintf("%s không thể lớn hơn %s", fieldName, fe.Param())
	case "min":
		err = fmt.Sprintf("%s không thể nhỏ hơn %s", fieldName, fe.Param())
	default:
		err = fmt.Sprintf("%s không hợp lệ", fieldName)
	}
	return err
}

func getFieldName(model reflect.Type, fieldName string) string {
	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		if fieldName == field.Name {
			translatedTag := field.Tag.Get("translated")
			if translatedTag == "" {
				return field.Name
			}
			return translatedTag
		}
	}
	return fieldName
}
