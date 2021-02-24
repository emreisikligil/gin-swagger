package api

import "github.com/gin-gonic/gin"

type APIResponse interface {
	WriteResponse(ctx *gin.Context)
}
