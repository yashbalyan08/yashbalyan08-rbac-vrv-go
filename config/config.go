package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func initDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("DB connected: ", DB)
}
