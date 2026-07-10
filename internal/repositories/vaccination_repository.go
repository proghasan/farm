package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type VaccinationRepository struct {
	DB *gorm.DB
}

func NewVaccinationRepository(db *gorm.DB) *VaccinationRepository {
	return &VaccinationRepository{DB: db}
}

func (r *VaccinationRepository) List(animalID string, page, perPage int) ([]models.AnimalVaccination, int64, error) {
	var vaccinations []models.AnimalVaccination
	tx := r.DB.Model(&models.AnimalVaccination{}).Preload("Animal").Preload("Vaccine")
	if animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Order("vaccination_date DESC").Find(&vaccinations).Error
	return vaccinations, total, err
}

func (r *VaccinationRepository) GetByID(id uint) (*models.AnimalVaccination, error) {
	var v models.AnimalVaccination
	err := r.DB.Preload("Animal").Preload("Vaccine").First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VaccinationRepository) Create(v *models.AnimalVaccination) error {
	return r.DB.Create(v).Error
}

func (r *VaccinationRepository) Update(v *models.AnimalVaccination, updates map[string]interface{}) error {
	return r.DB.Model(v).Updates(updates).Error
}

func (r *VaccinationRepository) Delete(id uint) error {
	return r.DB.Delete(&models.AnimalVaccination{}, id).Error
}

func (r *VaccinationRepository) Preload(v *models.AnimalVaccination) error {
	return r.DB.Preload("Animal").Preload("Vaccine").First(v, v.ID).Error
}
