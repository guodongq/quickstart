import asyncio
import json
from typing import Optional

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from contextlib import AsyncExitStack
from openai import OpenAI

import os
import urllib3

os.environ.pop("http_proxy", None)
os.environ.pop("https_proxy", None)
# Suppress warnings about Elasticsearch certificates
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class MCPClient:
    def __init__(self):
        """初始化 MCP 客户端"""
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()
        self.model = "llama3.2:latest"
        self.client = OpenAI(
            api_key="ollama",  # 读取OpenAI API Key
            base_url="http://localhost:11434/v1/",  # 使用local ollama
            timeout=3000,
        )

    async def connect_to_server(self, server_script_path: str):
        """Connect to an MCP server

        Args:
            server_script_path: Path to the server script (.py or .js)
        """
        is_python = server_script_path.endswith('.py')
        is_js = server_script_path.endswith('.js')
        if not (is_python or is_js):
            raise ValueError("Server script must be a .py or .js file")

        command = "python" if is_python else "node"
        server_params = StdioServerParameters(
            command=command,
            args=[server_script_path],
            env=None
        )

        stdio_transport = await self.exit_stack.enter_async_context(stdio_client(server_params))
        self.stdio, self.write = stdio_transport
        self.session = await self.exit_stack.enter_async_context(ClientSession(self.stdio, self.write))

        await self.session.initialize()

        # List available tools
        response = await self.session.list_tools()
        tools = response.tools
        print("\nConnected to server with tools:", [tool.name for tool in tools])

    async def process_query(self, query: str) -> str:
        """Process a query using LLM and available tools"""
        system_prompt = (
            "You are a helpful assistant."
            "You have the function of online search. "
            "Please MUST call web_search tool to search the Internet content before answering."
            "Please do not lose the user's question information when searching,"
            "and try to maintain the completeness of the question content as much as possible."
            "When there is a date related question in the user's question,"
            "please use the search function directly to search and PROHIBIT inserting specific time."
        )
        messages = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": query}
        ]

        # List available tools
        response = await self.session.list_tools()
        available_tools = [{
            "type": "function",
            "function": {
                "name": tool.name,
                "description": tool.description,
                "input_schema": tool.inputSchema
            }
        } for tool in response.tools]

        # Initial model response
        response = self.client.chat.completions.create(
            model=self.model,
            messages=messages,
            tools=available_tools,
        )

        content = response.choices[0]
        if content.finish_reason == "tool_calls":
            # Parse tool call
            tool_call = content.message.tool_calls[0]
            tool_name = tool_call.function.name
            tool_args = json.loads(tool_call.function.arguments)

            # Execute tool
            result = await self.session.call_tool(tool_name, tool_args)
            print(f"\n\n[Calling tool {tool_name} with args {tool_args}]\n\n")
            print(f"[Tool result: {result.content}]")

            # Append tool result to messages
            messages.append(content.message.model_dump())
            messages.append({
                "role": "tool",
                "content": result.content[0].text,
                "tool_call_id": tool_call.id,
            })

            # Generate final response
            response = self.client.chat.completions.create(
                model=self.model,
                messages=messages,
            )
            return response.choices[0].message.content

        return content.message.content

    async def chat_loop(self):
        """运行交互式聊天循环"""
        print("\n🤖 MCP 客户端已启动！输入 'quit' 退出")

        while True:
            try:
                query = input("\n你: ").strip()
                if query.lower() == 'quit':
                    break

                response = await self.process_query(query)  # 发送用户输入到 OpenAI API
                print(f"\n🤖 OpenAI: {response}")

            except Exception as e:
                print(f"\n⚠️ 发生错误: {str(e)}")
                import traceback
                traceback.print_exc()

    async def cleanup(self):
        """清理资源"""
        await self.exit_stack.aclose()


async def main():
    client = MCPClient()
    try:
        await client.connect_to_server(
            "/Users/guodongq/Workspaces/src/github.com/guodongq/quickstart/src/python/app/mcp_quickstart/stdio/server.py",
        )
        await client.chat_loop()
    finally:
        await client.cleanup()


if __name__ == '__main__':
    asyncio.run(main())
