import asyncio
import logging
from contextlib import AsyncExitStack, asynccontextmanager
from typing import AsyncGenerator

from mcp import ClientSession
from mcp.client.websocket import websocket_client

from apps.mcp.openai_client import TransportStrategy, main

logger = logging.getLogger(__name__)


class WebSocketTransportStrategy(TransportStrategy):
    def __init__(self, url: str = "ws://localhost:8000/mcp"):
        self.url = url

    @asynccontextmanager
    async def create_session(self) -> AsyncGenerator[ClientSession, None]:
        """Create a new session with the server."""
        exit_stack = AsyncExitStack()
        try:
            websocket_transport = await exit_stack.enter_async_context(
                websocket_client(self.url)
            )
            read, write = websocket_transport

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
    transport_strategy = WebSocketTransportStrategy()
    asyncio.run(main(transport_strategy))
