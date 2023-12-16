package seeder

import (
	"dbo-management-app/models"
	"dbo-management-app/service"
	"log"

	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db}
}

func (s *Seeder) SeedAdmin() {
	var count int64
	s.db.Model(&models.User{}).Where("name = ?", "Admin").Count(&count)

	if count == 0 {
		encryptedPassword, err := service.HashPassword("admin123")
		if err != nil {
			log.Fatalf("Error seeding user: %v", err.Error())
			return
		}
		admins := []models.User{
			{
				Name:     "Admin",
				Email:    "admin@mail.com",
				Password: string(encryptedPassword),
				Role:     "admin",
			},
		}

		s.db.AutoMigrate(&models.User{})
		for _, admin := range admins {
			result := s.db.Create(&admin)
			if result.Error != nil {
				log.Fatalf("Error seeding user: %v", result.Error)
			}
		}
	}

}
