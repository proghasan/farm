package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type PregnancyRepository struct {
	DB *gorm.DB
}

func NewPregnancyRepository(db *gorm.DB) *PregnancyRepository {
	return &PregnancyRepository{DB: db}
}

func (r *PregnancyRepository) List(animalID, status string, page, perPage int) ([]models.AnimalPregnancy, int64, error) {
	var items []models.AnimalPregnancy
	tx := r.DB.Model(&models.AnimalPregnancy{}).Preload("Animal").Preload("Breeder")
	if animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&items).Error
	return items, total, err
}

func (r *PregnancyRepository) GetByID(id uint) (*models.AnimalPregnancy, error) {
	var item models.AnimalPregnancy
	err := r.DB.Preload("Animal").Preload("Breeder").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PregnancyRepository) Create(item *models.AnimalPregnancy) error {
	return r.DB.Create(item).Error
}

func (r *PregnancyRepository) Update(item *models.AnimalPregnancy, updates map[string]interface{}) error {
	return r.DB.Model(item).Updates(updates).Error
}

func (r *PregnancyRepository) Delete(id uint) error {
	return r.DB.Delete(&models.AnimalPregnancy{}, id).Error
}

func (r *PregnancyRepository) Preload(item *models.AnimalPregnancy) error {
	return r.DB.Preload("Animal").Preload("Breeder").First(item, item.ID).Error
}

func (r *PregnancyRepository) ListByAnimalID(animalID uint) ([]models.AnimalPregnancy, error) {
	var items []models.AnimalPregnancy
	err := r.DB.Where("animal_id = ?", animalID).Order("created_at DESC").Find(&items).Error
	return items, err
}
