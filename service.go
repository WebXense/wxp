package wxp

import (
	"github.com/WebXense/sql"
	"github.com/gin-gonic/gin"
)

type Request[T any] struct {
	Ctx    *gin.Context
	Object *T
	Page   *sql.Pagination
	Sort   *sql.Sort
}

type Response[T any] struct {
	Data *T
}

type Service[T any] func(req *Request[T]) (interface{}, Error)
