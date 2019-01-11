## ----- VARIABLES -----
## Go module name.
MODULE = $(shell basename "$$(pwd)")
ifeq ($(shell ls -1 go.mod 2> /dev/null),go.mod)
	MODULE = $(shell cat go.mod | grep module | awk '{print $$2}')
endif

## Program version.
VERSION = "unset"
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif

## Custom Go linker flag.
LDFLAGS = -X $(MODULE)/internal/info.Version=$(VERSION)



## ----- TARGETS ------
## Generic:
.PHONY: default version setup install build clean run lint test review release \
        help

default: run
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

setup: go-setup ## Set up this project on a new device.
	@echo "Configuring githooks..." && \
	 git config core.hooksPath .githooks && \
	 echo done && \
	 echo "Installing client dependencies..." && cd client && yarn install

install: go-install ## Install project dependencies.
	@echo "Installing client dependencies..." && cd client && yarn install

build: go-build ## Build project.
clean: go-clean ## Clean build artifacts.
run: go-run ## Run project (development).
	@echo "Starting client..." && cd client && yarn start

lint: go-lint ## Lint and check code.
	@echo "Linting client code..." && cd client && yarn lint

test: go-test ## Run tests.
review: go-review ## Lint code and run tests.
	@echo "Linting client code..." && cd client && yarn lint

release: ## Release / deploy this project.
	@echo "No release procedure defined."

## Show usage for the targets in this Makefile.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	   awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


## Go:
.PHONY: go-deps go-bench go-setup go-install go-build go-clean go-run go-lint \
        go-test go-review

go-deps: ## Verify and tidy project dependencies.
	@echo "Verifying module dependencies..." && \
	 go mod verify && \
	 echo "Tidying module dependencies..." && \
	 go mod tidy && \
	 echo done

go-bench: ## Run benchmarks.
	@echo "Running benchmarks with 'go test -bench=.'..." && \
	 $(__TEST) -run=^$$ -bench=. -benchmem ./...

go-setup: go-install go-deps
go-install:
	@echo "Downloading module dependencies..." && \
	 go mod download && \
	 echo done

BUILDARGS = -ldflags "$(LDFLAGS)" $(BARGS)
BDIR = .
go-build:
	@echo "Building with 'go build $(BARGS)'..." && \
	 go build $(BUILDARGS) $(BDIR) && \
	 echo done

go-clean:
	@echo "Cleaning with 'go clean'..." && \
	 go clean $(BDIR) && \
	 echo done

go-run:
	@echo "Running with 'go run'..." && \
	 go run $(BUILDARGS) $(BDIR)

go-lint:
	@if command -v goimports > /dev/null; then \
	   echo "Formatting code with 'goimports'..." && goimports -w .; \
	 else \
	   echo "'goimports' not installed, formatting code with 'go fmt'..." && \
	   go fmt .; \
	 fi && \
	 if command -v golint > /dev/null; then \
	   echo "Linting code with 'golint'..." && golint ./...; \
	 else \
	   echo "'golint' not installed, skipping linting step."; \
	 fi && \
	 echo "Checking code with 'go vet'..." && go vet ./... && \
	 echo done

COVERFILE = coverage.out
TIMEOUT   = 20s
__TEST = go test -coverprofile="$(COVERFILE)" -covermode=atomic \
               -timeout="$(TIMEOUT)" $(TARGS)
go-test:
	@echo "Running tests with 'go test':" && $(__TEST) ./....

go-review: go-lint go-test


## git-secret:
.PHONY: secrets-hide secrets-reveal
secrets-hide: ## Hides modified secret files using git-secret.
	@echo "Hiding modified secret files..." && git secret hide -m

secrets-reveal: ## Reveals secret files that were hidden using git-secret.
	@echo "Revealing hidden secret files..." && git secret reveal


## CI:
.PHONY: ci-install ci-test ci-deploy
ci-install: go-install dk-pull
ci-test: dk-test
	@$(DKCMP_VERSION) up --no-start && make dk-tags

ci-deploy:
	@make dk-push && \
	 for deploy in $(DEPLOYS); do \
	   kubectl patch deployment "$$deploy" \
	     -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"$$(date +'%s')\"}}}}}"; \
	 done


## Docker:
.PHONY: dk-pull dk-push dk-build dk-build-push dk-clean dk-tags dk-up \
        dk-build-up dk-down dk-logs dk-test

DK    = docker
DKCMP = docker-compose
DKCMP_VERSION = VERSION="$(VERSION)" $(DKCMP)
DKCMP_LATEST  = VERSION=latest $(DKCMP)

dk-pull: ## Pull latest Docker images from registry.
	@echo "Pulling latest images from registry..." && \
	 $(DKCMP_LATEST) pull $(SVC)

dk-push: ## Push new Docker images to registry.
	@echo "Pushing images to registry..." && \
	 $(DKCMP_VERSION) push $(SVC) && \
	 $(DKCMP_LATEST) push $(SVC) && \
	 echo done

dk-build: ## Build and tag Docker images.
	@echo "Building images..." && \
	 $(DKCMP_VERSION) build --parallel --compress $(SVC) && \
	 echo done && make dk-tags

dk-clean: ## Clean up unused Docker data.
	@echo "Cleaning unused data..." && $(DK) system prune

dk-build-push: dk-build dk-push ## Build and push new Docker images.
dk-tags: ## Tag versioned Docker images with ':latest'.
	@echo "Tagging versioned images with ':latest'..." && \
	IMAGES="$$($(DKCMP_VERSION) config | egrep image | awk '{print $$2}')" && \
	for image in $$IMAGES; do \
	  if [ -z "$$($(DK) images -q "$$image" 2> /dev/null)" ]; then \
	    continue; \
	  fi && \
	  LAT_TAG="$$(echo "$$image" | sed -e 's/:.*$$/:latest/')" && \
	  $(DK) tag "$$image" "$$LAT_TAG"; \
	done && \
	echo done

__DK_UP = $(DKCMP_VERSION) up -d
dk-up: ## Start up containerized services.
	@echo "Bringing up services..." && $(__DK_UP) $(SVC) && echo done

dk-build-up: ## Build new images, then start them.
	@echo "Building and bringing up services..." && \
	 $(__DK_UP) --build $(SVC) && \
	 echo done

dk-down: ## Shut down containerized services.
	@echo "Bringing down services..." && \
	 $(DKCMP_VERSION) down $(SVC) && \
	 echo done

dk-logs: ## Show logs for containerized services.
	@$(DKCMP_VERSION) logs $(SVC)

DKCMP_TEST = $(DKCMP_VERSION) -f docker-compose.test.yml
dk-test: ## Test using 'docker-compose.test.yml'.
	@if [ -s docker-compose.test.yml ]; then \
	   echo "Running containerized service tests..." && \
	   $(DKCMP_TEST) up; \
	 fi
