package service

import (
	"context"
	"database/sql"

	"github.com/ondbyte/ogo/petstore_example/db"
	"github.com/ondbyte/ogo/petstore_example/db/pets"
	"github.com/ondbyte/ogo/petstore_example/models"
)

type Service struct {
	db *db.Db
}

func NewPetService(db *db.Db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) AddPet(ctx context.Context, ap *models.Pet) (*models.Pet, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	petx := s.db.Pets.WithTx(tx)
	defer tx.Rollback()
	p, err := petx.AddPet(ctx, pets.AddPetParams{
		Name: ap.Name,
		PhotoUrl: sql.NullString{
			String: ap.PhotoUrl,
		},
		Status: ap.Status,
		Category: sql.NullInt32{
			Int32: int32(ap.Category.ID),
		},
	})
	if err != nil {
		return nil, err
	}
	id, err := p.LastInsertId()
	if err != nil {
		return nil, err
	}
	for _, v := range ap.Tags {
		result, err := petx.AddPetTag(ctx, pets.AddPetTagParams{
			PetID: int32(id),
			TagID: v.ID,
		})
		if err != nil {
			return nil, err
		}
		_, err = result.LastInsertId()
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	ap.Id = id
	return ap, err
}
func (s *Service) AddCategory(ctx context.Context, ap *models.Category) (*models.Category, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	petx := s.db.Pets.WithTx(tx)
	defer tx.Rollback()
	result, err := petx.AddCategory(ctx, ap.Name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	ap.ID = id
	return ap, nil
}
