package main

import (
	"fmt"
	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/services"
	"log"
)

func main() {
	// Load config
	config.LoadConfig()
	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Failed to load config")
	}

	// Init DB
	if err := models.InitDB(cfg); err != nil {
		log.Fatal(err)
	}

	// Init Services
	services.SetupAdminAdminService()

	adminSvc := services.GetAdminAdminService()

	username := "devadmin"
	password := "password123"
	email := "dev@greenrideafrica.com"
	role := "super_admin"

	admin, errCode := adminSvc.CreateAdmin(username, email, password, role, "Dev", "Developer", "SYSTEM")
	if errCode != "" {
		fmt.Printf("Error creating admin: %s. It might already exist.\n", errCode)
		if errCode == "4001" { // UserAlreadyExists
			// Reset password
			existing := adminSvc.GetAdminByUsername(username)
			adminSvc.ResetPassword(existing.AdminID, password, "SYSTEM")
			fmt.Println("Password reset for devadmin.")
		}
	} else {
		fmt.Printf("Dev admin created: %s\n", admin.AdminID)
	}

	// Force must_change_password to 0
	if admin == nil {
		admin = adminSvc.GetAdminByUsername(username)
	}
	if admin != nil {
		models.GetDB().Model(admin).Update("must_change_password", 0)
		fmt.Println("Set must_change_password to 0 for devadmin.")
	}
}
