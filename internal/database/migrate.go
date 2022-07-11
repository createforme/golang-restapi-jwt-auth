package database

import (
	"github.com/createforme/golang-restapi-jwt-auth/internal/user"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	// MigrateDB - migrates our database and create our comment tables
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return err
	}
	return nil
}
