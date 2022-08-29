.PHONY: all
.DEFAULT_GOAL := help

clean: ## Clean

proto: ## Generate protobuf code from services
	@mkdir -p cryptellation/_genproto
	@for PROTOFILE in $(shell find ../../ -name *.proto); do \
		echo "Generating $$PROTOFILE"; \
		python3 \
			-m grpc_tools.protoc \
			--proto_path=$$(dirname $$PROTOFILE) \
			--python_out=cryptellation/_genproto \
			--grpc_python_out=cryptellation/_genproto \
			"$$PROTOFILE" || exit $? || exit $?; \
	done
	@sed -i 's/^import .*_pb2 as/from . \0/' cryptellation/_genproto/*pb2_grpc.py

lint: ## Lint the code

format: ## Format the code
	@black $(PWD) --exclude="_genproto/"

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
