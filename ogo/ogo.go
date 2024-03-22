package ogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Demo struct {
	Name string
	Age  uint
}

type Param[T any] struct {
	Val T
}

func (pp *Param[T]) set(a any) (b bool) {
	pp.Val, b = a.(T)
	return
}

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
	defer func() {
		err := recover()
		if err2, ok := err.(*Err); ok {
			panic(fmt.Sprintf("OGO: Handle: %v\n%v", path, err2.Error()))
		} else if err != nil {
			panic(err)
		}
	}()
	pathParams := PathParams(path)
	handlerVal := reflect.ValueOf(handlerFunc)
	handler := reflect.TypeOf(handlerFunc)
	if handler.NumIn() != 1 {
		panic(NewErr("handler should have exactly one struct argument"))
	}
	firstArg := handler.In(0).Elem()
	if firstArg.Kind() != reflect.Struct {
		panic(NewErr("first argument type of handler should be a pointer struct"))
	}
	jsonArg := &strings.Builder{}
	jsonArg.WriteRune('{')
	jsonPlaceHolders := map[string]StrValueExtractor{}
	fields := firstArg.NumField()
	for j := range fields {
		field := firstArg.Field(j)
		//val := r.PathValue(field.Name)
		strExtr, err := GetStrValueExtractor(field.Name, field.Type)
		if err != nil {
			panic(NewErr(err.Error()))
		}
		jsonArg.WriteString(strconv.Quote(field.Name))
		jsonArg.WriteRune(':')
		placeHolder := "$" + field.Name
		jsonArg.WriteString(placeHolder)
		if j < fields-1 {
			jsonArg.WriteRune(',')
		}
		jsonPlaceHolders[placeHolder] = strExtr

		// mark all path params as consumed one by one
		if pathParams[field.Name] == 1 {
			pathParams[field.Name]++
		}
	}
	//verify whether all path params has been consumed
	for k, v := range pathParams {
		if v != 2 {
			panic(NewErr(`path param '%v' isn't consumed,
make sure you add a field with name '%v' with any PathParam type to your argument struct
example: '%v %v,'`, k, k, k, int64PathParam,
			))
		}
	}
	jsonArg.WriteRune('}')
	reqToArg := func(r *http.Request) (*reflect.Value, error) {
		newHandlerArg := reflect.New(firstArg)
		jsonArg := jsonArg.String()
		for k, v := range jsonPlaceHolders {
			val, err := v(r)
			if err != nil {
				return nil, fmt.Errorf("err while constructing struct argument of type '%v' for handler path '%v', err: %v", firstArg, path, err)
			}
			jsonArg = strings.Replace(jsonArg, k, val, 1)
		}
		err := json.Unmarshal([]byte(jsonArg), newHandlerArg.Interface())
		if err != nil {
			return nil, fmt.Errorf(`err while constructing struct argument of type '%v' for handler path '%v'
error while unmarshalling json '%v'
err: %v`, firstArg, path, jsonArg, err)
		}
		return &newHandlerArg, nil
	}

	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		arg, err := reqToArg(r)
		if err != nil {
			fmt.Println(err, arg)
		}
		handlerVal.Call([]reflect.Value{*arg})
	})
}
