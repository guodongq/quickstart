import anyio
import pytest
from mcp.client.session import ClientSession
from mcp.server.fastmcp import FastMCP
from mcp.shared.memory import create_client_server_memory_streams

# Define a simple FastMCP server for the test
test_server = FastMCP(name="TestServer")


@test_server.tool()
def ping() -> str:
    return "pong"


@pytest.mark.anyio  # Mark test to be run with anyio
async def test_memory_transport():
    # 1. Use the memory stream generator
    async with create_client_server_memory_streams() as (
            (client_read, client_write),  # Client perspective
            (server_read, server_write)  # Server perspective
    ):
        print("Test: Memory streams created.")
        # Run server and client concurrently
        async with anyio.create_task_group() as tg:
            # 2. Start the server using its streams
            tg.start_soon(
                test_server.run, server_read, server_write,
                test_server.create_initialization_options()
            )
            print("Test: Server started in background task.")

            # 3. Create and run client using its streams
            async with ClientSession(client_read, client_write) as client:
                print("Test: Client session created. Initializing...")
                await client.initialize()
                print("Test: Client initialized. Calling 'ping' tool...")
                result = await client.call_tool("ping")
                print(f"Test: Client received result: {result}")
                # Assert the result is correct
                assert result.content[0].text == "pong"

            # Cancel server task when client is done (optional)
            tg.cancel_scope.cancel()

        print("Test: Finished.")
