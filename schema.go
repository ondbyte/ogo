package ogo

import "github.com/getkin/kin-openapi/openapi3"

func getSchemaForPtr(v any) *openapi3.SchemaRef {
	switch v.(type) {
	case *string:
		return openapi3.NewStringSchema().NewRef()
	case *int, *int16, *int32, *int64, *uint, *uint16, *uint32, *uint64:
		return openapi3.NewInt64Schema().NewRef()
	case *float32, *float64:
		return openapi3.NewFloat64Schema().NewRef()
	}
	return nil
}
