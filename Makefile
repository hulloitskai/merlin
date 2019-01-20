## ----- VARIABLES -----
## Program version.
VERSION = latest
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif



## ----- TARGETS ------
## Generic:
.PHONY: default version setup install build clean run lint test review release \
        help

default: run
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

setup: ## Set up this project on a new device.
	@echo "Configuring githooks..." && \
	 git config core.hooksPath .githooks && \
	 echo done
	@cd api && $(MAKE) install
	@cd client && $(MAKE) install

install: ## Install project dependencies.
	@cd api && $(MAKE) install
	@cd client && $(MAKE) install

build: ## Build project.
	@cd api && $(MAKE) build

clean: ## Clean build artifacts.
	@cd api && $(MAKE) clean
	@cd client && $(MAKE) clean

run: ## Run project (development).
	@if command -v parallel > /dev/null; then \
	   parallel --lb --tagstring '{= s:^cd ::; s: &&.*$$::; =}' ::: \
	     "cd api && make run" "cd client && make run"; \
	 else \
	   echo "Cannot run client and api in parallel ('GNU parallel' not found)."; \
	 fi

lint: ## Lint and check code.
	@cd api && $(MAKE) lint
	@cd client && $(MAKE) lint

test: ## Run tests.
	@cd api && $(MAKE) test

review: ## Lint code and run tests.
	@cd api && $(MAKE) review
	@cd client && $(MAKE) review

release: ## Release / deploy this project.
	@echo "No release procedure defined."

## Show usage for the targets in this Makefile.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	   awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


## CI:
.PHONY: ci-install ci-test ci-deploy
__KB = kubectl

ci-install:
	cd api && $(MAKE) install
	$(MAKE) dk-pull
ci-test: dk-test
	@$(__DKCMP_VER) up --no-start && $(MAKE) dk-tags
ci-deploy:
	@$(MAKE) dk-push && \
	 for deploy in $(DEPLOYS); do \
	   $(__KB) patch deployment "$$deploy" \
	     -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"$$(date +'%s')\"}}}}}"; \
	 done


# ## git-secret:
# .PHONY: secrets-hide secrets-reveal
# secrets-hide: ## Hides modified secret files using git-secret.
# 	@echo "Hiding modified secret files..." && git secret hide -m

# secrets-reveal: ## Reveals secret files that were hidden using git-secret.
# 	@echo "Revealing hidden secret files..." && git secret reveal


## Docker:
.PHONY: dk-pull dk-push dk-build dk-build-push dk-clean dk-tags dk-up \
        dk-build-up dk-down dk-logs dk-test
__DK    = docker
__DKCMP = docker-compose
__DKCMP_VER = VERSION="$(VERSION)" $(__DKCMP)
__DKCMP_LST  = VERSION=latest $(__DKCMP)

dk-pull: ## Pull latest Docker images from registry.
	@echo "Pulling latest images from registry..." && \
	 $(__DKCMP_LST) pull $(SVC)

dk-push: ## Push new Docker images to registry.
	@if git describe --exact-match --tags > /dev/null 2>&1; then \
	   echo "Pushing versioned images to registry (:$(VERSION))..." && \
	   $(__DKCMP_VER) push $(SVC); \
	 fi && \
	 echo "Pushing latest images to registry (:latest)..." && \
	 $(__DKCMP_LST) push $(SVC) && \
	 echo done

dk-build: ## Build and tag Docker images.
	@echo "Building images..." && \
	 $(__DKCMP_VER) build --parallel --compress $(SVC) && \
	 echo done && $(MAKE) dk-tags

dk-clean: ## Clean up unused Docker data.
	@echo "Cleaning unused data..." && $(__DK) system prune

dk-build-push: dk-build dk-push ## Build and push new Docker images.

dk-tags: ## Tag versioned Docker images with ':latest'.
	@echo "Tagging versioned images with ':latest'..." && \
	images="$$($(__DKCMP_VER) config | egrep image | awk '{print $$2}')" && \
	for image in $$images; do \
	  if [ -z "$$($(__DK) images -q "$$image" 2> /dev/null)" ]; then \
	    continue; \
	  fi && \
	  echo "$$image" | sed -e 's/:.*$$/:latest/' | \
	    xargs $(__DK) tag "$$image"; \
	done && \
	echo done

__DK_UP = $(__DKCMP_VER) up -d
dk-up: ## Start up containerized services.
	@echo "Bringing up services..." && $(__DK_UP) $(SVC) && echo done

dk-build-up: ## Build new images, then start them.
	@echo "Building and bringing up services..." && \
	 $(__DK_UP) --build $(SVC) && \
	 echo done

dk-down: ## Shut down containerized services.
	@echo "Bringing down services..." && \
	 $(__DKCMP_VER) down $(SVC) && \
	 echo done

dk-logs: ## Show logs for containerized services.
	@$(__DKCMP_VER) logs -f $(SVC)

__DKCMP_TEST = $(__DKCMP_VER) -f docker-compose.test.yml
dk-test: ## Test using 'docker-compose.test.yml'.
	@if [ -s docker-compose.test.yml ]; then \
	   echo "Running containerized service tests..." && \
	   $(__DKCMP_TEST) up; \
	 fi
