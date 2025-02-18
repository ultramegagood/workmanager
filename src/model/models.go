package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Вспомогательная функция для генерации UUID перед созданием
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

// ======= Пользователь =======
type User struct {
	BaseModel
	Name               string              `gorm:"not null" json:"name"`
	Email              string              `gorm:"uniqueIndex;not null" json:"email"` // Уникальный индекс для email	Password           string              `gorm:"not null" json:"-"`
	Role               string              `gorm:"default:user;not null" json:"role"`
	WorkTime           int                 `gorm:""json:"work_time"`
	Password           string              `gorm:"not null"json:"password"`
	VerifiedEmail      bool                `gorm:"default:false;not null" json:"verified_email"`
	ProjectPermissions []ProjectPermission `gorm:"foreignKey:UserID" json:"project_permissions"`
	Projects           []Project           `gorm:"many2many:project_users;" json:"projects"`
	Tasks              []Task              `gorm:"many2many:task_users;" json:"tasks"`
	Groups             []UserGroup         `gorm:"many2many:user_group_users"`
}
type Token struct {
	BaseModel
	Token   string    `gorm:"not null" json:"token"`
	UserID  uuid.UUID `gorm:"not null" json:"user_id"`
	Type    string    `gorm:"not null" json:"type"`
	Expires time.Time `gorm:"not null" json:"expires"`
	User    User      `gorm:"foreignKey:UserID;onDelete:CASCADE"`
}

// ======= Проекты и группы =======

type Project struct {
	BaseModel
	Title      string      `gorm:"not null" json:"title"`
	Users      []User      `gorm:"many2many:project_users;" json:"users"`
	UserGroups []UserGroup `gorm:"many2many:project_user_groups;" json:"user_groups"`
}

type ProjectUser struct {
	ProjectID uuid.UUID `gorm:"primaryKey" json:"project_id"`
	UserID    uuid.UUID `gorm:"primaryKey" json:"user_id"`
}

type Section struct {
	BaseModel
	Title     string     `gorm:"not null" json:"title"`
	ProjectID uuid.UUID  `gorm:"not null" json:"project_id"`
	UserGroup *uuid.UUID `json:"user_group,omitempty"`
	Project   Project    `gorm:"foreignKey:ProjectID;onDelete:CASCADE"`
	Tasks     []Task     `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	Order     int        `gorm:"not null;default:0" json:"order"`
}

// ======= Задачи =======
type Task struct {
	BaseModel
	ProjectID     uuid.UUID   `json:"project_id"`
	Project       Project     `gorm:"foreignKey:ProjectID;onDelete:CASCADE"`
	Title         string      `gorm:"not null" json:"title"`
	Description   string      `json:"description"`
	UserGroup     *uuid.UUID  `json:"user_group,omitempty"`
	Status        string      `json:"status"`
	Priority      string      `json:"priority"`
	DueDate       *time.Time  `json:"due_date,omitempty"`
	SectionID     uuid.UUID   `gorm:"not null" json:"section_id"`
	UserSectionID *uuid.UUID  `json:"user_section_id,omitempty"` // Новая связь с UserSection
	AssignedTo    *uuid.UUID  `json:"assigned_to,omitempty"`
	ParentTaskID  *uuid.UUID  `json:"parent_task_id,omitempty"`
	EstimatedTime int         `json:"estimated_time"`
	SpentTime     int         `json:"spent_time"`
	Users         []User      `gorm:"many2many:task_users;" json:"users"`
	UserGroups    []UserGroup `gorm:"many2many:task_user_groups;" json:"user_groups"`
}

// ======= Секции пользователя =======
type UserSection struct {
	BaseModel
	Title    string    `gorm:"not null" json:"title"`
	UserID   uuid.UUID `gorm:"not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID;onDelete:CASCADE"`
	Tasks    []Task    `gorm:"foreignKey:UserSectionID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	Order    int       `gorm:"not null;default:0" json:"order"`
}

type TaskUser struct {
	TaskID uuid.UUID `gorm:"primaryKey" json:"task_id"`
	UserID uuid.UUID `gorm:"primaryKey" json:"user_id"`
}

// ======= История задач =======

type TaskHistory struct {
	TaskID    uuid.UUID `gorm:"not null" json:"task_id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"not null" json:"action"`
	Timestamp time.Time `gorm:"autoCreateTime:milli" json:"timestamp"`
}

// ======= Комментарии =======
type Comment struct {
	BaseModel
	TaskID    uuid.UUID  `gorm:"not null" json:"task_id"`
	Task      Task       `gorm:"foreignKey:TaskID;onDelete:CASCADE"`
	UserID    uuid.UUID  `gorm:"not null" json:"user_id"`                        // Добавляем ID пользователя
	User      User       `gorm:"foreignKey:UserID;onDelete:CASCADE" json:"user"` // Связь с пользователем
	Body      string     `gorm:"not null" json:"body"`
	CitateID  *uuid.UUID `json:"citate_id,omitempty"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
	IsEdited  bool       `gorm:"default:false" json:"is_edited"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ======= Вложения (Attachments) =======

type Attachment struct {
	BaseModel
	TaskID          uuid.UUID  `gorm:"not null" json:"task_id"`
	Task            Task       `gorm:"foreignKey:TaskID;onDelete:CASCADE"`
	UserID          uuid.UUID  `gorm:"not null" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID;onDelete:CASCADE"`
	URL             string     `gorm:"not null" json:"url"`
	Type            string     `gorm:"not null" json:"type"`
	Size            int        `gorm:"not null" json:"size"`
	LinkedTaskID    *uuid.UUID `json:"linked_task_id,omitempty"`
	LinkedCommentID *uuid.UUID `json:"linked_comment_id,omitempty"`
}

// ======= Логи аудита (Audit Logs) =======

type AuditLog struct {
	BaseModel
	TaskID     uuid.UUID `gorm:"not null" json:"task_id"`
	Body       string    `gorm:"not null" json:"body"`
	ActionType string    `gorm:"not null" json:"action_type"`
	EntityType string    `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"not null" json:"entity_id"`
}

// ======= Разрешения =======

type ProjectPermission struct {
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;onDelete:CASCADE"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID;onDelete:CASCADE"`
	Role      string    `gorm:"not null" json:"role"`
}

type RolePermission struct {
	RoleID uuid.UUID `gorm:"not null" json:"role_id"`
	UserID uuid.UUID `gorm:"not null" json:"user_id"`
}

// ======= Группы пользователей =======

type UserGroup struct {
	BaseModel
	TeamTitle string    `gorm:"not null" json:"team_title"`
	OwnerID   uuid.UUID `gorm:"" json:"owner_id"` // Автор группы
	Owner     User      `gorm:"foreignKey:OwnerID"`
	Users     []User    `gorm:"many2many:user_group_users;foreignKey:ID;joinForeignKey:UserGroupID;References:ID;joinReferences:UserID" json:"users"`
	Projects  []Project `gorm:"many2many:project_user_groups;" json:"projects"`
	Tasks     []Task    `gorm:"many2many:task_user_groups;" json:"tasks"`
}

// ======= Роли пользователей в проекте =======

type UserProjectRole struct {
	BaseModel
	Name      string    `gorm:"not null" json:"name"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID;onDelete:CASCADE"`
}
