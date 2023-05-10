K8S_VERSION         := 1.26.1
MINIKUBE_NODES      := 3
NODE_MEMORY         := 2g
CLUSTER_NAME        := cryptellation-devel

DOCKER_IMAGE_PREFIX := lerenn/cryptellation

.PHONY: deploy/local
deploy/local: deploy/local/infra deploy/local/cryptellation ## Deploys everything locally

.PHONY: deploy/local/infra
deploy/local/infra: ## Deploys local development infrastructure environment
	@minikube start \
		--cpus=2 \
		--memory=${NODE_MEMORY} \
		--nodes=${MINIKUBE_NODES} \
		--profile ${CLUSTER_NAME} \
		--kubernetes-version=v${K8S_VERSION} \
		--extra-config=kubelet.cluster-dns=10.96.0.10
	@minikube -p $(CLUSTER_NAME) addons enable registry

.PHONY: deploy/local/cryptellation
deploy/local/cryptellation:
	@kubectl port-forward --namespace kube-system service/registry 5000:80 &> /dev/null & FORWARD_PID="$$!"; \
		$(MAKE) -C ../../build/package push/local; \
		kill -TERM $$FORWARD_PID
	@helm upgrade cryptellation ./cryptellation -f ./local.yml \
		-n cryptellation --create-namespace

.PHONY: local/status
local/status: ## Checks the status of the dev environment
	@minikube status --profile ${CLUSTER_NAME}

.PHONY: local/stop
local/stop:
	@minikube stop -p ${CLUSTER_NAME}

.PHONY: delete/local
delete/local: delete/local/cryptellation delete/local/infra ## Deletes the dev environment

.PHONY: delete/local/cryptellation
delete/local/cryptellation:
	@kubectl delete all --all --namespace cryptellation || true

.PHONY: delete/local/infra
delete/local/infra:
	@minikube delete --profile ${CLUSTER_NAME}

.PHONY: clean
clean: delete/local ## Clean everything

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

