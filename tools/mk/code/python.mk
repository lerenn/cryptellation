PYTHON_GATEWAY_CLIENT_DIR := ./clients/python/gateway
PYTHON_GATEWAY_VENV_PATH  := $(PYTHON_GATEWAY_CLIENT_DIR)/.venv

PYTHON_E2E_CLIENT_DIR := ./test/python
PYTHON_E2E_VENV_PATH  := $(PYTHON_E2E_CLIENT_DIR)/.venv

.PHONY: _python/e2e/dependencies
_python/e2e/dependencies: ## Install the dependencies
	@if [ ! -d "$(PYTHON_E2E_VENV_PATH)" ]; then \
		python3 -m venv --prompt cryptellation-e2e $(PYTHON_E2E_VENV_PATH) && \
		. $(PYTHON_E2E_VENV_PATH)/bin/activate && \
		pip install -Ur $(PYTHON_E2E_CLIENT_DIR)/requirements.txt; \
	fi


.PHONY: _python/gateway/dependencies
_python/gateway/dependencies: ## Install the dependencies
	@if [ ! -d "$(PYTHON_GATEWAY_VENV_PATH)" ]; then \
		python3 -m venv --prompt cryptellation $(PYTHON_GATEWAY_VENV_PATH) && \
		. $(PYTHON_GATEWAY_VENV_PATH)/bin/activate && \
		pip install -Ur $(PYTHON_GATEWAY_CLIENT_DIR)/requirements.dev.txt; \
	fi

.PHONY: python/generate
python/generate: _python/gateway/dependencies ## Generate the code
	@. $(PYTHON_GATEWAY_VENV_PATH)/bin/activate && openapi-python-client generate \
		--path ./api/gateway/v1.yaml \
		--output-path $(PYTHON_GATEWAY_CLIENT_DIR) \
		--config $(PYTHON_GATEWAY_CLIENT_DIR)/openapi.config.yaml \
		--overwrite

.PHONY: python/test
python/test: python/test/end-to-end ## Launch all tests

.PHONY: python/test/end-to-end
python/test/end-to-end: docker/gateway/up _python/e2e/dependencies ## Launch end-to-end tests
	@. $(PYTHON_E2E_VENV_PATH)/bin/activate && python3 $(PYTHON_E2E_CLIENT_DIR)/tests.py