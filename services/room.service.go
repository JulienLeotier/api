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

func saveFiles(fileHeaders []*multipart.FileHeader) ([]models.File, error) {
	var files []models.File
	for _, fileHeader := range fileHeaders {
		filename := filepath.Base(fileHeader.Filename)
		filePath := "./uploads/" + filename
		if err := saveUploadedFile(fileHeader, filePath); err != nil {
			return nil, err
		}
		file := models.File{URL: filePath}
		files = append(files, file)
	}
	return files, nil
}

func (s *roomService) CreateRoom(validationData models.RoomCreateDTO, tx *gorm.DB) (*models.RoomResponseDTO, error) {
	enigmaFiles, err := saveFiles(validationData.EnigmaFiles)
	if err != nil {
		return nil, err
	}
	revelationFiles, err := saveFiles(validationData.RevelationFiles)
	if err != nil {
		return nil, err
	}
	introFiles, err := saveFiles(validationData.IntroFiles)
	if err != nil {
		return nil, err
	}
	outroFiles, err := saveFiles(validationData.OutroFiles)
	if err != nil {
		return nil, err
	}
	ambianceMusicFiles, err := saveFiles(validationData.AmbianceMusicFiles)
	if err != nil {
		return nil, err
	}

	room := &models.Room{
		Name:               validationData.Name,
		Description:        validationData.Description,
		EnigmaFiles:        enigmaFiles,
		RevelationFiles:    revelationFiles,
		IntroFiles:         introFiles,
		OutroFiles:         outroFiles,
		AmbianceMusicFiles: ambianceMusicFiles,
		VariantID:          validationData.Variant,
		DetectiveID:        validationData.Detective,
	}
	roomSave, err := s.repository.CreateRoom(room, tx)
	if err != nil {
		return nil, err
	}

	return s.buildRoomResponse(roomSave), nil
}

func (s *roomService) GetRoomByID(roomID string, tx *gorm.DB) (*models.RoomResponseDTO, error) {
	room, err := s.repository.GetRoomByID(roomID, tx)
	if err != nil {
		return nil, err
	}
	return s.buildRoomResponse(room), nil
}

func (s *roomService) GetAllRooms(tx *gorm.DB) ([]models.RoomResponseDTO, error) {
	rooms, err := s.repository.GetAllRooms(tx)
	if err != nil {
		return nil, err
	}

	var responseRooms []models.RoomResponseDTO
	for _, room := range rooms {
		responseRooms = append(responseRooms, *s.buildRoomResponse(&room))
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

	room.Name = validationData.Name
	room.Description = validationData.Description
	room.VariantID = validationData.Variant
	room.DetectiveID = validationData.Detective

	enigmaFiles, err := saveFiles(validationData.EnigmaFiles)
	if err != nil {
		return nil, err
	}
	revelationFiles, err := saveFiles(validationData.RevelationFiles)
	if err != nil {
		return nil, err
	}
	introFiles, err := saveFiles(validationData.IntroFiles)
	if err != nil {
		return nil, err
	}
	outroFiles, err := saveFiles(validationData.OutroFiles)
	if err != nil {
		return nil, err
	}
	ambianceMusicFiles, err := saveFiles(validationData.AmbianceMusicFiles)
	if err != nil {
		return nil, err
	}

	room.EnigmaFiles = enigmaFiles
	room.RevelationFiles = revelationFiles
	room.IntroFiles = introFiles
	room.OutroFiles = outroFiles
	room.AmbianceMusicFiles = ambianceMusicFiles

	updatedRoom, err := s.repository.UpdateRoom(room, tx)
	if err != nil {
		return nil, err
	}

	return s.buildRoomResponse(updatedRoom), nil
}

func (s *roomService) buildRoomResponse(room *models.Room) *models.RoomResponseDTO {
	return &models.RoomResponseDTO{
		ID:                 room.ID,
		Name:               room.Name,
		Description:        room.Description,
		EnigmaFiles:        extractFileURLs(room.EnigmaFiles),
		RevelationFiles:    extractFileURLs(room.RevelationFiles),
		IntroFiles:         extractFileURLs(room.IntroFiles),
		OutroFiles:         extractFileURLs(room.OutroFiles),
		AmbianceMusicFiles: extractFileURLs(room.AmbianceMusicFiles),
		Variant:            room.Variant,
		Detective:          room.Detective,
	}
}

func extractFileURLs(files []models.File) []string {
	fileURLs := make([]string, len(files))
	for i, file := range files {
		fileURLs[i] = file.URL
	}
	return fileURLs
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
