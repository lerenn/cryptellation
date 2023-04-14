K8S_VERSION         := 1.26.1
MINIKUBE_NODES      := 3
NODE_MEMORY         := 2g
CLUSTER_NAME        := cryptellation-devel

DOCKER_IMAGE_PREFIX := digitalfeather/cryptellation

.PHONY: deploy
deploy: deploy/infra deploy/datastore deploy/event-broker ## Deploys everything, except Cryptellation System

.PHONY: deploy/all
deploy/all: deploy deploy/cryptellation ## Deploys everything, with Cryptellation System

.PHONY: deploy/infra
deploy/infra: ## Deploys local development infrastructure environment
	@minikube start \
		--cpus=2 \
		--memory=${NODE_MEMORY} \
		--nodes=${MINIKUBE_NODES} \
		--profile ${CLUSTER_NAME} \
		--kubernetes-version=v${K8S_VERSION} \
		--extra-config=kubelet.cluster-dns=10.96.0.10
	@minikube -p $(CLUSTER_NAME) addons enable registry

.PHONY: deploy/datastore
deploy/datastore:
	@kubectl create ns datastore || true
	@helm repo add cockroachdb https://charts.cockroachdb.com/; helm repo update
	@helm upgrade --install cockroachdb cockroachdb/cockroachdb --values datastore/cockroachdb.values.yaml --namespace datastore
	@kubectl create configmap init-databases-scripts --from-file=../../scripts/init-databases.sql --namespace datastore || true
	@kubectl apply -f datastore/cockroachdb.init-databases.job.yaml --namespace datastore || true

.PHONY: deploy/event-broker
deploy/event-broker:
	@kubectl create ns event-broker || true
	@helm repo add nats https://nats-io.github.io/k8s/helm/charts/; helm repo update
	@helm upgrade --install nats nats/nats --values event-broker/nats.values.yaml --namespace event-broker

.PHONY: deploy/cryptellation
deploy/cryptellation:
	@kubectl port-forward --namespace kube-system service/registry 5000:80 &> /dev/null & FORWARD_PID="$$!"; \
		$(MAKE) -C ../../build/package push/local; \
		kill -TERM $$FORWARD_PID
	@helm upgrade --install cryptellation ../../deployments/helm/cryptellation \
		-f ./cryptellation/values.yml \
		-n cryptellation \
		--create-namespace \
		--wait  | sed 's/^/    /'

.PHONY: status
status: ## Checks the status of the dev environment
	@minikube status --profile ${CLUSTER_NAME}

.PHONY: stop
stop:
	@minikube stop -p ${CLUSTER_NAME}

.PHONY: delete
delete: delete/cryptellation delete/event-broker delete/datastore delete/infra ## Deletes the dev environment

.PHONY: delete/cryptellation
delete/cryptellation:
	@kubectl delete all --all --namespace cryptellation || true

.PHONY: delete/event-broker
delete/event-broker: 
	@kubectl delete all --all --namespace event-broker || true

.PHONY: delete/datastore
delete/datastore:
	@kubectl delete all --all --namespace datastore || true

.PHONY: delete/infra
delete/infra:
	@minikube delete --profile ${CLUSTER_NAME}

.PHONY: clean
clean: delete ## Clean everything

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

