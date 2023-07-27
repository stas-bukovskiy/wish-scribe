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

func (h *Handler) singUp(ctx *gin.Context) {
	var request SingUpRequest

	if err := ctx.BindJSON(&request); err != nil {
		newHTTPErrorResponse(ctx, errs.NewError(errs.Validation, err))
		return
	}
	id, err := h.services.Authorization.CreateUser(userService.User{
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

func (h *Handler) singIn(ctx *gin.Context) {

}
