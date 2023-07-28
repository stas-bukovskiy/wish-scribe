package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"net/http"
	"strconv"
)

func (h *Handler) getById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newHTTPErrorResponse(ctx, errs.NewError(errs.InvalidRequest, "invalid id parameter"))
		return
	}

	user, err := h.services.User.GetById(uint(userId))
	if err != nil {
		newHTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
