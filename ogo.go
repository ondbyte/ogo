package ogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ondbyte/swagui"
)

type StatusCode int

type ErrorFormatter func(ReqErrs) (string, StatusCode)

type Ogo struct {
	info  *Info
	Ogo   bool
	Hmux  *http.ServeMux
	paths *openapi3.Paths
}

// ServeHTTP implements http.Handler.
func (m *Ogo) Run(addr string) error {
	m.serveSwaggerUi(addr)
	m.Ogo = false
	//fmt.Println("ogo running on: ", addr)
	return http.ListenAndServe(addr, m.Hmux)
}

func New(s OgoSettings) *Ogo {
	i := &Info{}
	s(i)
	return &Ogo{
		info:  i,
		Ogo:   true,
		paths: openapi3.NewPaths(),
		Hmux:  http.NewServeMux(),
	}
}

type mediaType string

var (
	Json mediaType = "application/json"
)

type RequestValidator[ValidatedData any, RespBody any] struct {
	root                  *RequestValidator[ValidatedData, RespBody]
	w                     http.ResponseWriter
	r                     *http.Request
	reqBody               *RequestBody //details about the incoming body
	method                string
	isRespBodyIsBasicType bool // is the response body is a basic type
	path                  string
	params                map[string]*Param
	query                 url.Values
	ogo                   bool
	operation             *openapi3.Operation
	validationErrHandler  validationErrHandler[RespBody]
}

type ValErr struct {
	Msg        string
	StatusCode int
}

type ReqErrs []*ValErr

func (rerrs ReqErrs) AddErr(r *ValErr) {
	if r != nil {
		rerrs = append(rerrs, r)
	}
}
func (v *RequestValidator[reqData, respData]) Description(d func() string) {
	if v.ogo {
		v.operation.Description = d()
	}
}

func (v *RequestValidator[reqData, respData]) Depricated(d func() bool) {
	if v.ogo {
		v.operation.Deprecated = d()
	}
}

/*
// return the possible response this api could respond
func (v *Validator[reqData, respData]) OnSuccess(d func() (status int, response *respData)) {
	if v.ogo {
		sc, rd := d()
		firstField := reflect.TypeOf(rd).Field(0)
		content := ""
		if _, isJson := firstField.Tag.Lookup("json"); isJson {
			content = "application/json"
		} else if _, isYaml := firstField.Tag.Lookup("yaml"); isYaml {
			content = "application/yaml"
		} else if _, isXml := firstField.Tag.Lookup("xml"); isXml {
			content = "application/xml"
		}
		v.operation.AddResponse(sc, &openapi3.Response{
			Content: openapi3.Content{
				content: &openapi3.MediaType{
					Example: rd,
				},
			},
		})
	}
} */

func (v *RequestValidator[reqData, respData]) HeaderParam(name string, ptr any, s ParamSettings) {
	v.param("header", name, ptr, s)
}

func (v *RequestValidator[reqData, respData]) CookieParam(name string, ptr any, s ParamSettings) {
	v.param("cookie", name, ptr, s)
}

func (v *RequestValidator[reqData, respData]) QueryParam(name string, ptr any, s ParamSettings) {
	v.param("query", name, ptr, s)
}

func (v *RequestValidator[reqData, respData]) PathParam(name string, ptr any, s ParamSettings) {
	v.param("path", name, ptr, s)
}

func (v *RequestValidator[reqData, respBody]) write(r *Response[respBody]) {
	defer func() {
		v.w = nil
	}()
	v.w.WriteHeader(r.Status)
	for _, h := range r.Headers {
		v.w.Header().Set(h.Key, h.Val)
	}
	if v.root.isRespBodyIsBasicType {
		v.w.Write([]byte(fmt.Sprintf("%v", r.Body)))
		return
	}
	switch r.MediaType {
	case Json, "":
		bs, err := json.Marshal(r.Body)
		if err != nil {
			panic(err)
		}
		v.w.Write(bs)
		return
	}
	panic(fmt.Sprint("response media type ", r.MediaType, " is not supported"))
}

func (v *RequestValidator[reqData, respData]) verifyPathParamIsInPath(
	name string,
) error {
	if strings.Contains(v.path, fmt.Sprintf("{%v}", name)) {
		return nil
	}
	return fmt.Errorf("path param '%v' is not in your path '%v'", name, v.path)
}

type Header struct {
	Key string
	Val string
}

func (v *RequestValidator[reqData, respBody]) possibleResponse(
	descr string,
	response *Response[respBody],
) {
	r := openapi3.NewResponse()
	r.Description = &descr
	if len(response.Headers) > 0 {
		r.Headers = openapi3.Headers{}
	}
	for _, h := range response.Headers {
		r.Headers[h.Key] = &openapi3.HeaderRef{
			Value: &openapi3.Header{
				Parameter: openapi3.Parameter{
					Name:   h.Key,
					Schema: getSchemaForPtr(&h.Val),
				},
			},
		}
	}
	var a interface{} = response.Body
	if a != nil {
		r.Content = openapi3.Content{
			string(response.MediaType): &openapi3.MediaType{
				Example: response.Body,
			},
		}
	}
	v.operation.AddResponse(response.Status, r)
}

