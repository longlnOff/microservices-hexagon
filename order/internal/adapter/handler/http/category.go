package http

import (
	"net/http"

)

// CategoryHandler represents the HTTP handler for category-related requests
type CategoryHandler struct {
	svc port.CategoryService
	logger port.LoggerRepository
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(svc port.CategoryService, logger port.LoggerRepository) *CategoryHandler {
	return &CategoryHandler{
		svc,
		logger,
	}
}

// createCategoryRequest represents a request body for creating a new category
type createCategoryRequest struct {
	Code 		string `json:"code" binding:"required" example:"large_language_model"`
	Name 		string `json:"name" binding:"required" example:"Large Language Model"`
	Description string `json:"description" binding:"required" example:"A large language model is a machine learning model that can understand and generate human-like language."`
}



// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	create a new category with name
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			createCategoryRequest	body		createCategoryRequest	true	"Create category request"
//	@Success		200						{object}	categoryResponse		"Category created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/categories [post]
//	@Security		BearerAuth
func (ch *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	category := domain.Category{
		Code: req.Code,
		Name: req.Name,
		Description: req.Description,
	}

	returned, err := ch.svc.CreateCategory(r.Context(), &category)
	if err != nil {
		return
	}

	if err := jsonResponse(w, http.StatusCreated, returned); err != nil {
		internalServerError(w, r, err)
		return
	}
}
