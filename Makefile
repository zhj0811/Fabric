# Copyright PeerFintech Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets

#   - all (default) - builds all binary
#   - apiserver - builds a native apiserver binary
#   - eventserver - builds a native eventserver binary
#   - apiserver-docker[-clean] - ensures the apiserver container is available[/cleaned]
#   - eventserver-docker[-clean] - ensures the eventserver container is available[/cleaned]
#   - clean - cleans the build area
#   - docker[-clean] - ensures all docker images are available[/cleaned]
#   - docker-list - generates a list of docker images that 'make docker' produces

BASE_VERSION = 1.1.1.0
PREV_VERSION = 1.0.4

PROJECT_COMPANY = peerfintech
PROJECT_NAME = factoring
PKGNAME = github.com/zhj0811/$(PROJECT_NAME)
CGO_FLAGS = CGO_CFLAGS=" "
GO_TAGS ?=
DBUILD = docker build

IMAGES = apiserver eventserver
IS_RELEASE = true

export GO_LDFLAGS GO_TAGS

ifneq ($(IS_RELEASE),true)
EXTRA_VERSION ?= snapshot-$(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(BASE_VERSION)-$(EXTRA_VERSION)
else
PROJECT_VERSION=$(BASE_VERSION)
endif

pkgmap.apiserver    := $(PKGNAME)/apiserver
pkgmap.eventserver  := $(PKGNAME)/eventserver

GO_LDFLAGS=-X $(PKGNAME)/common/metadata.Version=$(PROJECT_VERSION)
 
all: native 
native: apiserver eventserver

.PHONY: apiserver
apiserver: build/bin/apiserver
apiserver-docker: build/image/apiserver

.PHONY: eventserver
eventserver: build/bin/eventserver
eventserver-docker: build/image/eventserver

build/bin/%: $(PROJECT_FILES)
	@mkdir -p $(@D)
	@echo "$@"
	$(CGO_FLAGS) GOBIN=$(abspath $(@D)) go install -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))
	@echo "Binary available as $@"
	@touch $@

.PRECIOUS: build/image/%/payload
build/image/%/payload: build/bin/% images/packages/*
	mkdir -p $@
	cp $^ $@

.PRECIOUS: build/image/%/Dockerfile
build/image/%/Dockerfile: images/%/Dockerfile.in
	mkdir -p $(@D)
	cp $< $@

#build/image/%: Makefile build/image/%/payload build/image/%/Dockerfile
build/image/%: build/image/%/payload build/image/%/Dockerfile
	$(eval TARGET = $(@F))
	@echo "Building docker $(TARGET)-image"
	$(DBUILD) -t $(PROJECT_COMPANY)/$(PROJECT_NAME)-$(TARGET) $(@)
	docker tag $(PROJECT_COMPANY)/$(PROJECT_NAME)-$(TARGET) $(PROJECT_COMPANY)/$(PROJECT_NAME)-$(TARGET):$(BASE_VERSION)
	@touch $@

docker: $(patsubst %,build/image/%, $(IMAGES))

%-docker-list:
	$(eval TARGET = ${patsubst %-docker-list,%,${@}})
	@echo $(PROJECT_COMPANY)/$(PROJECT_NAME)-$(TARGET):$(BASE_VERSION)

docker-list: $(patsubst %,%-docker-list, $(IMAGES))

%-docker-clean:
	$(eval TARGET = ${patsubst %-docker-clean,%,${@}})
	-docker images -q $(PROJECT_COMPANY)/$(PROJECT_NAME)-$(TARGET):$(BASE_VERSION) | xargs -I '{}' docker rmi -f '{}'
	-@rm -rf build/image/$(TARGET) ||:

docker-clean: $(patsubst %,%-docker-clean, $(IMAGES))

.PHONY: clean
clean:
	-@rm -rf build ||:
