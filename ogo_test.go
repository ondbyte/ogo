package ogo_test

import (
	"net/http"
	"testing"

	"github.com/ondbyte/ogo"
)

func TestOgo(t *testing.T) {
	o := ogo.New(
		func(info *ogo.Info) {
			info.Title("Swagger Petstore")
			info.Description("This is a sample server Petstore server. You can find out more about Swagger at http://swagger.io or on irc.freenode.net, #swagger. For this sample, you can use the api key special-key to test the authorization filters.")
			info.TermsOfService("https://smartbear.com/terms-of-use/")
			info.Contact(
				&ogo.Contact{
					Name:  "Yadunandan",
					Url:   "http://yadunandan.xyz",
					Email: "iamyadunandan@gmail.com",
				},
			)
			info.License(
				&ogo.License{
					Name: "Apache 2.0",
					Url:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			)
			info.Version("0.0.1")
		},
	)
	type Body struct {
		Name string `json:"name"`
	}
	type ValidatedData struct {
		PathId   string
		HeaderId string
		QueryId  string
		CookieId string
		Body     Body
	}
	type ResponseBody struct {
		Msg string `json:"msg"`
	}
	ogo.SetupHandler[ValidatedData, ResponseBody](
		o,
		"POST",
		"/users/{pathId}",
		func(v *ogo.RequestValidator[ValidatedData, ResponseBody], reqData *ValidatedData) {
			v.Description(
				func() string {
					return "this endpoint returns the user details, please pass the 'userId' path param"
				})
			v.PathParam(
				"pathId", &reqData.PathId,
				func(param *ogo.Param) {
					param.Description("id of the user you are fetching")
					param.Required(http.StatusTeapot, "userId path param is required")
					param.Deprecated(true)
				},
			)
			v.HeaderParam(
				"headerId", &reqData.HeaderId,
				func(param *ogo.Param) {
					param.Description("id of the user you are fetching")
					param.Required(http.StatusTeapot, "headerId path query is required")
					param.Deprecated(true)
				},
			)
			v.QueryParam(
				"queryId", &reqData.QueryId,
				func(param *ogo.Param) {
					param.Description("id of the user you are fetching")
					param.Required(http.StatusTeapot, "queryId path query is required")
					param.Deprecated(true)
				},
			)

			v.Body(
				&reqData.Body,
				func(body *ogo.RequestBody) {
					body.MediaType(ogo.Json)
					body.Description("this is body's description, you can pass it as a json as we have set media type as json")
				},
			)
		},
		func(validatedStatus int, validatedErr string) (resp *ogo.Response[ResponseBody]) {
			if validatedErr != "" {
				return &ogo.Response[ResponseBody]{
					Status:    validatedStatus,
					MediaType: ogo.Json,
					Body: &ResponseBody{
						Msg: validatedErr,
					},
				}
			}
			return nil
		},
		func(reqData *ValidatedData) (resp *ogo.Response[ResponseBody]) {
			return &ogo.Response[ResponseBody]{
				Status:    http.StatusOK,
				MediaType: ogo.Json,
				Headers: []*ogo.Header{
					{
						Key: "key",
						Val: "value",
					},
				},
				Body: &ResponseBody{
					Msg: "success",
				},
			}
		},
	)

	o.Run(":8080")
}
