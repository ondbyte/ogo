package ogo

import (
	"fmt"
	"strings"
)

// int64
type Int64PathParam int64

var i Int64PathParam

var int64PathParam = strings.Split(fmt.Sprintf("%T", i), ".")[1]

func ToInt64PathParam(v string) (Int64PathParam, error) {
	return 0, nil
}

// float64
type Float64PathParam float64

var f Float64PathParam

var float64PathParam = strings.Split(fmt.Sprintf("%T", f), ".")[1]

// bool
type BoolPathParam bool

var b BoolPathParam

var boolPathParam = strings.Split(fmt.Sprintf("%T", b), ".")[1]
