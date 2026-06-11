package service

import (
	"errors"

	"github.com/Richa3822/recipe-api/internal/model"
	"github.com/Richa3822/recipe-api/internal/repository"
)

type RecipeService interface {
	CreateRecipe(userID uint, recipe *model.Recipe) error
	GetAllRecipes(userID uint, search string, page, limit int) ([]model.Recipe, int64, error)
	GetRecipeByID(id uint) (*model.Recipe, error)
	UpdateRecipe(userID uint, id uint, updated *model.Recipe) (*model.Recipe, error)
	DeleteRecipe(userID uint, id uint) error
}

type recipeService struct {
	repo repository.RecipeRepository
}

func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{repo: repo}
}

func (s *recipeService) CreateRecipe(userID uint, recipe *model.Recipe) error {
	if recipe.Title == "" {
		return errors.New("title is required")
	}
	if recipe.Ingredients == "" {
		return errors.New("ingredients are required")
	}
	if recipe.Steps == "" {
		return errors.New("steps are required")
	}
	recipe.UserID = userID // 👈 attach the logged-in user's ID
	return s.repo.Create(recipe)
}

func (s *recipeService) GetAllRecipes(userID uint, search string, page, limit int) ([]model.Recipe, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetAll(userID, search, offset, limit)
}

func (s *recipeService) GetRecipeByID(id uint) (*model.Recipe, error) {
	recipe, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("recipe not found")
	}
	return recipe, nil
}

func (s *recipeService) UpdateRecipe(userID uint, id uint, updated *model.Recipe) (*model.Recipe, error) {
	recipe, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("recipe not found")
	}

	// Ownership check — is this recipe yours?
	if recipe.UserID != userID {
		return nil, errors.New("you are not authorized to update this recipe")
	}

	if updated.Title != "" {
		recipe.Title = updated.Title
	}
	if updated.Description != "" {
		recipe.Description = updated.Description
	}
	if updated.Ingredients != "" {
		recipe.Ingredients = updated.Ingredients
	}
	if updated.Steps != "" {
		recipe.Steps = updated.Steps
	}
	if updated.CookTime != 0 {
		recipe.CookTime = updated.CookTime
	}
	if updated.Servings != 0 {
		recipe.Servings = updated.Servings
	}

	err = s.repo.Update(recipe)
	return recipe, err
}

func (s *recipeService) DeleteRecipe(userID uint, id uint) error {
	err := s.repo.Delete(id, userID)
	if err != nil {
		return errors.New("recipe not found or not authorized")
	}
	return nil
}
