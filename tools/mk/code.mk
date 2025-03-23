MK_TOOLS_PATH := ./tools/mk

include $(MK_TOOLS_PATH)/code/go.mk
include $(MK_TOOLS_PATH)/code/python.mk

.PHONY: code/check
code/check: code/generate code/lint code/test ## Generate, lint and test the code

.PHONY: code/check-todos
code/check-todos: ## Check the todos in the code
	@go run ./tools/invtodos .

.PHONY: code/generate
code/generate: go/generate ## Generate the code

.PHONY: code/lint
code/lint: go/lint ## Lint the code

.PHONY: code/test
code/test: code/test/unit code/test/integration code/test/end-to-end ## Launch all tests

.PHONY: code/test/unit
code/test/unit: go/test/unit ## Launch unit tests

.PHONY: code/test/integration
code/test/integration: go/test/integration ## Launch integration tests

.PHONY: code/test/end-to-end
code/test/end-to-end: go/test/end-to-end python/test/end-to-end ## Launch end-to-end tests