package ogo_test

import (
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

	o.Run(":8080")
}
