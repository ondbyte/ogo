package ogo

import "github.com/getkin/kin-openapi/openapi3"

type ParamOpt func(qs *Param)
type ParamSettings func(param *Param)

type PramOpts func()

type Param struct {
	parameter        *openapi3.Parameter
	validationStatus int
	validationErr    string
}

// makes this param as required, return the status code and error you need to respond when its missing
func (p *Param) Required(statusCode int, err string) *Param {
	p.validationStatus = statusCode
	p.validationErr = err
	p.parameter.Required = true
	return p
}

func (p *Param) Description(d string) *Param {
	p.parameter.Description = d
	return p
}

func (p *Param) Deprecated(d bool) *Param {
	p.parameter.Deprecated = d
	return p
}

func (p *Param) Example(e interface{}) *Param {
	p.parameter.Example = e
	return p
}
