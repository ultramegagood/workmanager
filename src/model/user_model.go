package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Email             string    `gorm:"uniqueIndex;not null" json:"email"`
	Password          string    `gorm:"not null" json:"-"`
	Role              string    `gorm:"default:user;not null" json:"role"`
	WorkTime          int       `gorm:"not null" json:"work_time"`
	UserGroupIDs      []uuid.UUID `gorm:"type:uuid[]" json:"user_group_ids"`
	ProjectPermissions []ProjectPermission `gorm:"foreignKey:UserID" json:"project_permissions"`
	CreatedAt         time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt         time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
