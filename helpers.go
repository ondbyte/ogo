package ogo

import (
	"fmt"
	"reflect"
)

func isBasicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.Bool, reflect.TypeOf([]byte("")).Kind():
		return true
	default:
		return false
	}
}

func isValidRespBodyType(t reflect.Type) (err error, isBasic bool) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if isBasicType(t) {
		return nil, true
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("'%v' is not valid type for response body", t), false
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			return fmt.Errorf("all fields in response body type '%v' should be exported", t.Name()), false
		}
	}
	return nil, false
}
