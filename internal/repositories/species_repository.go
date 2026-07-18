package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type SpeciesRepository struct {
	DB *gorm.DB
}

func NewSpeciesRepository(db *gorm.DB) *SpeciesRepository {
	return &SpeciesRepository{DB: db}
}

func (r *SpeciesRepository) All() ([]models.Species, error) {
	var species []models.Species
	err := r.DB.Model(&models.Species{}).Order("name asc").Find(&species).Error
	return species, err
}

func (r *SpeciesRepository) List(search string, page, perPage int) ([]models.Species, int64, error) {
	var species []models.Species
	tx := r.DB.Model(&models.Species{}).Preload("User").Order("id desc")
	if search != "" {
		tx = tx.Where("name LIKE ?", "%"+search+"%")
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&species).Error
	return species, total, err
}

func (r *SpeciesRepository) GetByID(id uint) (*models.Species, error) {
	var species models.Species
	err := r.DB.Preload("Breeds").First(&species, id).Error
	if err != nil {
		return nil, err
	}
	return &species, nil
}

func (r *SpeciesRepository) Create(species *models.Species) error {
	return r.DB.Create(species).Error
}

func (r *SpeciesRepository) Update(species *models.Species) error {
	return r.DB.Save(species).Error
}

func (r *SpeciesRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Species{}, id).Error
}
