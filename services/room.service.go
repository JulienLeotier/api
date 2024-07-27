package services

import (
	"geniale/models"
	"geniale/repositories"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type RoomService interface {
	CreateRoom(validationData models.RoomCreateDTO, tx *gorm.DB) (*models.RoomResponseDTO, error)
	GetRoomByID(roomID string, tx *gorm.DB) (*models.RoomResponseDTO, error)
	GetAllRooms(tx *gorm.DB) ([]models.RoomResponseDTO, error)
	DeleteRoom(roomID string, tx *gorm.DB) error
	UpdateRoom(roomID string, validationData models.RoomCreateDTO, tx *gorm.DB) (*models.RoomResponseDTO, error)
}

type roomService struct {
	repository repositories.RoomRepository
	userRepo   repositories.UserRepository
}

func NewRoomService(repo repositories.RoomRepository, userRepo repositories.UserRepository) RoomService {
	return &roomService{repository: repo, userRepo: userRepo}
}

func (s *roomService) CreateRoom(validationData models.RoomCreateDTO, tx *gorm.DB) (*models.RoomResponseDTO, error) {
	var files []models.File
	for _, fileHeader := range validationData.Files {
		// Save the file to the server
		filename := filepath.Base(fileHeader.Filename)
		filePath := "./uploads/" + filename
		if err := saveUploadedFile(fileHeader, filePath); err != nil {
			return nil, err
		}

		// Create File model
		file := models.File{URL: filePath}
		files = append(files, file)
	}

	room := &models.Room{
		Name:        validationData.Name,
		Description: validationData.Description,
		Files:       files,
		VariantID:   validationData.Variant,
		DetectiveID: validationData.Detective,
	}
	roomSave, err := s.repository.CreateRoom(room, tx)
	if err != nil {
		return nil, err
	}

	// Prepare response DTO
	responseFiles := make([]string, len(files))
	for i, file := range files {
		responseFiles[i] = file.URL
	}
	return &models.RoomResponseDTO{
		ID:          roomSave.ID,
		Name:        roomSave.Name,
		Description: roomSave.Description,
		Files:       responseFiles,
		Variant:     roomSave.Variant,
		Detective:   roomSave.Detective,
	}, nil
}

func (s *roomService) GetRoomByID(roomID string, tx *gorm.DB) (*models.RoomResponseDTO, error) {
	room, err := s.repository.GetRoomByID(roomID, tx)
	if err != nil {
		return nil, err
	}

	// Prepare response DTO
	responseFiles := make([]string, len(room.Files))
	for i, file := range room.Files {
		responseFiles[i] = file.URL
	}
	return &models.RoomResponseDTO{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Files:       responseFiles,
		Variant:     room.Variant,
		Detective:   room.Detective,
	}, nil
}

func (s *roomService) GetAllRooms(tx *gorm.DB) ([]models.RoomResponseDTO, error) {
	rooms, err := s.repository.GetAllRooms(tx)
	if err != nil {
		return nil, err
	}

	var responseRooms []models.RoomResponseDTO
	for _, room := range rooms {
		responseFiles := make([]string, len(room.Files))
		for i, file := range room.Files {
			responseFiles[i] = file.URL
		}
		responseRooms = append(responseRooms, models.RoomResponseDTO{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Files:       responseFiles,
			Variant:     room.Variant,
			Detective:   room.Detective,
		})
	}

	return responseRooms, nil
}

func (s *roomService) DeleteRoom(roomID string, tx *gorm.DB) error {
	return s.repository.DeleteRoom(roomID, tx)
}

func (s *roomService) UpdateRoom(roomID string, validationData models.RoomCreateDTO, tx *gorm.DB) (*models.RoomResponseDTO, error) {
	room, err := s.repository.GetRoomByID(roomID, tx)
	if err != nil {
		return nil, err
	}

	// Update room details
	room.Name = validationData.Name
	room.Description = validationData.Description
	room.VariantID = validationData.Variant
	room.DetectiveID = validationData.Detective

	var files []models.File
	for _, fileHeader := range validationData.Files {
		// Save the file to the server
		filename := filepath.Base(fileHeader.Filename)
		filePath := "./uploads/" + filename
		if err := saveUploadedFile(fileHeader, filePath); err != nil {
			return nil, err
		}

		// Create File model
		file := models.File{URL: filePath}
		files = append(files, file)
	}

	room.Files = files
	updatedRoom, err := s.repository.UpdateRoom(room, tx)
	if err != nil {
		return nil, err
	}

	// Prepare response DTO
	responseFiles := make([]string, len(files))
	for i, file := range files {
		responseFiles[i] = file.URL
	}
	return &models.RoomResponseDTO{
		ID:          updatedRoom.ID,
		Name:        updatedRoom.Name,
		Description: updatedRoom.Description,
		Files:       responseFiles,
		Variant:     updatedRoom.Variant,
		Detective:   updatedRoom.Detective,
	}, nil
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
