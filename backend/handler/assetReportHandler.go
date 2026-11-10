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

type AssetReportHandler struct {
	reportUsecase du.AssetReportUsecase
}

func NewAssetReportHandler(reportUsecase du.AssetReportUsecase) *AssetReportHandler {
	return &AssetReportHandler{reportUsecase: reportUsecase}
}

func (h *AssetReportHandler) Create(c *gin.Context) {
	var req dto.CreateReport
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.reportUsecase.CreateReport(c.Request.Context(), userID, req.ToEntity(), req.AssetCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "asset not found")
			return
		}
		response.Internal(c, "failed to create report")
		return
	}
	response.Created(c, nil, "report created successfully")
}

func (h *AssetReportHandler) GetAll(c *gin.Context) {
	reports, err := h.reportUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "failed to fetch reports")
		return
	}
	response.OK(c, dto.ToResponseReportList(reports), "reports retrieved successfully")
}

func (h *AssetReportHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	report, err := h.reportUsecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "report not found")
			return
		}
		response.Internal(c, "failed to fetch report")
		return
	}
	response.OK(c, dto.ToResponseReport(report), "report retrieved successfully")
}

func (h *AssetReportHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	var req dto.UpdateReport
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return
	}
	adminID := c.GetUint("user_id")
	if err := h.reportUsecase.UpdateReport(c.Request.Context(), adminID, req.ToEntity(uint(id)), req.Status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "report not found")
			return
		}
		if errors.Is(err, entity.ErrInvalidInput) {
			response.BadRequest(c, "status must be Baik or Buruk")
			return
		}
		response.Internal(c, "failed to update report")
		return
	}
	response.OK(c, nil, "report updated successfully")
}
