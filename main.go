package main

import (
	"log"
	"myproject/databases"
	_ "myproject/docs"
	"myproject/handlers"
	"myproject/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Pet Management API
// @version 1.0
// @description API для подбора домашних животных
// @host localhost:8080
// // @BasePath /v1
func main() {

	router := gin.Default()
	database, err := databases.Connect()
	petHandler := handlers.CreatePetHandler(database)
	userHandler := handlers.CreateUserHandler(database)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect()

	// Публичные маршруты
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", userHandler.Login)
	router.POST("/register", userHandler.Register)
	router.GET("/pets", petHandler.GetPet)
	router.GET("/pets/:id", petHandler.GetPets)

	// Защищенные маршруты (только для админов)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middlewares.Authenticate("admin"))
	{
		adminRoutes.POST("/pets", petHandler.CreatePet)
		adminRoutes.PUT("/pets", petHandler.UpdatePet)
		adminRoutes.DELETE("/pets", petHandler.DeletePet)
	}

	router.Run()
}
