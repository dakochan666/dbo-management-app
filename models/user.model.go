package models

import (
	"time"

	govalidator "github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not_null;type:varchar(255)" json:"name" valid:"required~name is required"`
	Email     string         `gorm:"not_null;type:varchar(255)" json:"email" valid:"required~email is required"`
	Password  string         `gorm:"not_null;type:varchar(255)" json:"password" valid:"required~password is required"`
	Role      string         `gorm:"not_null;type:varchar(255)" json:"role" valid:"required~role is required"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

type ReqUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		return errCreate
	}

	return
}
