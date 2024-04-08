package ogo

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

type RequestBody struct {
	mediaType mediaType
	descr     string
	mapper    func(bs []byte, ptr any)
}

func (body *RequestBody) MediaType(mt mediaType) *RequestBody {
	body.mediaType = mt
	return body
}

func (body *RequestBody) Description(d string) *RequestBody {
	body.descr = d
	return body
}

type BodySettings func(body *RequestBody)

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
		content := openapi3.Content{}
		example := reflect.New(t.Elem()).Interface()
		content[string(body.mediaType)] = &openapi3.MediaType{
			Example: example,
		}
		v.operation.RequestBody = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: content,
			},
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
	bs, err := io.ReadAll(v.r.Body)
	if err != nil {
		panic(err)
	}
	v.root.reqBody.mapper(bs, ptr)
}
