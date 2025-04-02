###########################################################################
# Copyright 2021. Ivan Eroshkin
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
###########################################################################

PROJECT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
COVERPROFILE_PATH=$(PROJECT_DIR)/build/coverage.txt

all: help

test-ci: # @HELP runs CI/CD pipeline locally
test-ci: go-vet go-lint go-test

go-vet: # @HELP examines Go code and reports suspicious constructs
	go vet ./...

go-lint-install: # @HELP installs linters (i.e., 'golint') locally
	go install golang.org/x/lint/golint@latest

go-lint: # @HELP runs linters against the Go codebase
go-lint: go-lint-install
	golint ./...

go-test: # @HELP runs unit tests to test the Go code
	mkdir -p $(PROJECT_DIR)/build
	go test -v ./... -race -coverprofile=$(COVERPROFILE_PATH) -covermode=atomic

go-tidy: # @HELP runs go mod commands (i.e., 'go mod tidy && go mod vendor')
	go mod tidy
	go mod vendor

clean: # @HELP cleans dependencies
	rm -r $(PROJECT_DIR)/build/
	rm -r $(PROJECT_DIR)/vendor/

help:
	@echo "Following 'make' targets can be invoked:"
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
	@echo ""

