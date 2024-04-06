package ogo

import "github.com/getkin/kin-openapi/openapi3"

type Info struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`

	Title          string   `json:"title" yaml:"title"` // Required
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"` // Required
}

func (i *Info) asOpenApi3Info() *openapi3.Info {
	return &openapi3.Info{
		Title:          i.Title,
		Description:    i.Description,
		TermsOfService: i.TermsOfService,
		Version:        i.Version,
		Contact:        i.Contact.toOpenApi3Contact(),
		License:        i.License.toOpenApi3License(),
	}
}

type Contact struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`

	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

func (c *Contact) toOpenApi3Contact() *openapi3.Contact {
	return &openapi3.Contact{
		Name:  c.Name,
		URL:   c.URL,
		Email: c.Email,
	}
}

type License struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`

	Name string `json:"name" yaml:"name"` // Required
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (k *License) toOpenApi3License() *openapi3.License {
	return &openapi3.License{
		Name: k.Name,
		URL:  k.URL,
	}
}
