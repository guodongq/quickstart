from mcp.server.fastmcp import FastMCP
import datetime
from zoneinfo import ZoneInfo

mcp = FastMCP("time")


@mcp.tool(
    name="get_current_time_in_timezone",
    description="Get current time in specified timezone",
)
def get_current_time_in_timezone(tz: str = "Asia/Shanghai") -> str:
    """Get current time in specified timezone"""
    try:
        timezone = ZoneInfo(tz)
    except Exception as e:
        return f"Error: {str(e)}"

    current_time = datetime.datetime.now(timezone)
    formatted_time = current_time.strftime("%I:%M:%S %p %Z")
    return formatted_time


if __name__ == "__main__":
    mcp.run(transport="stdio")
