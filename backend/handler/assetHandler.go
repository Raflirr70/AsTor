package handler

import (
	"astor/backend/entity"
	"astor/backend/handler/dto"
	"astor/backend/pkg/response"
	"astor/backend/pkg/validation"
	du "astor/backend/usecase/domain"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetHandler struct {
	assetUsecase du.AssetUsecase
}

func NewAssetHandler(assetUsecase du.AssetUsecase) *AssetHandler {
	return &AssetHandler{assetUsecase: assetUsecase}
}

func (h *AssetHandler) Create(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}
	asset, specs := req.ToEntity()
	adminID := c.GetUint("user_id")
	if err := h.assetUsecase.Create(c.Request.Context(), adminID, asset, specs); err != nil {
		if errors.Is(err, entity.ErrInvalidInput) {
			response.BadRequest(c, err.Error())
			return
		}
		log.Printf("Failed to create asset: %v", err)
		response.Internal(c, "Failed to create asset")
		return
	}
	response.Created(c, nil, "Asset created successfully")
}

func (h *AssetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", validation.ParseErrors(err)...)
		return
	}
	asset := req.ToEntity(uint(id))
	adminID := c.GetUint("user_id")
	if err := h.assetUsecase.Update(c.Request.Context(), adminID, asset); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		if errors.Is(err, entity.ErrInvalidInput) || errors.Is(err, entity.ErrNoChanges) {
			response.BadRequest(c, err.Error())
			return
		}
		response.Internal(c, "Failed to update asset")
		return
	}
	response.OK(c, nil, "Asset updated successfully")
}
func (h *AssetHandler) UpdateSpecification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	var req dto.UpdateAssetSpecification
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return
	}
	admin := c.GetUint("user_id")
	spec := req.ToEntity()
	spec.ID = uint(id)
	if err := h.assetUsecase.UpdateSpecification(c.Request.Context(), admin, spec); err != nil {
		response.Internal(c, "failed update asset specification")
		return
	}
	response.OK(c, nil, "Asset updated successfully")
}

func (h *AssetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	adminID := c.GetUint("user_id")
	asset := &entity.Assets{ID: uint(id)}
	if err := h.assetUsecase.Delete(c.Request.Context(), adminID, asset); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, entity.ErrNotFound.Error())
			return
		}
		response.Internal(c, "Failed to delete asset")
		return
	}

	response.OK(c, nil, "Asset deleted successfully")
}

func (h *AssetHandler) Checking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}

	asset, err := h.assetUsecase.FetchByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, entity.ErrNotFound.Error())
		return
	}

	asset.CheckedAt = time.Now()
	adminID := c.GetUint("user_id")

	if err := h.assetUsecase.Update(c.Request.Context(), adminID, asset); err != nil {
		response.Internal(c, "Failed to update checked_at")
		return
	}

	response.OK(c, nil, "Asset checked successfully")
}

func (h *AssetHandler) GetAll(c *gin.Context) {
	assets, err := h.assetUsecase.FetchAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "Failed to fetch assets")
		return
	}
	response.OK(c, dto.ToResponseAssetList(assets), "Success")
}
