package database

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/guodongq/quickstart/pkg/types"
//	"github.com/guodongq/quickstart/pkg/ddd"
//	"reflect"
//	"sort"
//	"sync"
//)
//
//type InMemoryRepository[T ddd.Entity, P ddd.PersistentEntity] struct {
//	data sync.Map
//	mu   sync.RWMutex
//}
//
//func NewInMemoryRepository[T domain.Entity, P domain.PersistentEntity]() *InMemoryRepository[T, P] {
//	return &InMemoryRepository[T, P]{}
//}
//
//func (r *InMemoryRepository[T, P]) Save(ctx context.Context, entity T) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	var persistentEntity P
//	persistentEntity.FromDomain(entity)
//
//	r.data.Store(entity.ID(), persistentEntity)
//	return nil
//}
//
//func (r *InMemoryRepository[T, P]) Get(ctx context.Context, id string) (T, error) {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	value, exists := r.data.Load(id)
//	if !exists {
//		var zero T
//		return zero, errors.New("entity not found")
//	}
//
//	persistentEntity := value.(P)
//	return persistentEntity.ToDomain().(T), nil
//
//}
//
//func (r *InMemoryRepository[T, P]) Update(ctx context.Context, id string, updateFn func(ctx context.Context, entity T) (T, error)) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	value, exists := r.data.Load(id)
//	if !exists {
//		return fmt.Errorf("entity with D: %v not found", id)
//	}
//
//	persistentEntity := value.(P)
//	domainEntity := persistentEntity.ToDomain().(T)
//
//	updateEntity, err := updateFn(ctx, domainEntity)
//	if err != nil {
//		return err
//	}
//
//	var updatedPersistentEntity P
//	updatedPersistentEntity.FromDomain(updateEntity)
//
//	r.data.Store(id, updatedPersistentEntity)
//	return nil
//}
//
//func (r *InMemoryRepository[T, P]) Delete(ctx context.Context, id string) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	if _, exists := r.data.Load(id); !exists {
//		return fmt.Errorf("entity with D: %v not found", id)
//	}
//	r.data.Delete(id)
//	return nil
//}
//
//func (r *InMemoryRepository[T, P]) GetAll(ctx context.Context, params types.QueryOptions) (domain.PagedResult[T], error) {
//	filters := params.Filter.ToFilter()
//	var filteredItems []T
//
//	// 加读锁
//	r.mu.RLock()
//	defer r.mu.RUnlock()
//
//	// 遍历 sync.Map，过滤数据
//	r.data.Range(func(key, value any) bool {
//		entity, ok := value.(T)
//		if !ok {
//			return true // 跳过类型不匹配的项
//		}
//
//		// 使用反射检查字段值
//		val := reflect.ValueOf(entity).Elem()
//		matchesAllFilters := true
//
//		for filterKey, filterValue := range filters {
//			fieldVal := val.FieldByName(filterKey)
//			if !fieldVal.IsValid() || fieldVal.Interface() != filterValue {
//				matchesAllFilters = false
//				break
//			}
//		}
//
//		if matchesAllFilters {
//			filteredItems = append(filteredItems, entity)
//		}
//		return true
//	})
//
//	// 排序逻辑
//	if params.Sort != nil {
//		sort.Slice(filteredItems, func(i, j int) bool {
//			valI := reflect.ValueOf(filteredItems[i]).Elem()
//			valJ := reflect.ValueOf(filteredItems[j]).Elem()
//
//			for _, sortField := range params.ComputedSort() {
//				fieldI := valI.FieldByName(sortField.OrderBy)
//				fieldJ := valJ.FieldByName(sortField.OrderBy)
//
//				if !fieldI.IsValid() || !fieldJ.IsValid() {
//					continue
//				}
//
//				// 比较字段值
//				switch fieldI.Kind() {
//				case reflect.String:
//					if fieldI.String() != fieldJ.String() {
//						if sortField.Direction == types.ASC {
//							return fieldI.String() < fieldJ.String()
//						} else {
//							return fieldI.String() > fieldJ.String()
//						}
//					}
//				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//					if fieldI.Int() != fieldJ.Int() {
//						if sortField.Direction == types.ASC {
//							return fieldI.Int() < fieldJ.Int()
//						} else {
//							return fieldI.Int() > fieldJ.Int()
//						}
//					}
//				case reflect.Float32, reflect.Float64:
//					if fieldI.Float() != fieldJ.Float() {
//						if sortField.Direction == types.ASC {
//							return fieldI.Float() < fieldJ.Float()
//						} else {
//							return fieldI.Float() > fieldJ.Float()
//						}
//					}
//				default:
//					// 不支持的类型，跳过
//					continue
//				}
//			}
//			return false
//		})
//	}
//
//	// 分页逻辑
//	start := params.Skip
//	end := start + params.Limit
//	totalCount := int64(len(filteredItems))
//	if end > totalCount {
//		end = totalCount
//	}
//
//	// 获取分页后的数据
//	pagedItems := filteredItems[start:end]
//
//	// 返回分页结果
//	return domain.PagedResult[T]{
//		Items:      pagedItems,
//		TotalCount: totalCount,
//		Pageable:   params.Pageable,
//	}, nil
//}
//
//func (r *InMemoryRepository[T, P]) Count(ctx context.Context, params types.QueryOptions) (int64, error) {
//	filters := params.Filter.ToFilter()
//	if len(filters) == 0 {
//		var count int64
//		r.data.Range(func(key, value any) bool {
//			count++
//			return true
//		})
//		return count, nil
//	}
//
//	var count int64
//	r.mu.RLock()
//	defer r.mu.RUnlock()
//
//	r.data.Range(func(key, value any) bool {
//		entity, ok := value.(T)
//		if !ok {
//			return true
//		}
//
//		val := reflect.ValueOf(entity).Elem()
//		matchesAllFilters := true
//
//		for filterKey, filterValue := range filters {
//			fieldVal := val.FieldByName(filterKey)
//			if !fieldVal.IsValid() && fieldVal.Interface() != filterValue {
//				matchesAllFilters = false
//				break
//			}
//		}
//
//		if matchesAllFilters {
//			count++
//		}
//
//		return true
//	})
//
//	return count, nil
//}
