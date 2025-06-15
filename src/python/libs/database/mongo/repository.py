from abc import ABC, abstractmethod
from contextlib import contextmanager
from typing import ClassVar, Dict

from pymongo.database import Database
from pymongo import MongoClient
from pymongo.errors import PyMongoError

from libs.database.mongo.error import MongoError


class MongoRepository(ABC):
    @abstractmethod
    def database(self) -> Database:
        pass

    @abstractmethod
    def close(self) -> None:
        pass


class PyMongoRepository(MongoRepository):
    _clients: ClassVar[Dict[str, MongoClient]] = {}

    def __init__(self, client: MongoClient, database_name: str):
        self._client = client
        self._database = client[database_name]

    def database(self) -> Database:
        """Get the MongoDB database instance."""
        return self._database

    def close(self) -> None:
        """Close the MongoDB client and remove from cache."""
        if self._client:
            client_key = f"{self._client.address}:{self._database.name}"
            if client_key in self._clients:
                del self._clients[client_key]
            self._client.close()

    @contextmanager
    def session(self, existing_session=None):
        """
        Provide a transactional session context.

        Args:
            existing_session: An existing MongoDB session to use, if any.

        Yields:
            A MongoDB session.

        Raises:
            MongoError: If the transaction fails.
        """
        if existing_session:
            with existing_session.start_transaction():
                yield existing_session
        else:
            with self.database().client.start_session() as session:
                try:
                    with session.start_transaction():
                        yield session
                except PyMongoError as e:
                    session.abort_transaction()
                    raise MongoError("Transaction failed", e) from e

    @classmethod
    def create_with_uri(
            cls,
            uri: str,
            database_name: str,
            **kwargs
    ) -> "PyMongoRepository":
        key = f"{uri}:{database_name}"
        if key not in cls._clients:
            cls._clients[key] = MongoClient(uri, **kwargs)
        return cls(cls._clients[key], database_name)
