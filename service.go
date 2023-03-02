package wxp

import "github.com/gin-gonic/gin"

type Service[T any] func(ctx *gin.Context, req *T) (interface{}, Error)
