package repository

import (
	"github.com/Richa3822/recipe-api/internal/model"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(recipe *model.Recipe) error
	GetAll(userID uint, search string, offset, limit int) ([]model.Recipe, int64, error)
	GetByID(id uint) (*model.Recipe, error)
	Update(recipe *model.Recipe) error
	Delete(id, userID uint) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *recipeRepository) GetAll(userID uint, search string, offset, limit int) ([]model.Recipe, int64, error) {
	var recipes []model.Recipe
	var total int64

	query := r.db.Model(&model.Recipe{}).Where("user_id = ?", userID) // 👈 only this user's recipes

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR ingredients ILIKE ?", searchTerm, searchTerm)
	}

	query.Count(&total)
	result := query.Offset(offset).Limit(limit).Find(&recipes)

	return recipes, total, result.Error
}

func (r *recipeRepository) GetByID(id uint) (*model.Recipe, error) {
	var recipe model.Recipe
	result := r.db.First(&recipe, id)
	return &recipe, result.Error
}

func (r *recipeRepository) Update(recipe *model.Recipe) error {
	return r.db.Save(recipe).Error
}

// DELETE — only deletes if both id AND user_id match
func (r *recipeRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Recipe{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
