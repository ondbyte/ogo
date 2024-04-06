package ogo_test

import (
	"net/http"
	"testing"

	"github.com/ondbyte/ogo/ogo"
)

func TestOgo(t *testing.T) {
	o := ogo.New()
	type ValidatedData struct {
		UserId  string
		Cookie1 string
	}
	type ResponseBody struct {
		Msg string `json:"msg"`
	}
	ogo.SetupHandler[ValidatedData, ResponseBody](
		o,
		"GET",
		"/users/me",
		func(v *ogo.RequestValidator[ValidatedData, ResponseBody], reqData *ValidatedData) {
			v.Description(
				func() string {
					return "this endpoint returns the user details, please pass the 'userId' path param"
				})
			v.QueryParam(
				"userId", &reqData.UserId,
				func(param *ogo.Param) {
					param.Description("id of the user you are fetching")
					param.Required(http.StatusTeapot, "userId query param is required")
					param.Deprecated(true)
				},
			)
			v.CookieParam(
				"cookie1", &reqData.Cookie1,
				func(param *ogo.Param) {
					param.Description("cookie1 will be test cookie")
					param.Required(http.StatusTeapot, "cookie1 cookie param is required, please pass it in you cookie")
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
