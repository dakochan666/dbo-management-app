package models

import (
	"errors"
	"time"

	govalidator "github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not_null" json:"user_id" valid:"-"`
	ProductID uint           `gorm:"not_null" json:"product_id" valid:"-"`
	Amount    int            `gorm:"not_null" json:"amount" valid:"required~amount is required"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if _, errCreate := govalidator.ValidateStruct(struct {
		Amount int `valid:"required"`
	}{Amount: o.Amount}); errCreate != nil {
		return errors.New("validation error: " + errCreate.Error())
	}

	return nil
}
