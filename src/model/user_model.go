package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type User struct {
	ID                 uuid.UUID           `gorm:"primaryKey;not null" json:"id"`
	Name               string              `gorm:"not null" json:"name"`
	Email              string              `gorm:"uniqueIndex;not null" json:"email"`
	Password           string              `gorm:"not null" json:"-"`
	Role               string              `gorm:"default:user;not null" json:"role"`
	WorkTime           int                 `gorm:"not null" json:"work_time"`
	ProjectPermissions []ProjectPermission `gorm:"many2many:user_project_permissions" json:"project_permissions"`
	CreatedAt          time.Time           `gorm:"autoCreateTime:milli" json:"-"`
	VerifiedEmail      bool                `gorm:"default:false;not null" json:"verified_email"`
	UpdatedAt          time.Time           `gorm:"autoUpdateTime:milli" json:"-"`
}


func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
