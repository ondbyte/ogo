package models

import "github.com/ondbyte/ogo/petstore_example/db/pets"

type Category pets.Category
type Tag pets.Tag

type Pet struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	PhotoUrl string   `json:"photoUrls"`
	Status   string   `json:"status"`
	Category Category `json:"category"`
	Tags     []Tag    `json:"tags"`
}
