package handler

import (
	"astor/backend/entity"
	"astor/backend/handler/dto"
	"astor/backend/pkg/response"
	"astor/backend/pkg/validation"
	du "astor/backend/usecase/domain"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase du.AuthUsecase
}

func NewAuthHandler(authUsecase du.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// tampung data request
	var req dto.RequestRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validasi Gagal", validation.ParseErrors(err)...)
		return
	}
	// ubah menjadi entity
	user := req.ToEntity()
	// jalankan logika register
	if err := h.authUsecase.Register(c.Request.Context(), user); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			response.BadRequest(c, "NPM already exists")
			return
		}
		response.Internal(c, "An error occurred during registration")
		return
	}
	response.Created(c, nil, "Created Success")
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.RequestLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validasi Gagal", validation.ParseErrors(err)...)
		return
	}
	user, token, err := h.authUsecase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid Authorization: "+err.Error())
		return
	}

	response.OK(c, dto.ToResponseLogin(user, token), "Login Success")
}
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenValue, exists := c.Get("token")
	fmt.Printf("token : %s\nExist : %s", tokenValue, exists)
	if !exists {
		response.Unauthorized(c, entity.ErrUnauthorized.Error())
		return
	}

	token, ok := tokenValue.(string)
	if !ok || token == "" {
		response.Unauthorized(c, entity.ErrUnauthorized.Error())
		return
	}

	if err := h.authUsecase.Logout(c.Request.Context(), token); err != nil {
		response.Internal(c, "Failed to log out")
		return
	}

	response.OK(c, nil, "Logout successful")
}
