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

type LaboratoryHandler struct {
	laboratoryUsecase du.LaboratoryUsecase
}

func NewLaboratoryHandler(laboratoryUsecase du.LaboratoryUsecase) *LaboratoryHandler {
	return &LaboratoryHandler{laboratoryUsecase: laboratoryUsecase}
}
func (h *LaboratoryHandler) Create(c *gin.Context) {
	var req dto.CreateLaboratory
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}
	adminID := c.GetUint("user_id")
	if err := h.laboratoryUsecase.Create(c.Request.Context(), adminID, req.ToEntity()); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			response.BadRequest(c, "Laboratory name already exists")
			return
		}
		response.Internal(c, "failed to create laboratory")
		return
	}
	response.Created(c, nil, "Created laboratory")
}
func (h *LaboratoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	var req dto.UpdateLaboratory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}
	lab := req.ToEntity()
	lab.ID = uint(id)
	adminID := c.GetUint("user_id")
	if err := h.laboratoryUsecase.Update(c.Request.Context(), adminID, lab); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			response.BadRequest(c, "Laboratory name already exists")
			return
		}
		response.Internal(c, "failed to update laboratory")
		return
	}
	response.OK(c, nil, "Updated laboratory")
}
func (h *LaboratoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	adminID := c.GetUint("user_id")

	if err := h.laboratoryUsecase.Delete(c.Request.Context(), adminID, &entity.Laboratories{ID: uint(id)}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "An error occurred during deletion")
		return
	}

	response.OK(c, nil, "Laboratory deleted successfully")
}
func (h *LaboratoryHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	adminID := c.GetUint("user_id")

	if err := h.laboratoryUsecase.Restore(c.Request.Context(), adminID, &entity.Laboratories{ID: uint(id)}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			response.BadRequest(c, "Laboratory is not deleted")
			return
		}
		response.Internal(c, "An error occurred during restore")
		return
	}

	response.OK(c, nil, "Laboratory restored successfully")
}
func (h *LaboratoryHandler) GetAll(c *gin.Context) {
	labs, err := h.laboratoryUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "failed to fetch laboratories")
		return
	}
	response.OK(c, dto.ToResponseLaboratoriesForAdminList(labs), "success")
}
