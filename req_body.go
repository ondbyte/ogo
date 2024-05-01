package ogo

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

type RequestBody struct {
	mediaType        mediaType
	validationStatus int
	validationErr    string
	descr            string
	mapper           func(bs []byte, ptr any)
}

func (body *RequestBody) MediaType(mt mediaType) *RequestBody {
	body.mediaType = mt
	return body
}

func (body *RequestBody) Description(d string) *RequestBody {
	body.descr = d
	return body
}

// makes this body required, returned response will be the passed status and err when body is empty
func (body *RequestBody) Required(status int, err string) *RequestBody {
	body.validationStatus = status
	body.validationErr = err
	return body
}

type BodySettings func(body *RequestBody)

// map the body of the request to a ptr
// also set the details of the body using BodySettings
func (v *RequestValidator[reqData, respData]) Body(ptr any, s BodySettings) {
	if v.ogo {
		if v.method == "GET" || v.method == "HEAD" {
			panic(fmt.Sprintf("GET/HEAD method cannot have a request body but handler for path '%v' does", v.path))
		}
		t := reflect.TypeOf(ptr)
		if t.Kind() != reflect.Pointer {
			panic("ptr should be a pointer")
		}
		if reflect.ValueOf(ptr).IsNil() {
			panic(`ptr cannot be empty, the pointer needs to be initialized to its zero value but its nil.
example :
ptr:=&YourBody{}
v.Body(ptr,....)
`)
		}
		body := &RequestBody{}
		s(body)
		example := reflect.New(t.Elem()).Interface()
		ref, err := openapi3gen.NewSchemaRefForValue(example, openapi3.Schemas{})
		if err != nil {
			panic(err)
		}
		bodyParam := &openapi3.Parameter{
			Name:        "body",
			In:          "body",
			Description: body.descr,
			Schema:      ref,
		}
		if body.validationErr != "" {
			bodyParam.Required = true
			v.possibleResponse(body.validationErr, &Response[respData]{
				Status: body.validationStatus,
			})
		}
		v.operation.Parameters = append(v.operation.Parameters, &openapi3.ParameterRef{
			Value: bodyParam,
		})
		v.reqBody = body
		switch body.mediaType {
		case Json, "":
			body.mapper = func(bs []byte, ptr any) {
				err := json.Unmarshal(bs, ptr)
				if err != nil {
					panic(err)
				}
			}
		default:
			panic(fmt.Sprintf("mediaType '%v' for requestBody isn't supported", body.mediaType))
		}
		return
	}
	bs, err := io.ReadAll(v.r.Body)
	if err != nil {
		panic(err)
	}
	v.root.reqBody.mapper(bs, ptr)
}
