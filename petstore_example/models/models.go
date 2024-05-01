package models

type Category struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type Tag struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	Id        uint     `json:"id"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status    string   `json:"status"`
	Category  Category `json:"category"`
	Tags      []Tag    `json:"tags"`
}
