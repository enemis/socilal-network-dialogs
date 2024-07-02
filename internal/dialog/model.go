package dialog

import (
	"github.com/google/uuid"
	"time"
)

type DialogType string

const (
	Direct    DialogType = "direct"
	Private   DialogType = "private"
	Published DialogType = "published"
)

type Participant struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Dialog struct {
	Id           uuid.UUID      `json:"id" db:"id"`
	Name         *string        `json:"name"`
	Type         DialogType     `json:"type" binding:"required"`
	Participants []*Participant `json:"participants"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	DeletedAt    *time.Time     `json:"deleted_at" db:"deleted_at"`
}

type Message struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	UserId    uuid.UUID  `json:"user_id" db:"user_id" binding:"required"`
	DialogId  uuid.UUID  `json:"dialog_id" db:"dialog_id" binding:"required"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Message   string     `json:"message" binding:"required"`
}
