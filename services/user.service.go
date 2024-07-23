package services

import (
	"errors"
	"strconv"

	"geniale/models"
	"geniale/repositories"
	"geniale/utils"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo      *repositories.UserRepository
	GroupRepo     *repositories.GroupRepository
	GroupUserRepo *repositories.GroupUserRepository
}

func NewUserService(userRepo *repositories.UserRepository, groupRepo *repositories.GroupRepository, groupUserRepo *repositories.GroupUserRepository) *UserService {
	return &UserService{
		UserRepo:      userRepo,
		GroupRepo:     groupRepo,
		GroupUserRepo: groupUserRepo,
	}
}

func (s *UserService) ExistingUser(user models.User) error {
	existingUser, err := s.UserRepo.FindByUsernameAndEmail(user.Username, user.Email)
	if err != nil {
		return err
	}

	existingUserCount := len(existingUser)
	if existingUserCount > 0 {
		if existingUserCount == 1 && existingUser[0].ID == user.ID {
			return nil
		}
		return errors.New("user already exists")
	}
	return nil
}
func (s *UserService) CreateTmpUser(validatedData models.UserCreateTmpDTO, tx *gorm.DB) (*models.User, error) {
	newUser := &models.User{
		Email:    validatedData.Email,
		Password: "",
		Username: validatedData.Username,
	}

	if err := s.ExistingUser(*newUser); err != nil {
		return nil, err
	}

	if err := tx.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil

}
func (s *UserService) CreateUser(validatedData models.UserCreateDTO, tx *gorm.DB) (*models.User, *models.GroupUser, error) {
	hashedPassword, _ := utils.HashPassword(validatedData.Password)
	newUser := &models.User{
		Email:    validatedData.Email,
		Password: hashedPassword,
		Username: validatedData.Username,
	}

	if err := s.ExistingUser(*newUser); err != nil {
		return nil, nil, err
	}

	if err := tx.Create(newUser).Error; err != nil {
		return nil, nil, err
	}

	group, err := s.GroupRepo.FindByName("user")
	if err != nil {
		return nil, nil, errors.New("group not found")
	}

	groupUser := &models.GroupUser{
		GroupID: group.ID,
		UserID:  newUser.ID,
	}

	if err := tx.Create(groupUser).Error; err != nil {
		return nil, nil, err
	}

	return newUser, groupUser, nil
}

func (s *UserService) LoginUser(validatedData models.LoginRequestDTO) (string, error) {
	foundUser, err := s.UserRepo.FindByEmail(validatedData.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(foundUser.Password, validatedData.Password) {
		return "", errors.New("incorrect password")
	}

	groupName, err := s.GroupRepo.FindGroupNameByUserID(foundUser.ID)
	if err != nil {
		return "", errors.New("could not get user group")
	}
	idToString := strconv.Itoa(int(foundUser.ID))
	token, err := utils.GenerateToken(foundUser.Email, groupName, idToString)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.UserRepo.FindByID(strconv.Itoa(int(id)))
}

func (s *UserService) FindByEmail(email string) (*models.User, error) {
	return s.UserRepo.FindByEmail(email)
}

func (s *UserService) FindByPasswordResetToken(token string, user *models.User) error {
	return s.UserRepo.FindByPasswordResetToken(token, user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.UserRepo.FindByID(id)
}

func (s *UserService) UpdateUser(user *models.User, tx *gorm.DB) error {
	return s.UserRepo.Update(user, tx)
}

func (s *UserService) SaveUserCode(user *models.User, code string, tx *gorm.DB) error {
	userCode := &models.UserCode{
		UserID: user.ID,
		Code:   code,
	}

	if err := tx.Create(userCode).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) CheckCode(user *models.User, code string, tx *gorm.DB) error {
	userCode := &models.UserCode{
		UserID: user.ID,
		Code:   code,
	}
	if err := tx.Where("user_id = ? AND code = ?", user.ID, code).First(userCode).Error; err != nil {
		return err
	}
	if err := tx.Delete(userCode).Error; err != nil {
		return err
	}
	if err := tx.Model(user).Update("email_checked", true).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetCode(user *models.User) (*models.UserCode, error) {
	result, err := s.UserRepo.GetCode(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.UserCode{}, nil
		}
		return nil, err
	}
	return result, nil
}
