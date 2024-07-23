package models

type Invitation struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	Email           string `json:"email" gorm:"unique;not null"`
	InvitationToken string `json:"invitation_token"`
	EntityID        uint   `json:"entity_id"`
	Entity          Entity `json:"entity"`
	Status          string `json:"status"`
}

type InvitationCreateDTO struct {
	Email    string `json:"email" validate:"required,email"`
	EntityID uint   `json:"entity_id" validate:"required"`
}

type InvitationCreateResponseDTO struct {
	Invitation Invitation `json:"invitation"`
	User       User       `json:"user"`
}

type InvitationResponseDTO struct {
	Message    string                      `json:"message"`
	Invitation InvitationCreateResponseDTO `json:"invitation"`
}

type InvitationListResponseDTO struct {
	Message     string       `json:"message"`
	Invitations []Invitation `json:"invitations"`
}

type InvitationUpdateDTO struct {
	Status string `json:"status" validate:"required"`
}

type InvitationAcceptDTO struct {
	Email           string `json:"email" validate:"required,email"`
	InvitationToken string `json:"invitation_token" validate:"required"`
}

type InvitationAcceptResponseDTO struct {
	Message string `json:"message"`
}

type InvitationDeleteDTO struct {
	Email           string `json:"email" validate:"required,email"`
	InvitationToken string `json:"invitation_token" validate:"required"`
}

type InvitationDeleteResponseDTO struct {
	Message string `json:"message"`
}

type InvitationResendDTO struct {
	Email           string `json:"email" validate:"required,email"`
	InvitationToken string `json:"invitation_token" validate:"required"`
}

type InvitationResendResponseDTO struct {
	Message string `json:"message"`
}

type InvitationCancelDTO struct {
	Email           string `json:"email" validate:"required,email"`
	InvitationToken string `json:"invitation_token" validate:"required"`
}

type InvitationCancelResponseDTO struct {
	Message string `json:"message"`
}
