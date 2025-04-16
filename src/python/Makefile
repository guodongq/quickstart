CURRENT_DIR=$(shell pwd)
MCP=${CURRENT_DIR}/.venv/bin/mcp
PYTHON=${CURRENT_DIR}/.venv/bin/python

.PHONY: mcp-stdio-server
mcp-stdio-server:
	$(MCP) dev app/mcp_quickstart/stdio/server.py

.PHONY: mcp-stdio-client
mcp-stdio-client:
	$(PYTHON) app/mcp_quickstart/stdio/client.py