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

type LoanHandler struct {
	loanUsecase du.LoanUsecase
}

func NewLoanHandler(loanUsecase du.LoanUsecase) *LoanHandler {
	return &LoanHandler{loanUsecase: loanUsecase}
}

func (h *LoanHandler) Create(c *gin.Context) {
	var req dto.CreateLoan
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", validation.ParseErrors(err)...)
		return
	}
	userID := c.GetUint("user_id")

	assetCodes := make([]string, 0, len(req.Items))
	quantities := make([]int, 0, len(req.Items))

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			item.Quantity = 1
		}
		assetCodes = append(assetCodes, item.AssetCode)
		quantities = append(quantities, item.Quantity)
	}

	if err := h.loanUsecase.Create(c.Request.Context(), userID, req.ToEntity(), quantities, assetCodes); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "asset not found")
			return
		}
		response.Internal(c, "failed to create loan")
		return
	}
	response.Created(c, nil, "loan created successfully")
}

func (h *LoanHandler) GetAll(c *gin.Context) {
	loans, err := h.loanUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Internal(c, "failed to fetch loans")
		return
	}
	res := make([]dto.ResponseLoan, 0, len(loans))
	for _, l := range loans {
		res = append(res, dto.ToResponseLoan(&l))
	}
	response.OK(c, res, "loans retrieved successfully")
}

func (h *LoanHandler) GetByStatus(c *gin.Context) {
	status := c.Query("status")
	loans, err := h.loanUsecase.GetByStatus(c.Request.Context(), status)
	if err != nil {
		response.Internal(c, "failed to fetch loans by status")
		return
	}
	res := make([]dto.ResponseLoan, 0, len(loans))
	for _, l := range loans {
		res = append(res, dto.ToResponseLoan(&l))
	}
	response.OK(c, res, "loans retrieved successfully")
}

func (h *LoanHandler) Approved(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	userID := c.GetUint("user_id")
	loan := &entity.Loans{ID: uint(id)}
	if err := h.loanUsecase.Approved(c.Request.Context(), userID, loan); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "loan not found")
			return
		}
		response.Internal(c, "failed to approve loan")
		return
	}
	response.OK(c, nil, "loan approved successfully")
}

func (h *LoanHandler) Returned(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, entity.ErrInvalidID.Error())
		return
	}
	status := c.Query("status")

	if status != "" {
		status = "Baik"
	}
	userID := c.GetUint("user_id")
	loan := &entity.Loans{ID: uint(id)}
	if err := h.loanUsecase.Returned(c.Request.Context(), userID, loan, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "loan not found")
			return
		}
		response.Internal(c, "failed to return loan")
		return
	}
	response.OK(c, nil, "loan approved successfully")
}
