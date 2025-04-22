from openai import OpenAI
from openai.types.chat.chat_completion import Choice

from .constants import Constants

config = Constants()
config.API_KEY = "ollama"
config.BASE_URL = "http://localhost:11434/v1/"
config.TIMEOUT = 3000
config.SYSTEM_PROMPT = (
    "You are a helpful AI assistant capable of answering user questions."
    "When necessary, you can use mcp tools to gather additional context "
    "and provide more accurate, complete responses."
)


class OpenAIClient:
    def __init__(self, model: str = "llama3.2:latest"):
        self.model = model
        self.client = OpenAI(
            api_key=config.API_KEY,
            base_url=config.BASE_URL,
            timeout=config.TIMEOUT,
        )

    def create_chat_completion(self, messages, tools=None) -> Choice:
        params = {
            "model": self.model,
            "messages": messages,
        }
        if tools is not None:
            params["tools"] = tools
        response = self.client.chat.completions.create(**params)
        return response.choices[0]
