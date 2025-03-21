PYTHON_CLIENT_DIR := ./clients/python
VENV_PATH         := $(PYTHON_CLIENT_DIR)/venv

.PHONY: _python/dependencies
_python/dependencies: ## Install the dependencies
	@if [ ! -d "$(VENV_PATH)" ]; then \
		python3 -m venv --prompt cryptellation $(PYTHON_CLIENT_DIR)/venv && \
		. $(VENV_PATH)/bin/activate && \
		pip install -Ur $(PYTHON_CLIENT_DIR)/requirements/dev.txt; \
	fi

.PHONY: python/generate
python/generate: _python/dependencies ## Generate the code
	@. $(VENV_PATH)/bin/activate && \
		openapi-python-client generate --path ./api/gateway/v1.yaml --output-path $(PYTHON_CLIENT_DIR)/cryptellation