from langgraph.prebuilt import create_react_agent
from langchain.chat_models import init_chat_model


def get_weather(city: str) -> str:
    """"Get weather for a given city."""
    return f"It's always sunny in {city}"


model = init_chat_model(
    "qwen2.5:14b",
    model_provider="ollama",
    temperature=0,
)

agent = create_react_agent(
    model=model,
    tools=[get_weather],
    prompt="You are a helpful assistant"
)

if __name__ == '__main__':
    # Run the agent
    result = agent.invoke(
        {"message": [{"role": "user", "content": "What is the weather in Shanghai?"}]}
    )
    print("\n", result)
