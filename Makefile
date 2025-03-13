include ./tools/mk/code.mk
include ./tools/mk/dagger.mk
include ./tools/mk/docker.mk
include ./tools/mk/helm.mk
include ./tools/mk/help.mk
include ./tools/mk/kind.mk

.PHONY: clean 
clean: docker/clean helm/clean kind/clean ## Clean the project
