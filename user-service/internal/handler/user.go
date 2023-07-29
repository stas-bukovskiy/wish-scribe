package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	"net/http"
	"strconv"
)

// @Summary      Get By ID
// @Description  Get user info by its id
// @Security 	 ApiKeyAuth
// @Tags         users
// @Accept       */*
// @Produce      json
// @Param        id path int true "user id"
// @Success      200  {object}  user_service.User
// @Failure      404,500  {object}  ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (h *Handler) getById(ctx *gin.Context) {
	log := h.logger.Named("getById")
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errs.NewHTTPErrorResponse(ctx, log, errs.NewError(errs.BadRequest, "invalid id parameter"))
		return
	}

	user, err := h.services.User.GetById(uint(userId))
	if err != nil {
		errs.NewHTTPErrorResponse(ctx, log, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
