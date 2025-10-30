package marshall

import (
	"errors"
	"reflect"
	"strings"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	ErrorNotPointer = errors.New("source and destination must be pointer")
	ErrorNotStruct  = errors.New("source and destination must be point to struct")
)

// MarshalModels marshalling from models to protocol buffers json struct response
func MarshalModels[T any, K any](sourceObj *T, destObj *K) (*K, error) {
	if sourceObj == nil || destObj == nil {
		return nil, ErrorNotPointer
	}

	sourceValue := reflect.ValueOf(sourceObj).Elem()
	destValue := reflect.ValueOf(destObj).Elem()

	if sourceValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Struct {
		return nil, ErrorNotStruct
	}

	sourceType := sourceValue.Type()
	destType := destValue.Type()

	destFields := make(map[string]int)
	for i := 0; i < destValue.NumField(); i++ {
		fieldName := normalizeFieldName(destType.Field(i).Name)
		destFields[fieldName] = i
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Field(i)
		sourceFieldName := normalizeFieldName(sourceType.Field(i).Name)

		if destFieldIndex, exists := destFields[sourceFieldName]; exists {
			destField := destValue.Field(destFieldIndex)

			if err := setFieldValue(sourceField, destField); err != nil {
				return nil, err
			}
		}
	}

	return destObj, nil
}

func setFieldValue(sourceField, destField reflect.Value) error {
	if !sourceField.IsValid() || !destField.IsValid() {
		return nil
	}

	if isZeroValue(sourceField) && isOptionalProtoField(destField) {
		return nil
	}

	switch destField.Type().String() {
	case "*wrapperspb.UInt64Value":
		if sourceField.Type().Kind() == reflect.Uint64 {
			if !isZeroValue(sourceField) {
				value := wrapperspb.UInt64(sourceField.Uint())
				destField.Set(reflect.ValueOf(value))
			}
		}
		return nil
	case "*wrapperspb.StringValue":
		if sourceField.Type().Kind() == reflect.String {
			if !isZeroValue(sourceField) {
				value := wrapperspb.String(sourceField.String())
				destField.Set(reflect.ValueOf(value))
			}
		}
		return nil
	case "*wrapperspb.UInt32Value":
		if sourceField.Type().Kind() == reflect.Uint32 {
			if !isZeroValue(sourceField) {
				value := wrapperspb.UInt32(uint32(sourceField.Uint()))
				destField.Set(reflect.ValueOf(value))
			}
		}
		return nil
	case "*wrapperspb.FloatValue":
		if sourceField.Type().Kind() == reflect.Float32 {
			if !isZeroValue(sourceField) {
				value := wrapperspb.Float(float32(sourceField.Float()))
				destField.Set(reflect.ValueOf(value))
			}
		}
		return nil
	case "*wrapperspb.":
	}

	if sourceField.Type().AssignableTo(destField.Type()) {
		destField.Set(sourceField)
	} else if sourceField.Type().ConvertibleTo(destField.Type()) {
		destField.Set(sourceField.Convert(destField.Type()))
	}

	return nil
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Array:
		return v.IsNil()
	default:
		return false
	}
}

func isOptionalProtoField(v reflect.Value) bool {
	fieldType := v.Type().String()
	return strings.Contains(fieldType, "*wrapperspb.") || v.Kind() == reflect.Ptr
}

func normalizeFieldName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, "_", ""))
}
