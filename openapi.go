package ogo

import "github.com/getkin/kin-openapi/openapi3"

type ServerInfo struct {
	description string
	url         string
}

func (si *ServerInfo) Description(d string) {
	si.description = d
}

func (si *ServerInfo) Url(u string) {
	si.url = u
}

type SwaggerInfo struct {
	Extensions     map[string]interface{} `json:"-" yaml:"-"`
	title          string                 `json:"title" yaml:"title"` // Required
	description    string                 `json:"description,omitempty" yaml:"description,omitempty"`
	termsOfService string                 `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	contact        *Contact               `json:"contact,omitempty" yaml:"contact,omitempty"`
	license        *License               `json:"license,omitempty" yaml:"license,omitempty"`
	version        string                 `json:"version" yaml:"version"` // Required
}

func (i *SwaggerInfo) asOpenApi3Info() *openapi3.Info {
	var contact *openapi3.Contact
	var license *openapi3.License
	if i.contact != nil {
		contact = i.contact.toOpenApi3Contact()
	}
	if i.license != nil {
		license = i.license.toOpenApi3License()
	}
	return &openapi3.Info{
		Title:          i.title,
		Description:    i.description,
		TermsOfService: i.termsOfService,
		Version:        i.version,
		Contact:        contact,
		License:        license,
	}
}

type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Url   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type License struct {
	Name string `json:"name" yaml:"name"` // Required
	Url  string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (i *SwaggerInfo) Title(title string) {
	i.title = title
}

func (i *SwaggerInfo) Description(description string) {
	i.description = description
}

func (i *SwaggerInfo) TermsOfService(termsOfService string) {
	i.termsOfService = termsOfService
}

func (i *SwaggerInfo) Contact(contact *Contact) {
	i.contact = contact
}

func (i *SwaggerInfo) License(license *License) {
	i.license = license
}

func (i *SwaggerInfo) Version(version string) {
	i.version = version
}

func (c *Contact) toOpenApi3Contact() *openapi3.Contact {
	return &openapi3.Contact{
		Name:  c.Name,
		URL:   c.Url,
		Email: c.Email,
	}
}

func (k *License) toOpenApi3License() *openapi3.License {
	return &openapi3.License{
		Name: k.Name,
		URL:  k.Url,
	}
}

type SwaggerInfoSettings func(info *SwaggerInfo)
type ServerSettings func(info *ServerInfo)
