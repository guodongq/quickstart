import uuid
from typing import Any
from pydantic import GetCoreSchemaHandler


class GenericUUID(uuid.UUID):
    @classmethod
    def next_id(cls):
        return cls(int=uuid.uuid4().int)

    @classmethod
    def __get_pydantic_core_schema__(
            cls,
            source_type: Any,
            handler: GetCoreSchemaHandler,
    ):
        return handler.generate_schema(uuid.UUID)

class ValueObject:
    """
    Base class for value objects
    """