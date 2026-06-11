package handler

import (
	"net/http"
	"strconv"

	"github.com/Richa3822/recipe-api/internal/model"
	"github.com/Richa3822/recipe-api/internal/service"
	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	service service.RecipeService
}

func NewRecipeHandler(service service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

// helper — extracts user_id from JWT context
func getUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	// JWT stores numbers as float64, convert to uint
	return uint(userID.(float64))
}

// POST /recipes
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.CreateRecipe(getUserID(c), &recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "recipe created successfully",
		"data":    recipe,
	})
}

// GET /recipes?search=pasta&page=1&limit=10
func (h *RecipeHandler) GetAllRecipes(c *gin.Context) {
	search := c.Query("search")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	recipes, total, err := h.service.GetAllRecipes(getUserID(c), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "recipes fetched successfully",
		"data":    recipes,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GET /recipes/:id
func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	recipe, err := h.service.GetRecipeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "recipe fetched successfully",
		"data":    recipe,
	})
}

// PUT /recipes/:id
func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updated model.Recipe
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	recipe, err := h.service.UpdateRecipe(getUserID(c), uint(id), &updated)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "recipe updated successfully",
		"data":    recipe,
	})
}

// DELETE /recipes/:id
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteRecipe(getUserID(c), uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe deleted successfully"})
}
