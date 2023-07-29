package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/errs"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

func (h *Handler) userIndemnity(c *gin.Context) {
	log := h.logger.Named("userIndemnity")
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		errs.NewHTTPErrorResponse(c, log, errs.NewError(errs.Unauthorized, "Missing authorization header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		errs.NewHTTPErrorResponse(c, log, errs.NewError(errs.Unauthorized, "Invalid authorization header"))
		return
	}

	user, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		errs.NewHTTPErrorResponse(c, log, err)
		return
	}

	c.Set(userCtx, user)
}
