from mcp.server.fastmcp import FastMCP
import datetime

mcp = FastMCP("time")


@mcp.tool()
def get_current_time() -> str:
    """Get current time"""
    current_time = datetime.datetime.now(
        datetime.timezone(datetime.timedelta(hours=8))
    )  # CST时间，假设东八区
    formatted_time = current_time.strftime("%I:%M:%S %p %Z")
    return formatted_time


if __name__ == "__main__":
    mcp.run(transport="sse")