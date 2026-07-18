package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type WeightRepository struct {
	DB *gorm.DB
}

func NewWeightRepository(db *gorm.DB) *WeightRepository {
	return &WeightRepository{DB: db}
}

func (r *WeightRepository) List(animalID string) ([]models.AnimalWeightHistory, error) {
	var weights []models.AnimalWeightHistory
	tx := r.DB.Model(&models.AnimalWeightHistory{})
	if animalID != "" {
		tx = tx.Where("animal_id = ?", animalID)
	}
	err := tx.Order("id DESC").Find(&weights).Error
	return weights, err
}

func (r *WeightRepository) GetByID(id uint) (*models.AnimalWeightHistory, error) {
	var w models.AnimalWeightHistory
	err := r.DB.Preload("Animal").First(&w, id).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeightRepository) GetLatestByAnimalID(animalID uint) (*models.AnimalWeightHistory, error) {
	var w models.AnimalWeightHistory
	err := r.DB.Where("animal_id = ?", animalID).Order("id DESC").First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeightRepository) Create(w *models.AnimalWeightHistory) error {
	return r.DB.Create(w).Error
}

func (r *WeightRepository) Update(w *models.AnimalWeightHistory, updates map[string]interface{}) error {
	return r.DB.Model(w).Updates(updates).Error
}

func (r *WeightRepository) Delete(id uint) error {
	return r.DB.Delete(&models.AnimalWeightHistory{}, id).Error
}

func (r *WeightRepository) Preload(w *models.AnimalWeightHistory) error {
	return r.DB.Preload("Animal").First(w, w.ID).Error
}
