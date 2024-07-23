package repositories

import (
	"geniale/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByPasswordResetToken(token string, user *models.User) error {
	return r.DB.Where("password_reset_token = ?", token).First(user).Error
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User, tx *gorm.DB) error {
	return tx.Save(user).Error
}

func (r *UserRepository) FindByUsernameAndEmail(username, email string) ([]*models.User, error) {
	var users []*models.User
	result := r.DB.Where("username = ? AND email = ?", username, email).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) GetCode(user *models.User) (*models.UserCode, error) {
	userCode := &models.UserCode{}
	result := r.DB.Where("user_id = ?", user.ID).First(userCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &models.UserCode{}, nil
		}
		return nil, result.Error
	}
	return userCode, nil
}
