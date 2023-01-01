.PHONY: all
.DEFAULT_GOAL := help

SERVICE_NAME = ticks
PYTHON_LIB_NAME = cryptellation_$(SERVICE_NAME)

.PHONY: clean
clean: ## Clean
	@echo "TODO"

.PHONY: proto
proto: ## Generate protobuf code from services
ifeq ($(shell python -c "import grpc_tools.protoc" &> /dev/null; echo $$?),1)
	@pip install grpcio-tools
endif
	@mkdir -p $(PYTHON_LIB_NAME)/_genproto
	@python3 \
		-m grpc_tools.protoc \
		--proto_path=../../api/ \
		--python_out=$(PYTHON_LIB_NAME)/_genproto \
		--grpc_python_out=$(PYTHON_LIB_NAME)/_genproto \
		$(SERVICE_NAME).proto
	@sed -i 's/^import .*_pb2 as/from . \0/' $(PYTHON_LIB_NAME)/_genproto/*pb2_grpc.py

.PHONY: lint
lint: ## Lint the code
	@echo "TODO"

.PHONY: format
format: ## Format the code
ifeq ($(shell python -c "import black" &> /dev/null; echo $$?),1)
	@pip install black
endif
	@black $(PWD) --exclude="_genproto/"

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
