package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/user-service/pkg/errs"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

func (h *Handler) userIndemnity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newHTTPErrorResponse(c, errs.NewError(errs.Unauthorized, "Missing authorization header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newHTTPErrorResponse(c, errs.NewError(errs.Unauthorized, "Invalid authorization header"))
		return
	}

	user, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newHTTPErrorResponse(c, err)
		return
	}

	c.Set(userCtx, user)
}
