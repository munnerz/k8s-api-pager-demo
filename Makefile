PACKAGE_NAME := github.com/munnerz/k8s-api-pager-demo
REGISTRY := quay.io/munnerz
IMAGE_NAME := pager
IMAGE_TAGS := latest
GOPATH ?= $HOME/go
HACK_DIR ?= hack
BUILD_TAG := build

# Get a list of all binaries to be built
CMDS := $(shell find ./cmd/ -type d -depth 1 -exec basename {} \; )
# Path to dockerfiles directory
DOCKERFILES := $(HACK_DIR)/build/dockerfiles
# A list of all types.go files in pkg/apis
TYPES_FILES = $(shell find pkg/apis -name types.go)
# docker_build_controller, docker_build_apiserver etc
DOCKER_BUILD_TARGETS = $(addprefix docker_build_, $(CMDS))
# docker_push_controller, docker_push_apiserver etc
DOCKER_PUSH_TARGETS = $(addprefix docker_push_, $(CMDS))


# Alias targets
###############

build: $(CMDS) docker_build
docker_build: $(DOCKER_BUILD_TARGETS)
docker_push: $(DOCKER_PUSH_TARGETS)
push: build docker_push

# Code generation
#################
# This target runs all required generators against our API types.
generate: $(TYPES_FILES)
	$(HACK_DIR)/update-codegen.sh

# Go targets
#################
go_verify: go_fmt go_test go_build

$(CMDS):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-w' -o $(DOCKERFILES)/pager-$@_linux_amd64 ./cmd/$@

go_test:
	go test $$(go list ./... | grep -v '/vendor/')

go_fmt:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi

# Docker targets
################
$(DOCKER_BUILD_TARGETS):
	$(eval DOCKER_BUILD_CMD := $(subst docker_build_,,$@))
	docker build -t $(REGISTRY)/$(IMAGE_NAME)-$(DOCKER_BUILD_CMD):$(BUILD_TAG) -f $(DOCKERFILES)/$(DOCKER_BUILD_CMD)/Dockerfile $(DOCKERFILES)
docker_build: $(DOCKER_BUILD_TARGETS)

$(DOCKER_PUSH_TARGETS):
	$(eval DOCKER_PUSH_CMD := $(subst docker_push_,,$@))
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker tag $(REGISTRY)/$(IMAGE_NAME)-$(DOCKER_PUSH_CMD):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME)-$(DOCKER_PUSH_CMD):$${tag} ; \
		docker push $(REGISTRY)/$(IMAGE_NAME)-$(DOCKER_PUSH_CMD):$${tag}; \
	done
