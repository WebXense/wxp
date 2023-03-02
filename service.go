package wxp

import (
	"github.com/WebXense/sql"
	"github.com/gin-gonic/gin"
)

type Request[T any] struct {
	Object *T
	Page   *sql.Pagination
	Sort   *sql.Sort
}

type Service[T any] func(ctx *gin.Context, req *Request[T]) (interface{}, Error)
