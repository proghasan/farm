package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type BreedRepository struct {
	DB *gorm.DB
}

func NewBreedRepository(db *gorm.DB) *BreedRepository {
	return &BreedRepository{DB: db}
}

func (r *BreedRepository) List(speciesID string, page, perPage int) ([]models.Breed, int64, error) {
	var breeds []models.Breed
	tx := r.DB.Model(&models.Breed{}).Preload("Species")
	if speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&breeds).Error
	return breeds, total, err
}

func (r *BreedRepository) GetByID(id uint) (*models.Breed, error) {
	var breed models.Breed
	err := r.DB.Preload("Species").First(&breed, id).Error
	if err != nil {
		return nil, err
	}
	return &breed, nil
}

func (r *BreedRepository) Create(breed *models.Breed) error {
	return r.DB.Create(breed).Error
}

func (r *BreedRepository) Update(breed *models.Breed, updates map[string]interface{}) error {
	return r.DB.Model(breed).Updates(updates).Error
}

func (r *BreedRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Breed{}, id).Error
}

func (r *BreedRepository) Preload(breed *models.Breed) error {
	return r.DB.Preload("Species").First(breed, breed.ID).Error
}
