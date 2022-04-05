# Copyright (c) 2019 StackRox Inc.
#
# Modifications Copyright (c) 2019 Elisa Oyj
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Makefile for building the RunAsUser Admission Controller
OPERATOR_NAME := runasuser-admission-controller
IMAGE ?= elisaoyj/$(OPERATOR_NAME)
ifeq ($(USE_JSON_OUTPUT), 1)
GOTEST_REPORT_FORMAT := -json
endif

.PHONY: clean deps test gofmt ensure run build build-image build-linux-amd64

clean:
	git clean -Xdf

deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint

test:
	GO111MODULE=on go test ./cmd -v -coverprofile=gotest-coverage.out $(GOTEST_REPORT_FORMAT) > gotest-report.out && cat gotest-report.out || (cat gotest-report.out; exit 1)
	GO111MODULE=off golint -set_exit_status ./cmd  > golint-report.out && cat golint-report.out || (cat golint-report.out; exit 1)
	GO111MODULE=on go vet -mod vendor ./cmd
	./hack/gofmt.sh
	git diff --exit-code go.mod go.sum

gofmt:
	./hack/gofmt.sh

ensure:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

run: build
	./bin/$(OPERATOR_NAME)

build-linux-amd64:
	rm -rf bin/linux/$(OPERATOR_NAME)
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o bin/linux/$(OPERATOR_NAME) ./cmd

build:
	rm -f bin/$(OPERATOR_NAME)
	GO111MODULE=on go build -v -i -o bin/$(OPERATOR_NAME) ./cmd

build-image: build-linux-amd64
	docker build -t $(IMAGE):latest .
