## ----- VARIABLES -----
## Program version.
VERSION = latest
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif



## ----- TARGETS ------
## Generic:
.PHONY: default version setup install build run test lint review clean release \
        help

default: run
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

setup: js-setup ## Set up this project on a new device.
	# @echo "Configuring githooks..." && \
	#  git config core.hooksPath .githooks && \
	#  echo done && \

install: js-install ## Install project dependencies.
build: js-build ## Build project.
run: js-run ## Run project (development).
test: ## Run tests.
	@echo "No testing procedure defined."
lint: js-lint ## Lint and check code.
review: js-review ## Lint code and run tests.
clean: js-clean ## Clean build artifacts.
release: ## Release / deploy this project.
	@echo "No release procedure defined."

## Show usage for the targets in this Makefile.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	   awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


# ## CI:
# .PHONY: ci-install ci-test ci-deploy
# ci-install: js-install
# ci-test: js-test
# ci-deploy:
# 	@echo "No deployment procedure defined."


## Javascript:
.PHONY: js-setup js-install js-build js-run js-test js-lint js-review js-clean \
        js-release

js-setup: js-install
js-install:
	@echo "Installing package dependencies..." && yarn install
js-build:
	@echo "Building static files for production..." && yarn build
js-run:
	@echo "Starting development server with..." && yarn dev
js-lint:
	@echo "Linting code with 'yarn lint'..." && yarn lint
js-review: js-lint
js-clean:
	@echo "Cleaning 'node_modules'..." && rm -rf node_modules && echo done
