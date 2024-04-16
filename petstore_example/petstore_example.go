package petstoreexample

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/ondbyte/ogo"
	"github.com/ondbyte/ogo/petstore_example/db"
	"github.com/ondbyte/ogo/petstore_example/models"
	"github.com/ondbyte/ogo/petstore_example/service"
)

// go:embed schema.sql
var dbSchema string

func Run() {
	db, err := db.InitDb(context.TODO(), dbSchema)
	if err != nil {
		panic(err)
	}
	petsService := service.NewPetService(db)
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
	if err != nil {
		panic(err)
	}
	ogo.SetupHandler[models.Pet, models.Pet](
		o,
		"POST",
		"/pet",
		func(v *ogo.RequestValidator[models.Pet, models.Pet], reqData *models.Pet) {
			v.Body(reqData, func(body *ogo.RequestBody) {
				body.Description("create a new pet in the store")
				body.MediaType(ogo.Json)
			})
		},
		func(validatedStatus int, validatedErr string) (resp *ogo.Response[models.Pet]) {
			if validatedErr != "" {
				return &ogo.Response[models.Pet]{
					Status: validatedStatus,
				}
			}
			return nil
		},
		func(reqData *models.Pet) (resp *ogo.Response[models.Pet]) {
			createdPet, err := petsService.AddPet(context.TODO(), reqData)
			if err != nil {
				return &ogo.Response[models.Pet]{
					Status: http.StatusInternalServerError,
				}
			}
			return &ogo.Response[models.Pet]{
				Status:    http.StatusOK,
				MediaType: ogo.Json,
				Body:      createdPet,
			}
		},
	)

	o.Run(":8080")
}
