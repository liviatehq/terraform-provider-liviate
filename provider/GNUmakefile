# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=liviate

OS      ?= $(shell go env GOOS)
ARCH    ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
EXE     := $(if $(filter windows,$(OS)),.exe,)
BINARY  := terraform-provider-$(PKG_NAME)$(EXE)
DIST    := terraform-provider-$(PKG_NAME)_v$(VERSION)$(EXE)
BIN_DIR := bin

# Terraform local plugin directory
ifeq ($(OS),windows)
  TERRAFORM_PLUGIN_DIR ?= $(if $(APPDATA),$(APPDATA)/terraform.d/plugins,$(HOME)/.terraform.d/plugins)
else
  TERRAFORM_PLUGIN_DIR ?= $(HOME)/.terraform.d/plugins
endif

PLUGIN_ARCH_DIR := $(TERRAFORM_PLUGIN_DIR)/registry.terraform.io/$(PKG_NAME)/$(PKG_NAME)/$(VERSION)/$(OS)_$(ARCH)

default: build

build: fmtcheck
	go build -o "$(BIN_DIR)/$(BINARY)" .

install: build
	mkdir -p "$(PLUGIN_ARCH_DIR)"
	cp "$(BIN_DIR)/$(BINARY)" "$(PLUGIN_ARCH_DIR)/$(DIST)"
	@echo "Installed $(DIST) -> $(PLUGIN_ARCH_DIR)"

deploy: install

uninstall:
	rm -rf "$(PLUGIN_ARCH_DIR)"

clean:
	rm -rf "$(BIN_DIR)"

sync-forge:
	@bash "$(CURDIR)/scripts/forge-sync.sh" push

sync-forge-dry-run:
	@bash "$(CURDIR)/scripts/forge-sync.sh" --dry-run push

pull-forge:
	@bash "$(CURDIR)/scripts/forge-sync.sh" pull

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 30m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@bash -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@bash -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build install deploy uninstall clean sync-forge sync-forge-dry-run pull-forge test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test

