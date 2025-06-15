class MongoError(Exception):
    """Base exception for MongoDB-related errors."""

    def __init__(self, message: str, cause: Exception = None):
        self.message = message
        self.cause = cause
        super().__init__(f"{message}: {str(cause)}" if cause else message)


class NotFoundError(MongoError):
    """Raised when a document is not found in the database."""

    def __init__(self, message: str = "Document not found", cause: Exception = None):
        super().__init__(message, cause)
