package ogo

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type paramType int

const (
	PATH paramType = iota
	QUERY
	HEADER
)

type ReflValueExtractor func(r *http.Request, v interface{}) error

func GetReflValueExtractor(paramName string, paramType reflect.Type) (b ReflValueExtractor, err error) {
	isPtr := paramType.Kind() == reflect.Pointer
	isRequired := !isPtr
	paramTypeName, tName := getPkgAndTypeAndTName(paramType)
	if paramType.PkgPath() != ourPkgName {
		return nil, fmt.Errorf("you should use only '%v' provided types in the handler argument struct but you have used %v", ourPkgName, paramType)
	}
	if strings.HasPrefix(tName, "*") {
		return nil, fmt.Errorf("you can only use non pointer type for generic type but you have used '%v'", tName)
	}
	switch paramTypeName {
	case "":
		{

		}
	case getTypeName(reflect.TypeOf(&Body[any]{})): //match body
		{
			b = func(r *http.Request, nw interface{}) error {
				bs, err := io.ReadAll(r.Body)
				if err != nil {
					return NewErr("err while io.ReadAll for body for param name '%v' and type '%v', err: %v", paramName, paramType, err)
				}
				if len(bs) == 0 {

					return nil
				}
				um := nw.(BodyUnmarshaller)
				return um.Unmarshal(bs, um)
			}
			return
		}
	case getTypeName(reflect.TypeOf(&PathParam[any]{})): //match param
		{
			b = func(r *http.Request, nw interface{}) error {
				v := r.PathValue(paramName)
				if v == "" {
					return fmt.Errorf("path param '%v' is required in request", paramName)
				}
				return setParam(nw, v)
			}
			return
		}
	case getTypeName(reflect.TypeOf(&QueryParam[any]{})): //match param
		{
			b = func(r *http.Request, nw interface{}) error {
				q := r.URL.Query()
				v := q.Get(paramName)
				if v == "" && isRequired {
					return fmt.Errorf("query param '%v' is required in request", paramName)
				}
				return setParam(nw, v)
			}
			return
		}
	case getTypeName(reflect.TypeOf(&Header[any]{})): //match param
		{
			b = func(r *http.Request, nw interface{}) error {
				v := r.Header.Get(paramName)
				if v == "" && isRequired {
					return fmt.Errorf("header '%v' is required in request", paramName)
				}
				return setParam(nw, v)
			}
			return
		}
	case getTypeName(reflect.TypeOf(&FormField[any]{})): //match param
		{
			b = func(r *http.Request, nw interface{}) error {
				v := r.FormValue(paramName)
				if v == "" && isRequired {
					return fmt.Errorf("form field '%v' is required in request", paramName)
				}
				return setParam(nw, v)
			}
			return
		}
	case getTypeName(reflect.TypeOf(&PostFormField[any]{})): //match param
		{
			b = func(r *http.Request, nw interface{}) error {
				v := r.PostFormValue(paramName)
				if v == "" && isRequired {
					return fmt.Errorf("post form field '%v' is required in request", paramName)
				}
				return setParam(nw, v)
			}
			return
		}
	}
	return nil, fmt.Errorf("type %v isn't supported for param name %v", paramType, paramName)
}
