K8S_VERSION    := 1.24.10
MINIKUBE_NODES := 1
NODE_MEMORY    := 2g
CLUSTER_NAME   := cryptellation-devel

.PHONY: deploy
deploy: deploy/infra ## Deploys everything, except Cryptellation System

.PHONY: deploy/infra
deploy/infra: ## Deploys local development infrastructure environment
	@minikube start \
		--cpus=2 \
		--memory=${NODE_MEMORY} \
		--nodes=${MINIKUBE_NODES} \
		--profile ${CLUSTER_NAME} \
		--kubernetes-version=v${K8S_VERSION} \
		--extra-config=kubelet.cluster-dns=10.96.0.10

.PHONY: delete
delete: delete/infra ## deletes the dev environment

.PHONY: delete/infra
delete/infra: ## deletes local dev infrastructure
	@minikube delete --profile ${CLUSTER_NAME}

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

