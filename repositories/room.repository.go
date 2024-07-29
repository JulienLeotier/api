package repositories

import (
	"geniale/models"
	"os"

	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(room *models.Room, tx *gorm.DB) (*models.Room, error)
	GetRoomByID(roomID string, tx *gorm.DB) (*models.Room, error)
	GetAllRooms(tx *gorm.DB) ([]models.Room, error)
	DeleteRoom(roomID string, tx *gorm.DB) error
	UpdateRoom(room *models.Room, tx *gorm.DB) (*models.Room, error)
}

type roomRepository struct {
	DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{DB: db}
}

func (r *roomRepository) CreateRoom(room *models.Room, tx *gorm.DB) (*models.Room, error) {
	if err := tx.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (r *roomRepository) GetRoomByID(roomID string, tx *gorm.DB) (*models.Room, error) {
	var room models.Room
	if err := tx.Preload("EnigmaFiles").Preload("RevelationFiles").Preload("IntroFiles").Preload("OutroFiles").Preload("AmbianceMusicFiles").Preload("Variant").Preload("Detective").First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetAllRooms(tx *gorm.DB) ([]models.Room, error) {
	var rooms []models.Room
	if err := tx.Preload("EnigmaFiles").Preload("RevelationFiles").Preload("IntroFiles").Preload("OutroFiles").Preload("AmbianceMusicFiles").Preload("Variant").Preload("Detective").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) DeleteRoom(roomID string, tx *gorm.DB) error {
	room, err := r.GetRoomByID(roomID, tx)
	if err != nil {
		return err
	}

	files := append(room.EnigmaFiles, room.RevelationFiles...)
	files = append(files, room.IntroFiles...)
	files = append(files, room.OutroFiles...)
	files = append(files, room.AmbianceMusicFiles...)

	for _, file := range files {
		if err := tx.Model(&file).Association("Rooms").Clear(); err != nil {
			return err
		}
		if err := os.Remove(file.URL); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := tx.Unscoped().Delete(&file).Error; err != nil {
			return err
		}
	}

	if err := tx.Unscoped().Delete(&room).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) UpdateRoom(room *models.Room, tx *gorm.DB) (*models.Room, error) {
	if err := tx.Save(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}
