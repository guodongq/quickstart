CURRENT_DIR=$(shell pwd)
MCP=${CURRENT_DIR}/.venv/bin/mcp
PYTHON=${CURRENT_DIR}/.venv/bin/python

### UV management
.PHONY: uv_sync
uv_sync:
	uv sync

### Ruff management
.PHONY: ruff-check
ruff-check:
	$(PYTHON) -m ruff check

.PHONY: ruff-fix
ruff-fix:
	$(PYTHON) -m ruff check --fix

.PHONY: ruff-format
ruff-format:
	$(PYTHON) -m ruff format

### MCP management
.PHONY: mcp-stdio-server
mcp-stdio-server:
	$(MCP) dev app/mcp_quickstart/stdio/server.py

.PHONY: mcp-stdio-client
mcp-stdio-client:
	uv run app/mcp_quickstart/stdio/client.py