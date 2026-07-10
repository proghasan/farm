package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type VaccineRepository struct {
	DB *gorm.DB
}

func NewVaccineRepository(db *gorm.DB) *VaccineRepository {
	return &VaccineRepository{DB: db}
}

func (r *VaccineRepository) List(speciesID string, page, perPage int) ([]models.Vaccine, int64, error) {
	var vaccines []models.Vaccine
	tx := r.DB.Model(&models.Vaccine{}).Preload("Species")
	if speciesID != "" {
		tx = tx.Where("species_id = ?", speciesID)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&vaccines).Error
	return vaccines, total, err
}

func (r *VaccineRepository) GetByID(id uint) (*models.Vaccine, error) {
	var vaccine models.Vaccine
	err := r.DB.Preload("Species").First(&vaccine, id).Error
	if err != nil {
		return nil, err
	}
	return &vaccine, nil
}

func (r *VaccineRepository) Create(vaccine *models.Vaccine) error {
	return r.DB.Create(vaccine).Error
}

func (r *VaccineRepository) Update(vaccine *models.Vaccine, updates map[string]interface{}) error {
	return r.DB.Model(vaccine).Updates(updates).Error
}

func (r *VaccineRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Vaccine{}, id).Error
}

func (r *VaccineRepository) Preload(vaccine *models.Vaccine) error {
	return r.DB.Preload("Species").First(vaccine, vaccine.ID).Error
}
