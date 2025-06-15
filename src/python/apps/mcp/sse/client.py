import asyncio
import logging
from contextlib import asynccontextmanager, AsyncExitStack
from typing import AsyncGenerator, Any

from mcp import ClientSession
from mcp.client.sse import sse_client

from apps.mcp.openai_client import TransportStrategy, main

logger = logging.getLogger(__name__)


class SSETransportStrategy(TransportStrategy):
    def __init__(
            self,
            url: str,
            headers: dict[str, Any] | None = None,
            timeout: float = 5,
            sse_read_timeout: float = 60 * 5,
    ):
        self.url = url
        self.headers = headers
        self.timeout = timeout
        self.sse_read_timeout = sse_read_timeout

    @asynccontextmanager
    async def create_session(self) -> AsyncGenerator[ClientSession, None]:
        exit_stack = AsyncExitStack()
        try:
            stdio_transport = await exit_stack.enter_async_context(
                sse_client(
                    self.url,
                    self.headers,
                    self.timeout,
                    self.sse_read_timeout,
                )
            )
            read, write = stdio_transport

            session = await exit_stack.enter_async_context(
                ClientSession(read, write)
            )
            await session.initialize()

            yield session
        except Exception as e:
            logger.error(f"Transport failed: {str(e)}")
            raise
        finally:
            await exit_stack.aclose()


if __name__ == '__main__':
    transport_strategy = SSETransportStrategy(
        url="http://localhost:8000/sse",
    )
    asyncio.run(main(transport_strategy))
