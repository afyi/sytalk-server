package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ctx1 context.Context = context.TODO()
)

// 释放
func close(ctx *gin.Context) {
	if err := recover(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    0,
			"message": err,
		})
	}
}
