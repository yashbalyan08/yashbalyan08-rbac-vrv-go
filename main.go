package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yashbalyan08/rbac-vrv-go/config"
	"github.com/yashbalyan08/rbac-vrv-go/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.InitDB()
	defer config.DB.Close()

	router := gin.Default()
	routes.InitRoutes(router)

	router.Run(":8080")
}
