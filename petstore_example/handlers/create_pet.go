package handlers

import (
	"context"
	"net/http"

	"github.com/ondbyte/ogo"
	"github.com/ondbyte/ogo/petstore_example/models"
	"github.com/ondbyte/ogo/petstore_example/service"
)

func CreatePet(o *ogo.Server, petService *service.PetService) {
	ogo.SetupHandler[models.Pet, models.Pet](
		o,
		"POST",
		"/pet",
		func(v *ogo.RequestValidator[models.Pet, models.Pet], reqData *models.Pet) {
			v.Body(reqData, func(body *ogo.RequestBody) {
				body.Required(405, "Invalid input")
				body.Description("Pet object that needs to be added to the store")
				body.MediaType(ogo.Json)
			})
			v.Summary(
				func() string {
					return "Add a new pet to the store"
				},
			)
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
			createdPet, err := petService.AddPet(context.TODO(), reqData)
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

}
