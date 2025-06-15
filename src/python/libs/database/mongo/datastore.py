from typing import TypeVar, Generic, List, Optional, Dict, Any, ClassVar, Type

from pymongo.collection import Collection
from pymongo.errors import PyMongoError, DuplicateKeyError

from libs.database.mongo.error import MongoError, NotFoundError
from libs.database.mongo.models import Model
from libs.database.mongo.repository import MongoRepository

ModelType = TypeVar("ModelType", bound=Model)


class DataStore(Generic[ModelType]):
    def __init__(
            self,
            mongo_repository: MongoRepository,
            collection_name: str,
            model_type: Type[ModelType]  # 显式传递模型类型
    ):
        self.mongo_repository = mongo_repository
        self.collection_name = collection_name
        self.model_type = model_type

    @property
    def collection(self) -> Collection:
        return self.mongo_repository.database()[self.collection_name]

    def save(self, model: ModelType) -> ModelType:
        try:
            # 排除None值减少存储空间
            doc = model.model_dump(by_alias=True, exclude_none=True)
            result = self.collection.insert_one(doc)
            if not model.id:
                model.id = str(result.inserted_id)
            return model
        except DuplicateKeyError as e:
            raise MongoError("Duplicate key error", e) from e
        except PyMongoError as e:
            raise MongoError("Failed to save model", e) from e

    def save_many(self, models: List[ModelType]) -> List[str]:
        try:
            docs = [m.model_dump(by_alias=True, exclude_none=True) for m in models]
            result = self.collection.insert_many(docs)
            inserted_ids = [str(id) for id in result.inserted_ids]
            # 更新模型中的ID
            for model, inserted_id in zip(models, inserted_ids):
                if not model.id:
                    model.id = inserted_id
            return inserted_ids
        except PyMongoError as e:
            raise MongoError("Failed to save multiple models", e) from e

    def get(self, id: str) -> ModelType:
        if result := self.collection.find_one({"_id": id}):
            return self._to_model(self.model_type, result)
        raise NotFoundError(f"Document with id {id} not found")

    def update(self, model: ModelType, upsert: bool = False) -> int:
        try:
            result = self.collection.update_one(
                {"_id": model.id},  # 直接使用id属性
                {"$set": model.model_dump(by_alias=True, exclude_none=True)},
                upsert=upsert
            )
            if not upsert and result.matched_count == 0:
                raise NotFoundError()
            return result.modified_count
        except PyMongoError as e:
            raise MongoError("Failed to update model", e) from e

    def delete(self, id: str) -> int:
        result = self.collection.delete_one({"_id": id})
        if result.deleted_count == 0:
            raise NotFoundError()
        return result.deleted_count

    def find(
            self,
            filter_query: Optional[Dict[str, Any]] = None,
            limit: int = 0,
            skip: int = 0,
            sort: Optional[List[tuple]] = None,
            **find_options
    ) -> List[ModelType]:
        try:
            cursor = self.collection.find(filter_query or {}, **find_options)
            if sort:
                cursor = cursor.sort(sort)
            if skip:
                cursor = cursor.skip(skip)
            if limit:
                cursor = cursor.limit(limit)
            return [self._to_model(self.model_type, doc) for doc in cursor]
        except PyMongoError as e:
            raise MongoError("Failed to find models", e) from e

    def count(self, filter_query: Optional[Dict[str, Any]] = None) -> int:
        try:
            return self.collection.count_documents(filter_query or {})
        except PyMongoError as e:
            raise MongoError("Failed to count documents", e) from e

    @staticmethod
    def _to_model(model_type: Type[ModelType], doc: Dict[str, Any]) -> ModelType:
        """将MongoDB文档转换为模型实例"""
        # 确保文档中的_id映射到模型的id字段
        if "_id" in doc:
            doc["id"] = str(doc["_id"])
        return model_type.parse_obj(doc)
