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
NAME := runasuser-admission-controller
IMAGE ?= elisaoyj/$(NAME)
.PHONY: prepare-build build-image build-linux-amd64

clean:
	git clean -Xdf

prepare-build:
	rm -rf bin/

build-linux-amd64: prepare-build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -i -o bin/linux/$(NAME) ./cmd/$(NAME)

build-image:
	docker build -t $(IMAGE):latest .
