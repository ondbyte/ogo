package ogo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

type ValidationErr struct {
	status int
	err    string
}

type RequestBody struct {
	*ValidationErr
	mediaType   mediaType
	descr       string
	mapper      func(bs []byte, ptr any)
	requestBody *openapi3.RequestBody
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
	body.status = status
	body.err = err
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
		typ := reflect.TypeOf(ptr)
		if typ.Kind() != reflect.Pointer {
			panic("ptr should be a pointer")
		}
		val := reflect.ValueOf(ptr)
		if val.IsNil() {
			panic(`ptr cannot be empty, the pointer needs to be initialized to its zero value but its nil.
example :
ptr:=&YourBody{}
v.Body(ptr,....)
`)
		}
		body := &RequestBody{
			ValidationErr: &ValidationErr{},
		}
		s(body)
		example := val.Interface()
		ref, err := openapi3gen.NewSchemaRefForValue(example, openapi3.Schemas{})
		if err != nil {
			panic(err)
		}
		body.requestBody = &openapi3.RequestBody{
			Description: body.descr,
			Content: openapi3.Content{
				string(body.mediaType): &openapi3.MediaType{
					Schema: ref,
				},
			},
		}

		if body.ValidationErr != nil {
			body.requestBody.Required = true
			v.possibleResponse(body.err, v.validationErrHandler(body.status, body.err))
		}
		v.operation.RequestBody = &openapi3.RequestBodyRef{
			Value: body.requestBody,
		}
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
	if v.r.Body == http.NoBody && v.root.reqBody.requestBody.Required {
		v.write(v.root.validationErrHandler(v.root.reqBody.status, v.root.reqBody.err))
		return
	}
	bs, err := io.ReadAll(v.r.Body)
	if err != nil {
		panic(err)
	}
	v.root.reqBody.mapper(bs, ptr)
}
