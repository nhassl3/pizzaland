package reflection

import "reflect"

// AllFieldsIsNil returns true if all exported fields in a struct are nil or zero.
func AllFieldsIsNil(msg any) bool {
	v := reflect.ValueOf(msg)

	// Dereference pointers
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		// Skip unexported fields (protobuf internals)
		if !f.CanInterface() {
			continue
		}

		val := f.Interface()
		switch f.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if !f.IsNil() {
				return false
			}
		case reflect.Struct:
			// Check nested structs recursively
			if !AllFieldsIsNil(val) {
				return false
			}
		default:
			if !reflect.DeepEqual(val, reflect.Zero(f.Type()).Interface()) {
				return false
			}
		}
	}

	return true
}
