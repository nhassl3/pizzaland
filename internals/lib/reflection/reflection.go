package reflection

import (
	"log/slog"
	"reflect"
)

func AllFieldsIsNil(msg any) bool {
	v := reflect.ValueOf(msg).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		slog.Info("Message:", slog.Any("A", f))
		if !f.IsNil() {
			return false
		}
	}
	return true
}
