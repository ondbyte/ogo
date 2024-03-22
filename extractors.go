package ogo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type paramType int

const (
	PATH paramType = iota
	QUERY
	HEADER
)

type ValueExtractor func(r *http.Request) (interface{}, error)
type StrValueExtractor func(r *http.Request) (string, error)

func GetValueExtractor(paramName string, paramType reflect.Type) (b ValueExtractor, err error) {
	switch paramType.String() {
	case "body", "Body":
		{
			b = func(r *http.Request) (interface{}, error) {
				bs, err := io.ReadAll(r.Body)
				if err != nil {
					return nil, fmt.Errorf("err while reading body for param name '%v' and type '%v', err: %v", paramName, paramType, err)
				}
				n := reflect.New(paramType).Interface()
				err = json.Unmarshal(bs, n)
				if err != nil {
					return nil, fmt.Errorf("error while unmarshalling the body to pram name '%v' and type '%v', err: %v", paramName, paramType, err)
				}
				return n, nil
			}
			return
		}
	case int64PathParam:
		b = func(r *http.Request) (interface{}, error) {
			v := r.PathValue(paramName)
			return strconv.ParseInt(v, 10, 64)
		}
		return
	default:
		return nil, fmt.Errorf("type %v isn't supported for param name %v", paramType, paramName)
	}
}

func GetStrValueExtractor(paramName string, paramType reflect.Type) (b StrValueExtractor, err error) {
	isPtr := paramType.Kind() == reflect.Pointer
	pt := strings.Split(paramType.String(), ".")[1]
	switch pt {
	case "body", "Body":
		{
			b = func(r *http.Request) (string, error) {
				bs, err := io.ReadAll(r.Body)
				if err != nil {
					return "", fmt.Errorf("err while reading body for param name '%v' and type '%v', err: %v", paramName, paramType, err)
				}
				if len(bs) == 0 {
					return "null", nil
				}
				return (string(bs)), nil
			}
			return
		}
	case int64PathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (string, error) {
			v := r.PathValue(paramName)
			return v, nil
		}
		return
	case float64PathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (string, error) {
			v := r.PathValue(paramName)
			return v, nil
		}
		return
	case boolPathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (string, error) {
			v := r.PathValue(paramName)
			return v, nil
		}
		return
	default:
		return nil, fmt.Errorf("type %v isn't supported for param name %v", paramType, paramName)
	}
}

type ReflValueExtractor func(r *http.Request) (*reflect.Value, error)

func GetReflValueExtractor(paramName string, paramType reflect.Type) (b ReflValueExtractor, err error) {
	isPtr := paramType.Kind() == reflect.Pointer
	pt := strings.Split(paramType.String(), ".")[1]
	switch pt {
	case "body", "Body":
		{
			b = func(r *http.Request) (*reflect.Value, error) {
				bs, err := io.ReadAll(r.Body)
				if err != nil {
					return nil, fmt.Errorf("err while reading body for param name '%v' and type '%v', err: %v", paramName, paramType, err)
				}
				if len(bs) == 0 {
					return nil, nil
				}
				nw := reflect.New(paramType)
				err = json.Unmarshal(bs, nw.Interface())
				if err != nil {
					return nil, fmt.Errorf("unable to unmarshal to your '%v' field of type '%v' due to err: %v", paramName, paramType, err)
				}
				return &nw, nil
			}
			return
		}
	case int64PathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (*reflect.Value, error) {
			v := r.PathValue(paramName)
			if v==""{
				return nil,fmt.Errorf("")
			}
			return v, nil
		}
		return
	case float64PathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (string, error) {
			v := r.PathValue(paramName)
			return v, nil
		}
		return
	case boolPathParam:
		if isPtr {
			panic(NewErr(fmt.Sprintf("type of '%v' cannot be '%v' ie a pointer, because a path parameter will always be present", paramName, paramType)))
		}
		b = func(r *http.Request) (string, error) {
			v := r.PathValue(paramName)
			return v, nil
		}
		return
	default:
		return nil, fmt.Errorf("type %v isn't supported for param name %v", paramType, paramName)
	}
}
