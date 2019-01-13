## ----- VARIABLES -----
## Program version.
VERSION = "unset"
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif



## ----- TARGETS ------
## Generic:
.PHONY: default version release help

default: charts-version
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

release: charts-version charts-package

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
__CDIRS = $$(ls -l | egrep '^d' | awk '{print $$NF}')
charts-version:
	@echo "Updating 'appVersion' in 'Chart.yaml' files..." && \
	 for cdir in $(__CDIRS); do \
	   cd $$cdir && \
	   if [ ! -s Chart.yaml ]; then continue; fi && \
	   cp Chart.yaml Chart.tmp.yaml && \
	   sed -e 's/^appVersion: .*$$/appVersion: "$(CVER)"/' \
	     Chart.tmp.yaml > Chart.yaml && \
	   rm Chart.tmp.yaml && \
	   git add Chart.yaml; \
	 done && \
	 echo done

charts-package:
	@echo "Packaging charts..." && \
	 for cdir in $(__CDIRS); do cd $$cdir && helm package .; done && \
	 echo done
