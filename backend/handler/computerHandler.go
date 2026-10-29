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

type ComputerHandler struct {
	computerUsecase du.ComputerUsecase
}

func NewComputerHandler(computerUsecase du.ComputerUsecase) *ComputerHandler {
	return &ComputerHandler{computerUsecase: computerUsecase}
}
func (h *ComputerHandler) Create(c *gin.Context) {
	var req dto.CreateComputer
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return // ← return wajib
	}
	input := du.UpdateComputerInput{Label: req.Label, Status: req.Status, ComponentIDs: req.Components}
	if err := h.computerUsecase.Create(c.Request.Context(), input); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		if errors.Is(err, entity.ErrInvalidInput) {
			response.BadRequest(c, err.Error())
			return
		}
		response.Internal(c, "failed to create computer")
		return
	}
	response.Created(c, nil, "computer created successfully")
}

func (h *ComputerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	var req dto.UpdateComputer
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return
	}
	input := du.UpdateComputerInput{Label: req.Label, Status: req.Status, ComponentIDs: req.Components}
	if err := h.computerUsecase.Update(c.Request.Context(), uint(id), input); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			response.NotFound(c, "computer not found")
			return
		} else if errors.Is(err, entity.ErrNoChanges) {
			response.BadRequest(c, "no changes detected")
			return
		} else {
			response.Internal(c, "failed updated computer")
			return
		}
	}
	response.OK(c, nil, "computer updated successfully")
}
func (h *ComputerHandler) GetAll(c *gin.Context) {
	comps, err := h.computerUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "invalide to find computer")
		return
	}
	res := make([]dto.ResponseComputer, 0, len(comps))
	for _, comp := range comps {
		res = append(res, dto.ToResponseComputer(&comp))
	}

	response.OK(c, res, "computers retrieved successfully")
}
func (h *ComputerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	adminID := c.GetUint("user_id")
	computer := &entity.Computers{ID: uint(id)}
	if err := h.computerUsecase.Delete(c.Request.Context(), adminID, computer); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to delete computer")
		return
	}
	response.OK(c, nil, "Computer deleted successfully")
}
func (h *ComputerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	computer, err := h.computerUsecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Internal(c, "invalid to find computer")
		return
	}
	response.OK(c, dto.ToResponseComputer(computer), "Success")
}
