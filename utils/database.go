package utils

import (
	"fmt"
	"log"
	"os"

	models "geniale/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	var err error
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	println(dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to auto migrate User model: %v", err)
	}

	err = DB.AutoMigrate(&models.UserCode{})
	if err != nil {
		log.Fatalf("Failed to auto migrate UserCode model: %v", err)
	}

	err = DB.AutoMigrate(&models.Right{})
	if err != nil {
		log.Fatalf("Failed to auto migrate Right model: %v", err)
	}
	err = DB.AutoMigrate(&models.Role{})
	if err != nil {
		log.Fatalf("Failed to auto migrate Role model: %v", err)
	}
	err = DB.AutoMigrate(&models.RightRole{})
	if err != nil {
		log.Fatalf("Failed to auto migrate RightRole model: %v", err)
	}
	err = DB.AutoMigrate(&models.RoleUser{})
	if err != nil {
		log.Fatalf("Failed to auto migrate RoleUser model: %v", err)
	}

	err = DB.AutoMigrate(&models.Image{})
	if err != nil {
		log.Fatalf("Failed to auto migrate Image model: %v", err)
	}

	err = DB.AutoMigrate(&models.Room{}, &models.File{})
	if err != nil {
		log.Fatalf("Failed to auto migrate Room model: %v", err)
	}

	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);")

	// Create default rights if they don't exist
	DB.Exec("INSERT INTO rights (name) SELECT 'create' WHERE NOT EXISTS (SELECT 1 FROM rights WHERE name = 'create');")
	DB.Exec("INSERT INTO rights (name) SELECT 'read' WHERE NOT EXISTS (SELECT 1 FROM rights WHERE name = 'read');")
	DB.Exec("INSERT INTO rights (name) SELECT 'update' WHERE NOT EXISTS (SELECT 1 FROM rights WHERE name = 'update');")
	DB.Exec("INSERT INTO rights (name) SELECT 'delete' WHERE NOT EXISTS (SELECT 1 FROM rights WHERE name = 'delete');")

	// Create default roles if they don't exist
	DB.Exec("INSERT INTO roles (name) SELECT 'admin' WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'admin');")
	DB.Exec("INSERT INTO roles (name) SELECT 'user' WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'user');")
	// Create default right_roles if they don't exist
	DB.Exec("INSERT INTO right_roles (role_id, right_id) SELECT roles.id, rights.id FROM roles, rights WHERE roles.name = 'admin' AND rights.name = 'create' AND NOT EXISTS (SELECT 1 FROM right_roles WHERE role_id = roles.id AND right_id = rights.id);")
	DB.Exec("INSERT INTO right_roles (role_id, right_id) SELECT roles.id, rights.id FROM roles, rights WHERE roles.name = 'admin' AND rights.name = 'read' AND NOT EXISTS (SELECT 1 FROM right_roles WHERE role_id = roles.id AND right_id = rights.id);")
	DB.Exec("INSERT INTO right_roles (role_id, right_id) SELECT roles.id, rights.id FROM roles, rights WHERE roles.name = 'admin' AND rights.name = 'update' AND NOT EXISTS (SELECT 1 FROM right_roles WHERE role_id = roles.id AND right_id = rights.id);")
	DB.Exec("INSERT INTO right_roles (role_id, right_id) SELECT roles.id, rights.id FROM roles, rights WHERE roles.name = 'admin' AND rights.name = 'delete' AND NOT EXISTS (SELECT 1 FROM right_roles WHERE role_id = roles.id AND right_id = rights.id);")
	DB.Exec("INSERT INTO right_roles (role_id, right_id) SELECT roles.id, rights.id FROM roles, rights WHERE roles.name = 'user' AND rights.name = 'read' AND NOT EXISTS (SELECT 1 FROM right_roles WHERE role_id = roles.id AND right_id = rights.id);")

}
