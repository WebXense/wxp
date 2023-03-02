package wxp

import (
	"github.com/WebXense/sql"
	"gorm.io/gorm"
)

type BaseRepository[T any] struct{}

func (r *BaseRepository[T]) Create(db *gorm.DB, entity *T) (*T, error) {
	return sql.Create(db, entity)
}

func (r *BaseRepository[T]) Update(db *gorm.DB, entity *T) (*T, error) {
	return sql.Update(db, entity)
}

func (r *BaseRepository[T]) Delete(db *gorm.DB, entity *T) error {
	return sql.Delete(db, entity)
}

func (r *BaseRepository[T]) FindAll(db *gorm.DB, page *sql.Pagination, sort *sql.Sort) ([]T, error) {
	return sql.FindAll[T](db, nil, page, sort)
}

func (r *BaseRepository[T]) FindByID(db *gorm.DB, id uint) (*T, error) {
	return sql.FindOne[T](db, sql.Eq("id", id))
}
