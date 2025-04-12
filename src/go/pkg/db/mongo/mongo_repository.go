package mongo

//
//import (
//	"context"
//	"errors"
//	"github.com/guodongq/abc/pkg/library/core/types"
//	"github.com/guodongq/abc/pkg/library/modules/provider/mongodb"
//	"github.com/guodongq/abc/pkg/library/modules/seedwork/domain"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"reflect"
//)
//
//type MongoRepository[T domain.Entity, P domain.PersistentEntity] struct {
//	mongoClient mongodb.MongoRepository
//}
//
//func NewMongoRepository[T domain.Entity, P domain.PersistentEntity](
//	mongoClient mongodb.MongoRepository,
//) MongoRepository[T, P] {
//	return MongoRepository[T, P]{
//		mongoClient: mongoClient,
//	}
//}
//
//func (r MongoRepository[T, P]) collection(ctx context.Context, entity domain.Entity) *mongo.Collection {
//	var tableName string
//
//	// 检查 T 是否实现了 TableName 方法
//	tableNameMethod := reflect.ValueOf(entity).MethodByName("TableName")
//	if tableNameMethod.IsValid() {
//		// 调用 TableName 方法获取表名
//		results := tableNameMethod.Call(nil)
//		if len(results) > 0 && results[0].Kind() == reflect.String {
//			tableName = results[0].String()
//		}
//	}
//
//	// 如果没有实现 TableName 方法，使用默认逻辑
//	if tableName == "" {
//		tableName = reflect.TypeOf(entity).Elem().Name()
//	}
//	return r.mongoClient.Database(ctx).Collection(tableName)
//}
//
//func (r MongoRepository[T, P]) Save(ctx context.Context, entity T) error {
//	collection := r.collection(ctx, entity)
//
//	// 将领域对象转换为持久化对象
//	var persistentEntity P
//	persistentEntity.FromDomain(entity)
//
//	// 插入到 MongoDB
//	_, err := collection.InsertOne(ctx, persistentEntity)
//	return err
//}
//
//func (r MongoRepository[T, P]) Get(ctx context.Context, id string) (T, error) {
//	collection := r.collection(ctx, *new(T))
//
//	// 查询 MongoDB
//	var persistentEntity P
//	err := collection.FindOne(ctx, types.Empty().SetID(id)).Decode(&persistentEntity)
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			var zero T
//			return zero, errors.New("entity not found")
//		}
//		var zero T
//		return zero, err
//	}
//
//	// 将持久化对象转换为领域对象
//	return persistentEntity.ToDomain().(T), nil
//}
//
//func (r MongoRepository[T, P]) Update(ctx context.Context, id string, updateFn func(ctx context.Context, entity T) (T, error)) error {
//	return r.mongoClient.RunTransaction(ctx, func(sctx mongo.SessionContext) error {
//		collection := r.collection(ctx, *new(T))
//
//		// 获取当前实体
//		var persistentEntity P
//		err := collection.FindOne(ctx, types.Empty().SetID(id)).Decode(&persistentEntity)
//		if err != nil {
//			if errors.Is(err, mongo.ErrNoDocuments) {
//				return errors.New("entity not found")
//			}
//			return err
//		}
//
//		// 将持久化对象转换为领域对象
//		domainEntity := persistentEntity.ToDomain().(T)
//
//		// 调用更新函数
//		updatedEntity, err := updateFn(ctx, domainEntity)
//		if err != nil {
//			return err
//		}
//
//		// 将更新后的领域对象转换为持久化对象
//		var updatedPersistentEntity P
//		updatedPersistentEntity.FromDomain(updatedEntity)
//
//		// 更新 MongoDB
//		_, err = collection.ReplaceOne(ctx, types.Empty().SetID(id), updatedPersistentEntity)
//		return err
//	})
//}
//
//func (r MongoRepository[T, P]) Delete(ctx context.Context, id string) error {
//	collection := r.collection(ctx, *new(T))
//
//	// 删除 MongoDB
//	_, err := collection.DeleteOne(ctx, types.Empty().SetID(id))
//	return err
//}
//
//func (r MongoRepository[T, P]) GetAll(ctx context.Context, params types.QueryOptions) (domain.PagedResult[T], error) {
//	collection := r.collection(ctx, *new(T))
//
//	// 构建查询过滤器
//	filter := types.F{}
//	if params.Filter != nil {
//		filter = params.Filter.ToFilter()
//	}
//
//	// 查询总数
//	totalCount, err := collection.CountDocuments(ctx, filter)
//	if err != nil {
//		return domain.PagedResult[T]{}, err
//	}
//
//	// 分页选项
//	findOptions := options.Find().
//		SetSkip(params.Skip).
//		SetLimit(params.Limit)
//
//	// 排序选项
//	if params.Sort != nil {
//		sortFields := bson.D{}
//		for _, field := range params.ComputedSort() {
//			order := 1 // 升序
//			if types.DESC == field.Direction {
//				order = -1 // 降序
//			}
//			sortFields = append(sortFields, bson.E{Key: field.OrderBy, Value: order})
//		}
//		findOptions.SetSort(sortFields)
//	}
//
//	// 查询 MongoDB
//	cursor, err := collection.Find(ctx, filter, findOptions)
//	if err != nil {
//		return domain.PagedResult[T]{}, err
//	}
//	defer cursor.Close(ctx)
//
//	// 解析结果
//	var items []T
//	for cursor.Next(ctx) {
//		var persistentEntity P
//		if err := cursor.Decode(&persistentEntity); err != nil {
//			return domain.PagedResult[T]{}, err
//		}
//		items = append(items, persistentEntity.ToDomain().(T))
//	}
//
//	// 返回分页结果
//	return domain.PagedResult[T]{
//		Items:      items,
//		TotalCount: totalCount,
//		Pageable:   params.Pageable,
//	}, nil
//}
//
//// Count 获取实体总数
//func (r MongoRepository[T, P]) Count(ctx context.Context, params types.QueryOptions) (int64, error) {
//	collection := r.collection(ctx, *new(T))
//
//	// 构建查询过滤器
//	filter := types.F{}
//	if params.Filter != nil {
//		filter = params.Filter.ToFilter()
//	}
//
//	// 查询总数
//	return collection.CountDocuments(ctx, filter)
//}
