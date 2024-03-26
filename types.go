package ogo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Param[T any] struct {
	Val T
}

type PathParam[T any] Param[T]
type QueryParam[T any] Param[T]
type Header[T any] Param[T]
type FormField[T any] Param[T]
type PostFormField[T any] Param[T]

// Implement the set method for PathParam[int64].
func setParam(pp any, val string) (err error) {
	if val == "" {
		return nil
	}
	switch pp := pp.(type) {
	case *PathParam[int64]:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		pp.Val = i
	case *QueryParam[int64]:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		pp.Val = i
	case *Header[int64]:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		pp.Val = i
	case *FormField[int64]:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		pp.Val = i
	case *PostFormField[int64]:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		pp.Val = i
	default:
		return fmt.Errorf("type '%T' is not supported for PathParam", pp)
	}
	return nil
}

type Body[T any] struct {
	Val T
}

// Unmarshal implements BodyUnmarshaller.
func (b *Body[T]) Unmarshal(bs []byte, v any) error {
	return json.Unmarshal(bs, &b.Val)
}

type BodyUnmarshaller interface {
	Unmarshal(bs []byte, v any) error
}

var _ BodyUnmarshaller = &Body[any]{}

// returns only name part of the type of any types
// ex: Body from Body[any]
func getTypeName(t reflect.Type) (typeName string) {
	a, _ := getPkgAndTypeAndTName(t)
	return a
}

func getPkgAndTypeAndTName(t reflect.Type) (typeName, tName string) {
	removed := strings.Replace(t.String(), "*", "", 1)
	splitted := strings.Split(removed, ".")
	splitted2 := strings.Split(splitted[1], "[")
	tName = strings.TrimSuffix(splitted2[1], "]")
	return splitted2[0], tName
}

var ourPkgName = reflect.TypeOf(&PathParam[any]{}).PkgPath()
