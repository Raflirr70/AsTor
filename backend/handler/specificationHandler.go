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

type SpecificationHandler struct {
	specUsecase du.SpecificationUsecase
}

func NewSpecificationHandler(specUsecase du.SpecificationUsecase) *SpecificationHandler {
	return &SpecificationHandler{specUsecase: specUsecase}
}

func (h *SpecificationHandler) Create(c *gin.Context) {
	var req dto.CreateSpecification
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.specUsecase.Create(c.Request.Context(), adminID, req.ToEntity()); err != nil {
		response.Internal(c, "Failed to create specification")
		return
	}

	response.Created(c, nil, "Specification created successfully")
}

func (h *SpecificationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	var req dto.UpdateSpecification
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}

	spec := req.ToEntity()
	spec.ID = uint(id)
	adminID := c.GetUint("user_id")

	if err := h.specUsecase.Update(c.Request.Context(), adminID, spec); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to update specification")
		return
	}

	response.OK(c, nil, "Specification updated successfully")
}

func (h *SpecificationHandler) GetAll(c *gin.Context) {
	specs, err := h.specUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "Failed to fetch specifications")
		return
	}

	response.OK(c, dto.ToResponseSpecificationList(specs), "Success")
}
func (h *SpecificationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	spec, err := h.specUsecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "invalid to featch specification")
		return
	}
	response.OK(c, dto.ToResponseSpecification(spec), "Success")
}

func (h *SpecificationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	adminID := c.GetUint("user_id")

	if err := h.specUsecase.Delete(c.Request.Context(), adminID, &entity.Specifications{ID: uint(id)}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to delete specification")
		return
	}

	response.OK(c, nil, "Specification deleted successfully")
}
