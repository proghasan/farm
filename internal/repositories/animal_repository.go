package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type AnimalRepository struct {
	DB *gorm.DB
}

func NewAnimalRepository(db *gorm.DB) *AnimalRepository {
	return &AnimalRepository{DB: db}
}

func (r *AnimalRepository) List(search string, filters map[string]string, page, perPage int) ([]models.Animal, int64, error) {
	var animals []models.Animal
	tx := r.DB.Model(&models.Animal{}).Preload("Breed.Species").Preload("User").Order("id desc")
	if search != "" {
		tx = tx.Where("tag_no LIKE ? OR color LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if v, ok := filters["species_id"]; ok && v != "" {
		tx = tx.Where("species_id = ?", v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		tx = tx.Where("status = ?", v)
	}
	if v, ok := filters["gender"]; ok && v != "" {
		tx = tx.Where("gender = ?", v)
	}
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&animals).Error
	return animals, total, err
}

func (r *AnimalRepository) GetByID(id uint) (*models.Animal, error) {
	var animal models.Animal
	err := r.DB.
		Preload("Breed.Species").
		Preload("User").
		Preload("WeightHistories").
		Preload("AnimalVaccinations.Vaccine").
		First(&animal, id).Error
	if err != nil {
		return nil, err
	}
	return &animal, nil
}

func (r *AnimalRepository) GetProfile(id uint) (*models.Animal, error) {
	var animal models.Animal
	err := r.DB.
		Preload("Breed.Species").
		Preload("User").
		Preload("Father").
		Preload("Mother").
		Preload("WeightHistories").
		Preload("AnimalVaccinations.Vaccine").
		First(&animal, id).Error
	if err != nil {
		return nil, err
	}
	return &animal, nil
}

func (r *AnimalRepository) Create(animal *models.Animal) error {
	return r.DB.Create(animal).Error
}

func (r *AnimalRepository) Update(animal *models.Animal) error {
	return r.DB.Save(animal).Error
}

func (r *AnimalRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Animal{}, id).Error
}

func (r *AnimalRepository) Preload(animal *models.Animal) error {
	return r.DB.Preload("Breed.Species").Preload("User").First(animal, animal.ID).Error
}

func (r *AnimalRepository) UpdateCurrentWeight(animalID uint, weight float64) error {
	return r.DB.Model(&models.Animal{}).Where("id = ?", animalID).Update("current_weight", weight).Error
}
