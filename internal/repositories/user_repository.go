package repositories

import (
	"farm/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) List(page, perPage int) ([]models.User, int64, error) {
	var users []models.User
	tx := r.DB.Model(&models.User{}).Select("id, name, email, phone, username, role, status, created_at, updated_at")
	var total int64
	tx.Count(&total)
	offset := (page - 1) * perPage
	err := tx.Offset(offset).Limit(perPage).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Select("id, name, email, phone, username, role, status, created_at, updated_at").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByIDFull(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *models.User, updates map[string]interface{}) error {
	return r.DB.Model(user).Updates(updates).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindByLogin(login string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ? OR phone = ? OR username = ?", login, login, login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
