package mysql

import "database/sql"

type GenericRepository struct {
	Conn *sql.DB
}

func NewGenericRepository(conn *sql.DB) *GenericRepository {
	return &GenericRepository{
		Conn: conn,
	}
}

//
//
//import (
//"context"
//"fmt"
//"github.com/guodongq/abc/pkg/library/core/types"
//"github.com/guodongq/abc/pkg/library/modules/seedwork/domain"
//)
//import "gorm.io/gorm"
//
//type GormRepository[T domain.Entity, P domain.PersistentEntity] struct {
//	db *gorm.DB
//}
//
//func NewGormRepository[T domain.Entity, P domain.PersistentEntity](db *gorm.DB) *GormRepository[T, P] {
//	return &GormRepository[T, P]{db: db}
//}
//
//func (r *GormRepository[T, P]) Save(ctx context.Context, entity T) error {
//	var persistentEntity P
//	persistentEntity.FromDomain(entity)
//	return r.db.WithContext(ctx).Create(&persistentEntity).Error
//}
//
//func (r *GormRepository[T, P]) Update(ctx context.Context, id string, updateFn func(ctx context.Context, entity T) (T, error)) error {
//	domainEntity, err := r.Get(ctx, id)
//	if err != nil {
//		return err
//	}
//
//	var persistentEntity P
//	persistentEntity.FromDomain(domainEntity)
//	return r.db.WithContext(ctx).Save(&persistentEntity).Error
//}
//
//func (r *GormRepository[T, P]) Delete(ctx context.Context, id string) error {
//	return r.db.WithContext(ctx).Delete(&id).Error
//}
//
//func (r *GormRepository[T, P]) Get(ctx context.Context, id string) (T, error) {
//	var persistentEntity P
//	err := r.db.WithContext(ctx).First(&persistentEntity, "id = ?", id).Error
//	if err != nil {
//		var zero T
//		return zero, err
//	}
//	return persistentEntity.ToDomain().(T), nil
//}
//
//func (r *GormRepository[T, P]) GetAll(ctx context.Context, params types.QueryOptions) (domain.PagedResult[T], error) {
//	var entities []T
//	query := r.db.WithContext(ctx)
//
//	// Apply filters
//	for key, value := range params.Filter.ToFilter() {
//		query = query.Where(fmt.Sprintf("%s = ?", key), value)
//	}
//
//	// Apply sorting
//	for _, sortable := range params.Sortable.ComputedSort() {
//		query = query.Order(fmt.Sprintf("%s %s", sortable.OrderBy, sortable.Direction))
//	}
//
//	// Pagination
//	var totalCount int64
//	query.Model(new(T)).Count(&totalCount)
//	query = query.Offset(int(params.Skip)).Limit(int(params.Limit))
//	err := query.Find(&entities).Error
//
//	return domain.PagedResult[T]{
//		Items:      entities,
//		TotalCount: totalCount,
//		Pageable:   params.Pageable,
//	}, err
//}
//
//func (r *GormRepository[T, P]) Count(ctx context.Context, params types.QueryOptions) (int64, error) {
//	var count int64
//	query := r.db.WithContext(ctx)
//
//	// Apply filters
//	for key, value := range params.Filter.ToFilter() {
//		query = query.Where(fmt.Sprintf("%s = ?", key), value)
//	}
//
//	err := query.Model(new(T)).Count(&count).Error
//	return count, err
//}
