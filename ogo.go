package ogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/ondbyte/swagui"
)

type StatusCode int

type ErrorFormatter func(ReqErrs) (string, StatusCode)

type Server struct {
	info       *SwaggerInfo
	ogo        bool
	hmux       *http.ServeMux
	paths      *openapi3.Paths
	serverInfo *ServerInfo
}

// ServeHTTP implements http.Handler.
func (m *Server) Run(port uint, serverSettings ServerSettings) error {
	i := &ServerInfo{}
	if serverSettings != nil {
		serverSettings(i)
	}
	url := fmt.Sprintf("http://localhost:%v", port)
	if i.url == "" {
		i.url = url
	}
	m.serverInfo = i
	m.serveSwaggerUi(url)
	m.ogo = false
	//fmt.Println("ogo running on: ", addr)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), m.hmux)
}

func NewServer(s SwaggerInfoSettings) *Server {
	i := &SwaggerInfo{}
	if s != nil {
		s(i)
	}
	return &Server{
		info:  i,
		ogo:   true,
		paths: openapi3.NewPaths(),
		hmux:  http.NewServeMux(),
	}
}

type mediaType string

var (
	Json mediaType = "application/json"
)

type RequestValidator[ValidatedData any, RespBody any] struct {
	root    *RequestValidator[ValidatedData, RespBody]
	w       http.ResponseWriter
	r       *http.Request
	reqBody *RequestBody //details about the incoming body
	// err encountered while validation
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

// set the description of the endpoint
func (v *RequestValidator[reqData, respData]) Description(d func() string) {
	if v.ogo {
		v.operation.Description = d()
	}
}

// sets the summary of the the endpoint
func (v *RequestValidator[reqData, respData]) Summary(d func() string) {
	if v.ogo {
		v.operation.Summary = d()
	}
}

// mark this endpoint as deprecated
func (v *RequestValidator[reqData, respData]) Depricated(d func() bool) {
	if v.ogo {
		v.operation.Deprecated = d()
	}
}

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

// this can be used to set responses that are unknown while validating,
// for examples the responses that arise from your internal handling of the request data
func (v *RequestValidator[reqData, respData]) PossibleResponse(fn func() (descr string, response *Response[respData])) {
	if v.ogo {
		v.possibleResponse(fn())
	}
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
		var bs []byte
		if r.Body != nil {
			data, err := json.Marshal(r.Body)
			if err != nil {
				panic(err)
			}
			bs = data
		} else if r.RawBody != nil {
			bs = r.RawBody
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
	ref, err := openapi3gen.NewSchemaRefForValue(response.Body, openapi3.Schemas{})
	if err != nil {
		panic(err)
	}
	if response.MediaType != "" && response.Body != nil {
		media := &openapi3.MediaType{
			Schema: ref,
		}
		if response.Body != nil {
			media.Example = response.Body
		}
		r.Content = openapi3.Content{
			string(response.MediaType): media,
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
			//path param is required by default
			p.parameter.Required = true
			if err := v.verifyPathParamIsInPath(name); err != nil {
				panic(err)
			}
		case "query", "header", "cookie":
			if p.requiredErr != "" {
				resp := v.validationErrHandler(p.requiredStatus, p.requiredErr)
				if resp == nil {
					panic(fmt.Sprintf("you have not handled requiredStatus and requiredErr in your validationErrHandler for path '%v'", v.path))
				}
				v.possibleResponse(
					/* fmt.Sprintf("if '%v' param '%v' is missing, api will respond with the following", paramType, name) */
					p.requiredErr, resp,
				)
			}
		}

		if p.invalidTypeErr != "" {
			resp := v.validationErrHandler(p.invalidTypeStatus, p.invalidTypeErr)
			if resp == nil {
				panic(fmt.Sprintf("you have not handled invalidTypeStatus and invalidTypeErr in your validationErrHandler for path '%v'", v.path))
			}
			v.possibleResponse(p.invalidTypeErr, resp)
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
			v.write(v.root.validationErrHandler(p.requiredStatus, p.requiredErr))
			return
		}
	case "cookie":
		c, _ := v.r.Cookie(name)
		if c == nil && p.parameter.Required {
			v.write(v.root.validationErrHandler(p.requiredStatus, p.requiredErr))
			return
		}
		if c != nil {
			d = c.Value
		}
	case "query":
		d = v.query.Get(name)
		if d == "" && p.parameter.Required {
			v.write(v.root.validationErrHandler(p.requiredStatus, p.requiredErr))
			return
		}
	case "path":
		d = v.r.PathValue(name)
		if d == "" {
			v.write(v.root.validationErrHandler(p.requiredStatus, p.requiredErr))
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
func (m *Server) serveSwaggerUi(url string) {
	t := &openapi3.T{
		Info:    m.info.asOpenApi3Info(),
		OpenAPI: "3.0.2",
	}
	t.Servers = append(t.Servers, &openapi3.Server{
		URL:         m.serverInfo.url,
		Description: m.serverInfo.description,
	})
	t.Paths = m.paths
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	ep := "/swagger_doc/*"
	url = fmt.Sprintf("%v%v", url, ep)
	m.hmux.Handle(ep, swagui.Handle(b, swagui.Json))
	fmt.Println("swagger UI is at: ", url)
	fmt.Println("swagger Spec: ", string(b))
}

type Response[respBody any] struct {
	Status    int
	Headers   []*Header
	MediaType mediaType
	Body      *respBody

	// also useful when response are just string not any structured data like json,yml etc,
	// this will be considered only if the Body is nil
	RawBody []byte
}

type validator[ValidatedData any, respBody any] func(v *RequestValidator[ValidatedData, respBody], reqData *ValidatedData)
type validationErrHandler[respBody any] func(validatedStatus int, validatedErr string) (resp *Response[respBody])
type handler[ValidatedData any, respBody any] func(reqData *ValidatedData) (resp *Response[respBody])

func SetupHandler[ValidatedData any, respBody any | string](
	mux *Server,
	method string,
	path string,
	validator validator[ValidatedData, respBody],
	validationErrHandler validationErrHandler[respBody],
	handler handler[ValidatedData, respBody],
) {
	// these two check make sure generic types are non pointer
	if reflect.TypeOf(new(ValidatedData)).Elem().Kind() == reflect.Pointer {
		panic(fmt.Sprintf("SetupHandler[#1,#2] for path '%v', type #1 is pointer type, make it a non pointer type", path))
	}
	if reflect.TypeOf(new(respBody)).Elem().Kind() == reflect.Pointer {
		panic(fmt.Sprintf("SetupHandler[#1,#2] for path '%v', type #2 is pointer type, make it a non pointer type", path))
	}

	// this makes sure response is of basic type(ex: a string) or a struct
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
		op.Responses = openapi3.NewResponses()
		op.Responses.Delete("default")
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
	root.possibleResponse("default success response", &Response[respBody]{
		Status:    http.StatusOK,
		Body:      new(respBody),
		MediaType: Json,
	})
	validator(root, new(ValidatedData))
	mux.hmux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := &RequestValidator[ValidatedData, respBody]{
			ogo:   mux.ogo,
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
