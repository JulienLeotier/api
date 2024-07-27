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
	// Save room
	if err := tx.Create(room).Error; err != nil {
		return nil, err
	}

	// Save files and associate them with the room
	for _, file := range room.Files {
		if err := tx.FirstOrCreate(&file, models.File{URL: file.URL}).Error; err != nil {
			return nil, err
		}
		if err := tx.Model(room).Association("Files").Append(&file); err != nil {
			return nil, err
		}
	}

	return room, nil
}

func (r *roomRepository) GetRoomByID(roomID string, tx *gorm.DB) (*models.Room, error) {
	var room models.Room
	if err := tx.Preload("Files").Preload("Variant").Preload("Detective").First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetAllRooms(tx *gorm.DB) ([]models.Room, error) {
	var rooms []models.Room
	if err := tx.Preload("Files").Preload("Variant").Preload("Detective").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) DeleteRoom(roomID string, tx *gorm.DB) error {
	// Get the room and its files
	room, err := r.GetRoomByID(roomID, tx)
	if err != nil {
		return err
	}

	// Delete the associated files
	for _, file := range room.Files {
		// Remove associations from the join table
		if err := tx.Model(&file).Association("Rooms").Clear(); err != nil {
			return err
		}

		// Remove the file from the file system
		if err := os.Remove(file.URL); err != nil && !os.IsNotExist(err) {
			return err
		}
		// Delete the file from the database
		if err := tx.Unscoped().Delete(&file).Error; err != nil {
			return err
		}
	}

	// Delete the room
	if err := tx.Unscoped().Delete(&room).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) UpdateRoom(room *models.Room, tx *gorm.DB) (*models.Room, error) {
	// Update room details
	if err := tx.Save(room).Error; err != nil {
		return nil, err
	}

	// Save new files and associate them with the room
	for _, file := range room.Files {
		if err := tx.FirstOrCreate(&file, models.File{URL: file.URL}).Error; err != nil {
			return nil, err
		}
		if err := tx.Model(room).Association("Files").Append(&file); err != nil {
			return nil, err
		}
	}

	return room, nil
}
