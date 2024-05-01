package service

import (
	"context"
	"fmt"

	"github.com/ondbyte/ogo/petstore_example/models"
)

type PetService struct {
	pets       map[uint]*models.Pet
	categories map[uint]*models.Category
	tags       map[uint]*models.Tag
}

func NewPetService(
	pets map[uint]*models.Pet,
	category map[uint]*models.Category,
	tags map[uint]*models.Tag,
) *PetService {
	return &PetService{
		pets:       pets,
		categories: category,
		tags:       tags,
	}
}

func (s *PetService) AddPet(ctx context.Context, ap *models.Pet) (*models.Pet, error) {
	for _, p := range s.pets {
		if p.Id == ap.Id || p.Name == ap.Name {
			return nil, fmt.Errorf("pet exists")
		}
	}
	var validCategory bool
	for _, c := range s.categories {
		if c.Id == ap.Category.Id && c.Name == ap.Category.Name {
			validCategory = true
		}
	}
	if !validCategory {
		return nil, fmt.Errorf("invalid category while add pet")
	}
	if !validCategory {
		return nil, fmt.Errorf("invalid category while add pet")
	}
	var invalidTags string
	for _, t1 := range s.tags {
		for _, t2 := range ap.Tags {
			if t1.Id == t2.Id && t1.Name == t2.Name {
				continue
			}
			invalidTags += fmt.Sprintf("\ninvalid tag %v while add pet", t2)
			break
		}
	}
	if invalidTags != "" {
		return nil, fmt.Errorf(invalidTags)
	}

	s.pets[ap.Id] = ap
	return ap, nil
}
