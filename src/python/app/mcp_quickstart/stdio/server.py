import json
import httpx
from typing import Any
from mcp.server.fastmcp import FastMCP

import os
import urllib3

os.environ.pop("http_proxy", None)
os.environ.pop("https_proxy", None)
# Suppress warnings about Elasticsearch certificates
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Initialize FastMCP server
mcp = FastMCP("WeatherServer")

# Constants
OPENWEATHER_API_BASE = "https://api.openweathermap.org/data/2.5/weather"
API_KEY = "bd5e378503939ddaee76f12ad7a97608"  # 请替换为你自己的 OpenWeather API Key
USER_AGENT = "weather-app/1.0"


async def fetch_weather(city: str) -> dict[str, Any] | None:
    """Make a request to the OpenWeather API with proper error handling."""
    # weather = '{"coord": {"lon": 116.3972, "lat": 39.9075}, "weather": [{"id": 800, "main": "Clear", "description": "晴", "icon": "01d"}], "base": "stations", "main": {"temp": 23.94, "feels_like": 22.65, "temp_min": 23.94, "temp_max": 23.94, "pressure": 1005, "humidity": 10, "sea_level": 1005, "grnd_level": 1000}, "visibility": 10000, "wind": {"speed": 4.63, "deg": 343, "gust": 9.31}, "clouds": {"all": 6}, "dt": 1744612286, "sys": {"type": 1, "id": 9609, "country": "CN", "sunrise": 1744580304, "sunset": 1744627839}, "timezone": 28800, "id": 1816670, "name": "Beijing", "cod": 200}'
    # return json.loads(weather)
    params = {"q": city, "appid": API_KEY, "units": "metric", "lang": "zh_cn"}
    headers = {"User-Agent": USER_AGENT}

    async with httpx.AsyncClient(verify=False) as client:
        try:
            response = await client.get(
                OPENWEATHER_API_BASE, params=params, headers=headers, timeout=30.0
            )
            response.raise_for_status()
            return response.json()
        except httpx.HTTPStatusError as e:
            return {"error": f"HTTP error: {e.response.status_code}"}
        except Exception as e:
            return {"error": f"Request failed: {str(e)}"}


def format_weather(data: dict[str, Any] | str) -> str:
    """
    将天气数据格式化为易读文本。
    :param data: 天气数据（可以是字典或 JSON 字符串）
    :return: 格式化后的天气信息字符串
    """
    # 如果传入的是字符串，则先转换为字典
    if isinstance(data, str):
        try:
            data = json.loads(data)
        except Exception as e:
            return f"Failed to parse the weather data: {e}"

    # 如果数据中包含错误信息，直接返回错误提示
    if "error" in data:
        return f"⚠️ {data['error']}"

    # 提取数据时做容错处理
    city = data.get("name", "Unknown")
    country = data.get("sys", {}).get("country", "Unknown")
    temp = data.get("main", {}).get("temp", "N/A")
    humidity = data.get("main", {}).get("humidity", "N/A")
    wind_speed = data.get("wind", {}).get("speed", "N/A")
    # weather 可能为空列表，因此用 [0] 前先提供默认字典
    weather_list = data.get("weather", [{}])
    description = weather_list[0].get("description", "Unknown")

    return (
        f"🌍 {city}, {country}\n"
        f"🌡 Temperature: {temp}°C\n"
        f"💧 Humidity: {humidity}%\n"
        f"🌬 Wind Speed: {wind_speed} m/s\n"
        f"🌤 Weather: {description}\n"
    )


@mcp.tool()
async def query_weather(city: str) -> str:
    """
    Input the English name of a specified city and return the result of today's weather query.
    :param city: The name of the city (e.g, Beijing).
    :return: Formatted weather information.
    """
    data = await fetch_weather(city)
    return format_weather(data)


if __name__ == "__main__":
    # 以标准 I/O 方式运行 MCP 服务器
    mcp.run(transport="stdio")
