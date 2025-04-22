from typing import Optional, Any

from mcp import ClientSession, StdioServerParameters
from mcp.client.sse import sse_client
from contextlib import AsyncExitStack


class MCPClient:
    def __init__(self):
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()

    async def connect_to_server(
            self,
            url: str,
            headers: dict[str, Any] | None = None,
            timeout: float = 5,
            sse_read_timeout: float = 60 * 5,
    ):
        sse_transport = await self.exit_stack.enter_async_context(
            sse_client(url, headers, timeout, sse_read_timeout)
        )
        read, write = sse_transport
        self.session = await self.exit_stack.enter_async_context(
            ClientSession(read, write)
        )

        await self.session.initialize()

        response = await self.session.list_tools()
        print("Connected to server with tools:", [tool.name for tool in response.tools])
