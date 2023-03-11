package wxp

import (
	"github.com/WebXense/sql"
	"github.com/WebXense/sql/stm"
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

func (r *BaseRepository[T]) DeleteBy(db *gorm.DB, where *stm.Statement) error {
	return sql.DeleteBy[T](db, where)
}

func (r *BaseRepository[T]) FindOne(db *gorm.DB, where *stm.Statement) (*T, error) {
	return sql.FindOne[T](db, where)
}

func (r *BaseRepository[T]) FindAll(db *gorm.DB, page *sql.Pagination, sort *sql.Sort, where *stm.Statement) ([]T, error) {
	return sql.FindAll[T](db, where, page, sort)
}

func (r *BaseRepository[T]) Count(db *gorm.DB, where *stm.Statement) (int64, error) {
	return sql.Count[T](db, where)
}

func (r *BaseRepository[T]) FindByID(db *gorm.DB, id uint) (*T, error) {
	return r.FindOne(db, sql.Eq("id", id))
}
