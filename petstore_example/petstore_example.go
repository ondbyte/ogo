package petstoreexample

import (
	_ "embed"

	"github.com/ondbyte/ogo"
	"github.com/ondbyte/ogo/petstore_example/handlers"
	"github.com/ondbyte/ogo/petstore_example/models"
	"github.com/ondbyte/ogo/petstore_example/service"
)

func demoPetService() *service.PetService {
	return service.NewPetService(
		map[uint]*models.Pet{},
		map[uint]*models.Category{
			0: &models.Category{
				Id:   0,
				Name: "dog",
			},
			2: &models.Category{
				Id:   1,
				Name: "dog",
			},
		},
		map[uint]*models.Tag{
			0: &models.Tag{
				Id:   0,
				Name: "brown",
			},
			2: &models.Tag{
				Id:   1,
				Name: "ginger",
			},
		},
	)
}

func Run() {
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

	petsService := demoPetService()
	handlers.CreatePet(o, petsService)
	handlers.GetPet(o, petsService)
	o.Run(":8080")
}
