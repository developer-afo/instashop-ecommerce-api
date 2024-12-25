// seed/seed.go
package seed

import (
	"fmt"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type SeederInterface interface {
	Seed()
}

type seeder struct {
	dbConn database.DatabaseInterface
}

func NewSeeder(dbConn database.DatabaseInterface) SeederInterface {
	return &seeder{dbConn: dbConn}
}

func (s *seeder) Seed() {
	s.SeedAdmin()
	s.SeedOrderStatuses()
}

func (s *seeder) SeedAdmin() {
	hashing := helper.NewHashing()
	adminEmail := "admin@instashop.com.ng"
	testEmail := "test@email.com"
	hashedPassword, err := hashing.HashPassword("password")
	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return
	}

	adminExists := s.dbConn.Connection().Where("email = ?", adminEmail).First(&models.User{}).RowsAffected > 0
	testUserExists := s.dbConn.Connection().Where("email = ?", testEmail).First(&models.User{}).RowsAffected > 0

	// Create Admin User
	adminUser := models.User{
		FirstName:       "Instashop",
		LastName:        "Admin",
		Email:           adminEmail,
		IsEmailVerified: true,
		Password:        hashedPassword,
		Role:            "admin",
	}

	adminUser.Prepare()

	// Create test user
	testUser := models.User{
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@email.com",
		IsEmailVerified: true,
		Password:        hashedPassword,
		Role:            "customer",
	}

	testUser.Prepare()

	if adminExists {
		fmt.Println("Admin already exists in the database. Skipping seeding...")
	} else {
		if err := s.dbConn.Connection().Create(&adminUser).Error; err != nil {
			fmt.Println("Failed to create admin user:", err)
		}
		fmt.Println("Admin user created successfully.")
	}

	if testUserExists {
		fmt.Println("Test user already exists in the database. Skipping seeding...")
	} else {
		if err := s.dbConn.Connection().Create(&testUser).Error; err != nil {
			fmt.Println("Failed to create test user:", err)
		}
		fmt.Println("Test user created successfully.")
	}

}

func (s *seeder) SeedOrderStatuses() {
	orderStatuses := []models.OrderStatus{
		{Name: "Order Placed", ShortName: "order_placed"},
		{Name: "Awaiting Confirmation", ShortName: "awaiting_confirmation"},
		{Name: "Order Processing", ShortName: "order_processing"},
		{Name: "Out for delivery", ShortName: "out_for_delivery"},
		{Name: "Delivered", ShortName: "delivered"},
		{Name: "Cancelled", ShortName: "cancelled"},
	}

	for _, orderStatus := range orderStatuses {
		statusExists := s.dbConn.Connection().Where("short_name = ?", orderStatus.ShortName).First(&models.OrderStatus{}).RowsAffected > 0
		if statusExists {
			fmt.Printf("%s Order status already exists in the database. Skipping seeding...\n", orderStatus.Name)
		} else {
			orderStatus.Prepare()
			if err := s.dbConn.Connection().Create(&orderStatus).Error; err != nil {
				fmt.Println("Failed to create order status:", err)
			}
			fmt.Println("Order statuses created successfully.")
		}
	}

}
