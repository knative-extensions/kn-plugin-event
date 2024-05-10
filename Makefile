CI ?= false
PROW_JOB_ID ?= 0
PROW = $(shell [ "$(CI)" = "true" ] && [ "$(PROW_JOB_ID)" != "0" ] && echo "true" || echo "false")

ifeq ($(PROW),true)
	export HOME = /tmp
	export TEST_IMAGES_DIR = /usr/bin
endif

build:
	./mage build
.PHONY: build

unit:
	./mage test
.PHONY: unit

e2e:
	openshift/e2e-tests.sh
.PHONY: e2e
