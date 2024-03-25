package ogo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Mux struct {
	mux *http.ServeMux
}

func NewOgoMux(mux *http.ServeMux) *Mux {
	return &Mux{
		mux: mux,
	}
}

type Err struct {
	msg string
}

func (e *Err) Error() string {
	return e.msg
}

func NewErr(format string, a ...any) error {
	return &Err{
		msg: fmt.Sprintf(format, a...),
	}
}

func (s *Mux) Handle(path string, handlerFunc any) {
	var err error
	logger := NewHandlerLogger(path)
	pathParams := PathParams(path)
	handler := reflect.TypeOf(handlerFunc)
	if handler.NumIn() != 1 {
		logger.logAndPanic("handler should have exactly one struct argument")
	}
	firstArg := handler.In(0).Elem()
	if firstArg.Kind() != reflect.Struct {
		logger.logAndPanic("first argument type of handler should be a pointer struct")
	}
	fields := firstArg.NumField()
	ptrFields := make([]bool, fields)
	structFields := make([]*reflect.StructField, fields)
	fieldSetters := make([]ReflValueExtractor, fields)
	for j := range fields {
		field := firstArg.Field(j)
		if !field.IsExported() {
			logger.logAndPanic(
				"all fields of the handler argument struct needs to be exported, ie must start with a capital letter change '%v' to something like '%v'",
				field.Name,
				strings.ToTitle(field.Name),
			)
		}
		structFields[j] = &field
		fieldSetters[j], err = GetReflValueExtractor(field.Name, field.Type)
		if err != nil {
			logger.logAndPanic(err.Error())
		}
		ptrFields[j] = field.Type.Kind() == reflect.Pointer
		// mark all path params as consumed one by one
		if pathParams[field.Name] == 1 {
			pathParams[field.Name]++
		}
	}
	//verify whether all path params has been consumed
	for k, v := range pathParams {
		if v != 2 {
			logger.logAndPanic(`path param '%v' isn't consumed
make sure you add a field with name '%v' with any PathParam type to your argument struct
example: '%v %v[int],'`, k, k, k, getTypeName(reflect.TypeOf(&PathParam[any]{})),
			)
		}
	}
	reqToArg := func(r *http.Request) (*reflect.Value, error) {
		newHandlerArgStruct := reflect.New(firstArg)
		structValue := newHandlerArgStruct.Elem()
		for i, setter := range fieldSetters {
			field := structValue.Field(i)
			var err error
			if ptrFields[i] {
				field.Set(reflect.New(structFields[i].Type.Elem()))
				err = setter(r, field.Interface())
			} else {
				addr := field.Addr()
				err = setter(r, addr.Interface())
			}
			if err != nil {
				return nil, fmt.Errorf("unable to set field '%v' of arg struct %v due to err: %v", field, firstArg, err)
			}
		}
		return &newHandlerArgStruct, nil
	}
	rHandlerFn := reflect.ValueOf(handlerFunc)
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		arg, err := reqToArg(r)
		if err != nil {
			logger.log(err.Error())
		}
		rHandlerFn.Call([]reflect.Value{*arg})
	})
}
