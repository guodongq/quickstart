import asyncio
from typing import Optional, Any

from mcp import ClientSession
from mcp.client.sse import sse_client
import logging
from langchain_core.tools import BaseTool


class MCPClient:
    def __init__(self):
        self.session: Optional[ClientSession] = None

    async def connect_to_server(
            self,
            url: str,
            headers: dict[str, Any] | None = None,
            timeout: float = 5,
            sse_read_timeout: float = 60 * 5,
    ):
        pass


async def main():
    client = MCPClient()
    await client.connect_to_server(
        url="http://localhost:8000/sse",
    )


if __name__ == '__main__':
    asyncio.run(main())
