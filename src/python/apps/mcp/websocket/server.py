import logging

import anyio
from starlette.applications import Starlette
from starlette.routing import WebSocketRoute
from starlette.websockets import WebSocket
from mcp.server.websocket import websocket_server
import datetime
from zoneinfo import ZoneInfo
import uvicorn
from mcp.server.fastmcp.server import FastMCP


class WebSocketMCP(FastMCP):
    async def run_ws_async(self):
        async def websocket_endpoint(websocket: WebSocket):
            # 1. Use the websocket_server context manager
            async with websocket_server(
                    websocket.scope, websocket.receive, websocket.send
            ) as (read_stream, write_stream):
                # 2. It yields streams connected to this specific WebSocket
                logging.info(f"Server: WebSocket client connected. Running server logic.")
                # 3. Pass streams to the server's run method
                await self._mcp_server.run(
                    read_stream,
                    write_stream,
                    self._mcp_server.create_initialization_options()
                )
            logging.info("Server: WebSocket client disconnected.")

        starlette_app = Starlette(
            debug=self.settings.debug,
            routes=[
                WebSocketRoute("/mcp", endpoint=websocket_endpoint)
            ]
        )

        config = uvicorn.Config(
            starlette_app,
            host=self.settings.host,
            port=self.settings.port,
            log_level=self.settings.log_level.lower(),
        )
        server = uvicorn.Server(config)
        await server.serve()


mcp = WebSocketMCP("time")


@mcp.tool(
    name="get_current_time",
    description="Get current time in specified tz",
)
def get_current_time(tz: str = "Asia/Shanghai") -> str:
    """Get current time in specified timezone"""
    try:
        timezone = ZoneInfo(tz)
    except Exception as e:
        return f"Error: {str(e)}"

    current_time = datetime.datetime.now(timezone)
    formatted_time = current_time.strftime("%I:%M:%S %p %Z")
    return formatted_time


if __name__ == "__main__":
    anyio.run(mcp.run_ws_async)
