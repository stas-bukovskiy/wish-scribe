package handler

import (
	"github.com/gin-gonic/gin"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"net/http"
)

type SingUpRequest struct {
	Name     string `json:"name" binging:"required"`
	Email    string `json:"email" binging:"required"`
	Password string `json:"password" binging:"required"`
}

type SingInRequest struct {
	Email    string `json:"email" binging:"required"`
	Password string `json:"password" binging:"required"`
}

// @Summary      Sign Up
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param request body SingUpRequest true "account data"
// @Success      200  {object}  userService.User
// @Failure      404,500  {object}  ErrorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) singUp(ctx *gin.Context) {
	var request SingUpRequest

	if err := ctx.BindJSON(&request); err != nil {
		newHTTPErrorResponse(ctx, errs.NewError(errs.Validation, err))
		return
	}
	id, err := h.services.User.CreateUser(userService.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		newHTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Sign In
// @Description  Login and get access token
// @Tags         auth
// @Accept       */*
// @Produce      json
// @Param request body SingInRequest true "login data"
// @Success      200  {string}  token
// @Failure      404,500  {object}  ErrorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) singIn(ctx *gin.Context) {
	var request SingInRequest

	if err := ctx.BindJSON(&request); err != nil {
		newHTTPErrorResponse(ctx, errs.NewError(errs.Validation, err))
		return
	}
	id, err := h.services.Authorization.GenerateToken(request.Email, request.Password)
	if err != nil {
		newHTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"token": id,
	})
}
