package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"net/http"
)

func (h *Handler) verifyToken(ctx *gin.Context) {
	token := ctx.Param("token")
	user, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		if errs.KindIs(errs.Unauthorized, err) {
			newHTTPErrorResponse(ctx, errs.NewError(errs.NotExist, "Invalid access token"))
		}
		newHTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
