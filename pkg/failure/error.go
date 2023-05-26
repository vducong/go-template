package failure

import (
	"errors"
	"fmt"
	"promotion/pkg/tracing"
	"strings"

	"gorm.io/gorm"
)

func IsSQLRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsFSNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "code = NotFound")
}

func ErrWithTrace(err error) error {
	return fmt.Errorf("%w \n at %s", err, tracing.Trace(4))
}
