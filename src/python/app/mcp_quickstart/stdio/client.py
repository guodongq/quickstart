import asyncio
import logging
from contextlib import AsyncExitStack, asynccontextmanager
from typing import Dict, AsyncGenerator

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

from app.mcp_quickstart import TransportStrategy, main

logger = logging.getLogger(__name__)


class StdioTransportStrategy(TransportStrategy):
    def __init__(
            self,
            command: str = "python",
            server_script_path: str = None,
            env: Dict[str, str] | None = None
    ):
        self.command = command
        self.server_script_path = server_script_path
        self.env = env

    @asynccontextmanager
    async def create_session(self) -> AsyncGenerator[ClientSession, None]:
        """Create a new session with the server."""
        server_params = StdioServerParameters(
            command=self.command,
            args=[self.server_script_path],
            env=self.env,
        )

        exit_stack = AsyncExitStack()
        try:
            stdio_transport = await exit_stack.enter_async_context(
                stdio_client(server_params)
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
    transport_strategy = StdioTransportStrategy(
        server_script_path="/Users/guodongq/Workspaces/src/github.com/guodongq/quickstart/src/python/app/mcp_quickstart/stdio/server.py"
    )
    asyncio.run(main(transport_strategy))
