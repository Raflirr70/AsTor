package handler

import (
	"astor/backend/entity"
	"astor/backend/handler/dto"
	"astor/backend/pkg/response"
	"astor/backend/pkg/validation"
	du "astor/backend/usecase/domain"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	categoryUsecase du.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase du.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.categoryUsecase.Create(c.Request.Context(), adminID, req.ToEntity(), req.SpecificationIDs); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			response.BadRequest(c, "Category already exists")
			return
		}
		response.Internal(c, "Failed to create category")
		return
	}

	response.Created(c, nil, "Category created successfully")
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	var req dto.UpdateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}

	category := req.ToEntity()
	category.ID = uint(id)
	adminID := c.GetUint("user_id")

	if err := h.categoryUsecase.Update(c.Request.Context(), adminID, category, req.SpecificationIDs); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to update category")
		return
	}

	response.OK(c, nil, "Category updated successfully")
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	adminID := c.GetUint("user_id")

	if err := h.categoryUsecase.Delete(c.Request.Context(), adminID, &entity.Categories{ID: uint(id)}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to delete category")
		return
	}

	response.OK(c, nil, "Category deleted successfully")
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "Failed to fetch categories")
		return
	}

	response.OK(c, dto.ToResponseCategoryList(categories), "Success")
}
