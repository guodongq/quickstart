import logging
from abc import ABC, abstractmethod
from typing import Optional

from pydantic import BaseModel as PydanticBaseModel, Field

logger = logging.getLogger(__name__)


class Model(ABC):
    """Abstract base class for models."""

    @property
    @abstractmethod
    def id(self) -> str:
        pass


class BaseModel(PydanticBaseModel, Model):
    """Base model with Pydantic v2 configuration."""
    id: Optional[str] = Field(default=None, alias="_id")

    model_config = {
        "populate_by_name": True,
        "use_enum_values": True,
        "arbitrary_types_allowed": True
    }
