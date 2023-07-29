package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	"net/http"
)

// @Summary      Verify Token
// @Description  Verify access token and get user info
// @Tags         token
// @Accept       */*
// @Produce      json
// @Param token path string true "token to verify"
// @Success      200  {object}  user_service.User
// @Failure      404,500  {object}  ErrorResponse
// @Router       /api/v1/tokens/{token} [get]
func (h *Handler) verifyToken(ctx *gin.Context) {
	log := h.logger.Named("verifyToken")
	token := ctx.Param("token")
	user, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		if errs.KindIs(errs.Unauthorized, err) {
			errs.NewHTTPErrorResponse(ctx, log, errs.NewError(errs.BadRequest, "Invalid access token"))
			return
		}
		errs.NewHTTPErrorResponse(ctx, log, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
