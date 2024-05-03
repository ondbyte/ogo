package handlers

import (
	"net/http"

	"github.com/ondbyte/ogo"
	"github.com/ondbyte/ogo/petstore_example/models"
	"github.com/ondbyte/ogo/petstore_example/service"
)

func GetPet(o *ogo.Server, petService *service.PetService) {
	ogo.SetupHandler[uint, models.Pet](
		o,
		"GET",
		"/pet/{petId}",
		func(v *ogo.RequestValidator[uint, models.Pet], petId *uint) {
			v.Summary(
				func() string {
					return "Find pet by ID"
				},
			)
			v.PathParam(
				"petId", petId,
				func(param *ogo.Param) {
					param.Description("ID of the pet to return")
					param.IfInvalidType(http.StatusBadRequest, "Invalid ID supplied")
				},
			)
			v.PossibleResponse(func() (descr string, response *ogo.Response[models.Pet]) {
				return "Pet not found", &ogo.Response[models.Pet]{
					Status: http.StatusNotFound,
				}
			})
		},
		func(validatedStatus int, validatedErr string) (resp *ogo.Response[models.Pet]) {
			if validatedErr != "" {
				return &ogo.Response[models.Pet]{
					Status:  validatedStatus,
					RawBody: []byte("Invalid ID supplied"),
				}
			}
			return nil
		},
		func(petId *uint) (resp *ogo.Response[models.Pet]) {
			return nil
		},
	)
}
