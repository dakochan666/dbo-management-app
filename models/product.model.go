package models

import (
	"time"

	govalidator "github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not_null;type:varchar(255)" json:"name" valid:"required~name is required"`
	Stock       int            `gorm:"not_null" json:"stock" valid:"required~stock is required"`
	Description string         `gorm:"not_null;type:varchar(255)" json:"description" valid:"required~description is required"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Orders []Order `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)
	if errCreate != nil {
		return errCreate
	}

	return
}
