import asyncio
import json
import logging
from contextlib import AsyncExitStack
from typing import Optional, List, Dict, Any

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from utils.openai_util import OpenAIClient, config


class MCPClient:
    def __init__(self):
        """Initialize MCP Client with optional tools."""
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()
        self.client = OpenAIClient()
        self.available_tools: List[Dict[str, Any]] = []

    async def connect_to_server(self, server_script_path: str):
        """Connect to an MCP server and cache available tools.

        Args:
            server_script_path: Path to the server script (.py or .js)
        """
        is_python = server_script_path.endswith(".py")
        is_js = server_script_path.endswith(".js")
        if not (is_python or is_js):
            raise ValueError("Server script must be a .py or .js file")

        command = "python" if is_python else "node"
        server_params = StdioServerParameters(
            command=command, args=[server_script_path], env=None
        )

        try:
            stdio_transport = await self.exit_stack.enter_async_context(
                stdio_client(server_params)
            )
            read, write = stdio_transport
            self.session = await self.exit_stack.enter_async_context(
                ClientSession(read, write)
            )
            await self.session.initialize()

            # Cache tools on connection
            response = await self.session.list_tools()
            self.available_tools = [
                {
                    "type": "function",
                    "function": {
                        "name": tool.name,
                        "description": tool.description,
                        "input_schema": tool.inputSchema,
                    },
                }
                for tool in response.tools
            ]
            logging.info(
                f"Connected to server with tools: {[tool['function']['name'] for tool in self.available_tools]}"
            )
        except Exception as e:
            logging.error(f"Failed to connect to server: {e}")
            await self.cleanup()
            raise

    async def process_query(self, query: str) -> str:
        """Process a query using LLM and available tools

        Args:
            query: The user query to process.

        Returns:
            The response from the LLM or tool.
        """
        if not self.session:
            raise RuntimeError("Server not connected")

        messages = [
            {"role": "system", "content": config.SYSTEM_PROMPT},
            {"role": "user", "content": query},
        ]

        try:
            # Initial model response (only pass tools if available)
            response = self.client.create_chat_completion(
                messages=messages,
                tools=self.available_tools if self.available_tools else None,
            )

            if response.finish_reason == "tool_calls":
                # Parse tool call
                tool_call = response.message.tool_calls[0]
                tool_name = tool_call.function.name
                tool_args = json.loads(tool_call.function.arguments)

                # Execute tool
                logging.info(f"[Calling tool {tool_name} with args {tool_args}]")
                result = await self.session.call_tool(tool_name, tool_args)
                logging.info(f"[Tool result: {result.content}]")

                if not result.content or not isinstance(result.content, list):
                    raise ValueError("Invalid tool response format")

                # Append tool result to messages
                messages.append(response.message.model_dump())
                messages.append(
                    {
                        "role": "tool",
                        "content": result.content[0].text,
                        "tool_call_id": tool_call.id,
                    }
                )

                # Generate final response
                response = self.client.create_chat_completion(messages=messages)

            return response.message.content
        except json.JSONDecodeError:
            logging.error("Failed to parse tool arguments")
            return "Error: Invalid tool request format."
        except Exception as e:
            logging.error(f"Tool execution failed: {e}")
            return f"Error: {str(e)}"

    async def chat_loop(self):
        """Interactive chat loop with the user"""
        logging.info("🤖 MCP Client started! Type 'quit' to exit.")

        while True:
            try:
                query = input("> ").strip()
                if query.lower() == "quit":
                    break

                response = await self.process_query(query)
                logging.info(f"🤖 Assistant: {response}")

            except KeyboardInterrupt:
                break
            except Exception as e:
                logging.error(f"⚠️ Error: {str(e)}")
                import traceback

                traceback.print_exc()

    async def cleanup(self):
        """Cleanup resources and disconnect session."""
        await self.exit_stack.aclose()


async def main():
    client = MCPClient()
    try:
        await client.connect_to_server(
            "/Users/guodongq/Workspaces/src/github.com/guodongq/quickstart/src/python/app/mcp_quickstart/stdio/server.py",
        )
        await client.chat_loop()
    except Exception as e:
        logging.error(f"Fatal error: {e}")
    finally:
        await client.cleanup()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
