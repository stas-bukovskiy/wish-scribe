package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
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
func (h *Handler) signUp(ctx *gin.Context) {
	log := h.logger.Named("signUp")
	var request SingUpRequest

	if err := ctx.BindJSON(&request); err != nil {
		errs.NewHTTPErrorResponse(ctx, log, errs.NewError(errs.Validation, err))
		return
	}
	id, err := h.services.User.CreateUser(userService.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		errs.NewHTTPErrorResponse(ctx, log, err)
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
func (h *Handler) signIn(ctx *gin.Context) {
	log := h.logger.Named("signIn")

	var request SingInRequest

	if err := ctx.BindJSON(&request); err != nil {
		errs.NewHTTPErrorResponse(ctx, log, errs.NewError(errs.Validation, err))
		return
	}
	id, err := h.services.Authorization.GenerateToken(request.Email, request.Password)
	if err != nil {
		errs.NewHTTPErrorResponse(ctx, log, err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"token": id,
	})
}