func (v *RequestValidator[reqData, respData]) param(
	paramType string,
	name string,
	ptr any,
	settings ParamSettings,
) {
	if v.ogo {
		s := getSchemaForPtr(ptr)
		if s == nil {
			if reflect.TypeOf(ptr).Kind() != reflect.Pointer {
				panic(fmt.Sprintf("please pass a pointer as value for 'ptr' field, you have passed a non pointer to %v param '%v'", paramType, name))
			}
			panic(fmt.Sprintf("unsupported type '%T' for schema", ptr))
		}
		p := &Param{
			parameter: &openapi3.Parameter{
				Name: name, In: paramType, Schema: s,
			},
		}
		settings(p)
		v.operation.AddParameter(p.parameter)
		v.params[name] = p

		switch paramType {
		case "path":
			if err := v.verifyPathParamIsInPath(name); err != nil {
				panic(err)
			}
		case "query", "header", "cookie":
			if p.validationErr != "" {
				resp := v.validationErrHandler(p.validationStatus, p.validationErr)
				if resp == nil {
					panic(fmt.Sprintf("you have not handled validationStatus and validationErr in your validationErrHandler for path '%v'", v.path))
				}
				v.possibleResponse(fmt.Sprintf("if '%v' param '%v' is missing, api will respond with the following", paramType, name), resp)
			}
		}
		return
	}
	p := v.root.params[name]
	var err error
	var d string
	switch paramType {
	case "header":
		d = v.r.Header.Get(name)
		if d == "" && p.parameter.Required {
			v.write(v.root.validationErrHandler(p.validationStatus, p.validationErr))
			return
		}
	case "cookie":
		c, _ := v.r.Cookie(name)
		if c == nil && p.parameter.Required {
			v.write(v.root.validationErrHandler(p.validationStatus, p.validationErr))
			return
		}
		if c != nil {
			d = c.Value
		}
	case "query":
		d = v.query.Get(name)
		if d == "" && p.parameter.Required {
			v.write(v.root.validationErrHandler(p.validationStatus, p.validationErr))
			return
		}
	case "path":
		d = v.r.PathValue(name)
		if d == "" {
			v.write(v.root.validationErrHandler(p.validationStatus, p.validationErr))
			return
		}
	}
	switch ptr := ptr.(type) {
	case *string:
		*ptr = d
	case *int64:
		err = v.setInt64(d, ptr)
	case *float64:
		err = v.setFloat64(d, ptr)
	case *bool:
		err = v.setBool(d, ptr)
	default:
		err = json.Unmarshal([]byte(d), ptr)
	}
	if err != nil {
		panic(err)
	}
}

type RespContext[reqData any, respData any] struct {
	ReqData *reqData
}

func (ctx *RespContext[reqData, respData]) Write(response respData) {

}
func (m *Ogo) serveSwaggerUi(addr string) {
	t := &openapi3.T{
		Info:    m.info.asOpenApi3Info(),
		OpenAPI: "3.0.2",
	}
	t.Servers = append(t.Servers, &openapi3.Server{
		URL:         "http://localhost:8080",
		Description: "xyz",
	})
	t.Paths = m.paths
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	ep := "/swagger_doc/*"
	url := fmt.Sprintf("%v%v", addr, ep)
	m.Hmux.Handle(ep, swagui.Handle(b, swagui.Json))
	fmt.Println("swagger UI is at: ", url)
	fmt.Println("swagger Spec: ", string(b))
}

type Response[respBody any] struct {
	Status    int
	Headers   []*Header
	MediaType mediaType
	Body      *respBody
}
type validator[ValidatedData any, respBody any] func(v *RequestValidator[ValidatedData, respBody], reqData *ValidatedData)
type validationErrHandler[respBody any] func(validatedStatus int, validatedErr string) (resp *Response[respBody])
type handler[ValidatedData any, respBody any] func(reqData *ValidatedData) (resp *Response[respBody])

func SetupHandler[ValidatedData any, respBody any](
	mux *Ogo,
	method string,
	path string,
	validator validator[ValidatedData, respBody],
	validationErrHandler validationErrHandler[respBody],
	handler handler[ValidatedData, respBody],
) {
	if reflect.TypeOf(new(ValidatedData)).Elem().Kind() == reflect.Pointer {
		panic(fmt.Sprintf("SetupHandler[#1,#2] for path '%v', type #1 is pointer type, make it a non pointer type", path))
	}
	if reflect.TypeOf(new(respBody)).Elem().Kind() == reflect.Pointer {
		panic(fmt.Sprintf("SetupHandler[#1,#2] for path '%v', type #2 is pointer type, make it a non pointer type", path))
	}
	err, isBasic := isValidRespBodyType(reflect.TypeOf(new(respBody)).Elem())
	if err != nil {
		panic(err)
	}
	pi := mux.paths.Find(path)
	if pi == nil {
		pi = &openapi3.PathItem{}
		mux.paths.Set(path, pi)
	}
	op := pi.GetOperation(method)
	if op == nil {
		op = openapi3.NewOperation()
		pi.SetOperation(method, op)
	}
	root := &RequestValidator[ValidatedData, respBody]{
		method:                method,
		isRespBodyIsBasicType: isBasic,
		ogo:                   true,
		operation:             op,
		params:                map[string]*Param{},
		path:                  path,
		validationErrHandler:  validationErrHandler,
	}
	validator(root, new(ValidatedData))
	mux.Hmux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := &RequestValidator[ValidatedData, respBody]{
			ogo:   mux.Ogo,
			root:  root,
			query: r.URL.Query(),
			w:     w,
			r:     r,
		}
		data := new(ValidatedData)
		validator(v, data)
		if v.w == nil {
			// already written
			return
		}
		v.write(handler(data))
	})
}
