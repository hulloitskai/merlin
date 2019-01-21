## ----- VARIABLES -----
## Program version.
VERSION = latest
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif



## ----- TARGETS ------
## Generic:
.PHONY: default version lint clean release help

default: charts-version
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

lint: charts-lint ## Lint and check code.
clean: charts-clean ## Clean build artifacts.
release: charts-version charts-package ## Release / deploy this project.

## Show usage for the targets in this Makefile.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	   awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


## Helm Charts:
.PHONY: charts-version charts-package

CVER = unset
ifneq ($(VERSION),unset)
	CVER = $(shell git describe --abbrev=0 --tags | cut -c 2-)
endif
__CDIRS = $(shell \
	unset vcdirs && \
	for cdir in $$(ls -l | egrep '^d' | awk '{print $$NF}'); do \
	  if [ -s "$$cdir/Chart.yaml" ]; then \
	    vcdirs="$$vcdirs $$cdir"; \
	  fi; \
	done && \
	echo "$$vcdirs" \
)

charts-version: ## Set chart appVersion string to the current project version.
	@echo "Updating 'appVersion' in 'Chart.yaml' files..." && \
	 for cdir in $(__CDIRS); do \
	   cd $$cdir && cp Chart.yaml Chart.tmp.yaml && \
	   sed -e 's/^appVersion: .*$$/appVersion: "$(CVER)"/' \
	     Chart.tmp.yaml > Chart.yaml && \
	   rm Chart.tmp.yaml && \
	   git add Chart.yaml; \
	 done && \
	 echo done

charts-lint:
	@echo "Linting charts with Helm..." && \
	 for cdir in $(__CDIRS); do helm lint $$cdir; done && \
	 echo done

charts-package: charts-lint ## Package charts into release archives.
	@echo "Packaging charts..." && \
	 for cdir in $(__CDIRS); do helm package "$$cdir"; done && \
	 echo done

charts-clean:
	@echo "Cleaning up chart packages..." && \
	 (ls -1 *.tgz 2> /dev/null && rm *.tgz || true) && \
	 echo done
