package main

import (
	"fmt"
	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
)

func main() {
	// 1. Load config
	fmt.Println("Loading configuration...")
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize DB
	cfg := config.Get()
	fmt.Println("Connecting to database...")
	if err := models.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2.5. Run migrations
	fmt.Println("Running migrations...")
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 3. Setup services
	services.SetupService()

	// 4. Admin Credentials
	username := "admin"
	email := "admin@greenrideafrica.com"
	password := "admin123" // Recommended: change this after first login
	role := models.AdminRoleSuperAdmin

	adminSvc := services.GetAdminAdminService()

	// Check if already exists
	if adminSvc.IsUsernameExists(username) {
		fmt.Printf("Admin with username '%s' already exists. Resetting password to 'admin123'...\n", username)
		admin := adminSvc.GetAdminByUsername(username)
		errCode := adminSvc.ResetPassword(admin.AdminID, password, "SYSTEM_RECOVERY")
		if errCode != protocol.Success {
			log.Fatalf("Failed to reset password. Error Code: %v", errCode)
		}
		fmt.Println("Success! Password updated to: admin123")
	} else {
		fmt.Printf("Creating new Super Admin...\nUsername: %s\nPassword: %s\n", username, password)
		_, errCode := adminSvc.CreateAdmin(username, email, password, role, "IT", "System Admin", "SYSTEM_RECOVERY")
		if errCode != protocol.Success {
			log.Fatalf("Failed to create admin. Error Code: %v", errCode)
		}
		fmt.Println("Success! Admin created successfully.")
	}

	fmt.Println("\nYou can now login at admin.greenrideafrica.com using these credentials.")
}
