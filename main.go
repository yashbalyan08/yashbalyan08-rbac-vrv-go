package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yashbalyan08/rbac-vrv-go/config"
	"github.com/yashbalyan08/rbac-vrv-go/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Connecting to DB...")
		log.Fatal("Error", err)
	}
	config.InitDB()
	defer config.DB.Close()

	router := gin.Default()
	router.Static("/static", "./static")

	router.LoadHTMLGlob("pages/*")
	routes.InitRoutes(router)

	router.Run(":10000")
}
