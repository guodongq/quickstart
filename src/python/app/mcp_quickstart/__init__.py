import json
import logging
from abc import ABC, abstractmethod
from contextlib import AsyncExitStack
from types import SimpleNamespace
from typing import List, Dict, Any, AsyncGenerator

from openai import OpenAI
from mcp import ClientSession
from contextlib import asynccontextmanager

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(filename)s:%(lineno)d - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S"
)

logger = logging.getLogger(__name__)


class AppConstants(SimpleNamespace):
    REQUIRED_KEYS = {"API_KEY", "BASE_URL", "TIMEOUT", "SYSTEM_PROMPT"}

    def __init__(self, **kwargs):
        missing_keys = self.REQUIRED_KEYS - set(kwargs.keys())
        if missing_keys:
            raise ValueError(f"Missing required keys: {missing_keys}")
        super().__init__(**kwargs)

    def __setattr__(self, key, value):
        if hasattr(self, key):
            raise AttributeError(f"Constant '{key}' cannot be modified")
        super().__setattr__(key, value)


CONFIG = AppConstants(
    API_KEY="ollama",
    BASE_URL="http://localhost:11434/v1/",
    TIMEOUT=3000,
    SYSTEM_PROMPT=(
        "You are a helpful AI assistant capable of answering user questions."
        "When necessary, you can use mcp tools to gather additional context "
        "and provide more accurate, complete responses."
    )
)


class TransportStrategy(ABC):
    @asynccontextmanager
    @abstractmethod
    async def create_session(self) -> AsyncGenerator[ClientSession, None]:
        yield


class MCPClient:
    def __init__(self, model: str = "llama3.2:latest"):
        self.model: str = model
        self.client: OpenAI = OpenAI(
            api_key=CONFIG.API_KEY,
            base_url=CONFIG.BASE_URL,
            timeout=CONFIG.TIMEOUT,
        )
        self.available_mcp_tools: List[Dict[str, Any]] = []
        self.exit_stack: AsyncExitStack = AsyncExitStack()
        self.session: ClientSession | None = None

    async def connect_to_server(self, transport_strategy: TransportStrategy):
        self.session = await self.exit_stack.enter_async_context(
            transport_strategy.create_session()
        )

        tool_response = await self.session.list_tools()
        self.available_mcp_tools = [
            {
                "type": "function",
                "function": {
                    "name": tool.name,
                    "description": tool.description,
                    "input_schema": tool.inputSchema,
                }
            }
            for tool in tool_response.tools
        ]

    async def _process_question(self, question: str) -> str:
        if not self.session:
            raise RuntimeError("Session not initialized")

        messages = [
            {"role": "system", "content": CONFIG.SYSTEM_PROMPT},
            {"role": "user", "content": question}
        ]

        try:
            response = self.client.chat.completions.create(
                model=self.model,
                messages=messages,
                tools=self.available_mcp_tools if self.available_mcp_tools else None,
            )

            if not response.choices:
                raise ValueError("No choices in response")

            content = response.choices[0]
            if content.finish_reason == "tool_calls" and content.message.tool_calls:
                # Parse tool call
                tool_call = content.message.tool_calls[0]
                tool_name = tool_call.function.name
                try:
                    tool_args = json.loads(tool_call.function.arguments)
                except json.JSONDecodeError as e:
                    raise ValueError(f"Invalid tool arguments: {e}")

                # Execute tool call
                logger.info(f"[Calling tool {tool_name} with args {tool_args}]")
                result = await self.session.call_tool(tool_name, tool_args)
                logger.info(f"[Tool result: {result.content}]")

                # Append tool result to messages
                messages.append(content.message.model_dump())
                messages.append({
                    "role": "tool",
                    "content": result.content[0].text if result.content else "",
                    "tool_call_id": tool_call.id,
                })

                # Generate final response
                response = self.client.chat.completions.create(
                    model=self.model,
                    messages=messages
                )
                if not response.choices:
                    raise ValueError("No choices in final response")

            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"Error processing question: {e}")
            return f"Sorry, an error occurred while processing your request: {str(e)}"

    async def chat_loop(self):
        logger.info("OPEN AI Client started. Type 'quit' to exit.")
        while True:
            try:
                question = input("> ").strip()
                if question.lower() in ("quit", "exit"):
                    break

                response = await self._process_question(question)
                logger.info(f"Response: {response}")
            except KeyboardInterrupt:
                logger.info("Received interrupt, exiting...")
                break
            except Exception as e:
                logger.error(f"Error: {str(e)}")

                import traceback
                traceback.print_exc()

    async def cleanup(self):
        await self.exit_stack.aclose()


async def main(transport_strategy: TransportStrategy):
    client = MCPClient()
    try:
        await client.connect_to_server(transport_strategy)
        await client.chat_loop()
    except KeyboardInterrupt:
        logger.info("Exiting...")
    except Exception as e:
        logger.error(f"Error occurred: {e}")
    finally:
        await client.cleanup()
