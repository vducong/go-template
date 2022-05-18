package failure

import (
	"reflect"
)

type BindJSONErr struct {
	Model reflect.Type
	Err   error
}

func (e BindJSONErr) String() string {
	return TranslateBindingJSONError(e.Model, e.Err)
}
