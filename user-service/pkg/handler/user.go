package handler

import "github.com/gin-gonic/gin"

func (h *Handler) getById(ctx *gin.Context) {
	ctx.AbortWithStatus(200)
}
