package main

import (
	"log"

	"github.com/Richa3822/recipe-api/db"
	"github.com/Richa3822/recipe-api/internal/handler"
	"github.com/Richa3822/recipe-api/internal/middleware"
	"github.com/Richa3822/recipe-api/internal/repository"
	"github.com/Richa3822/recipe-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	db.Connect()

	// Wire up auth
	userRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Wire up recipes
	recipeRepo := repository.NewRecipeRepository(db.DB)
	recipeService := service.NewRecipeService(recipeRepo)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	router := gin.Default()

	// Public routes — no token needed
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes — token required
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware()) // 👈 applies to all routes below
	{
		api.POST("/recipes", recipeHandler.CreateRecipe)
		api.GET("/recipes", recipeHandler.GetAllRecipes)
		api.GET("/recipes/:id", recipeHandler.GetRecipeByID)
		api.PUT("/recipes/:id", recipeHandler.UpdateRecipe)
		api.DELETE("/recipes/:id", recipeHandler.DeleteRecipe)
	}

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
