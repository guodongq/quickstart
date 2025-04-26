import datetime
from zoneinfo import ZoneInfo

import socketio
from mcp.server.fastmcp.server import FastMCP
from starlette.applications import Starlette
from starlette.routing import Mount
import uvicorn
import anyio


class SocketIOMCP(FastMCP):
    def __init__(self, sio: socketio.AsyncServer, **kwargs):
        super().__init__(**kwargs)
        self.sio = sio

    async def run_socketio_async(self):
        @self.sio.on("connect")
        async def connect(sid):
            await self.sio.emit("connect", "connected", room=sid)

        @self.sio.on("disconnect")
        async def disconnect(sid):
            pass

        @self.sio.on("mcp_transport")
        async def mcp_transport(sid, data):
            pass

        @self.sio.on("mcp_register")
        async def mcp_register(sid, data):
            pass

        starlette_app = Starlette(
            debug=self.settings.debug,
            routes=[
                Mount("/socket.io/", app=socketio.ASGIApp(self.sio))
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


mcp = SocketIOMCP(
    name="time",
    sio=socketio.AsyncServer(async_mode="asgi", cors_allowed_origins="*"),
)


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
    anyio.run(mcp.run_socketio_async)
