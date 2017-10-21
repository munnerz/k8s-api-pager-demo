PACKAGE_NAME := github.com/srossross/k8s-test-controller

# A temporary directory to store generator executors in
BINDIR ?= bin
GOPATH ?= $HOME/go
HACK_DIR ?= hack

GOOS := $(shell go env GOHOSTOS)
GOARCH := $(shell go env GOHOSTARCH)
CGO_ENABLED := 0
LDFLAGS := -X github.com/srossross/k8s-test-controller/main.VERSION=$(shell echo $${CIRCLE_TAG:-?}) \
  -X github.com/srossross/k8s-test-controller/main.BUILD_TIME=$(shell date -u +%Y-%m-%d)

USERNAME := $(shell echo ${CIRCLE_PROJECT_USERNAME})
REPONAME := $(shell echo ${CIRCLE_PROJECT_REPONAME})

# A list of all types.go files in pkg/apis
TYPES_FILES = $(shell find pkg/apis -name types.go)

# This step pulls the Kubernetes repo so we can build generators in another
# target. Soon, github.com/kubernetes/kube-gen will be live meaning we don't
# need to pull the entirety of the k8s source code.
.get_deps:
	@echo "Grabbing dependencies..."
	@go get -d -u k8s.io/kubernetes/ || true
	@go get -d github.com/kubernetes/repo-infra || true
	# Once k8s.io/kube-gen is live, we should be able to remove this dependency
	# on k8s.io/kubernetes. https://github.com/kubernetes/kubernetes/pull/49114
	cd ${GOPATH}/src/k8s.io/kubernetes; git checkout 25d3523359ff17dda6deb867a7c3dd6c8b7ea705;
	@touch $@

# Targets for building k8s code generators
#################################################
.generate_exes: .get_deps \
				$(BINDIR)/defaulter-gen \
                $(BINDIR)/deepcopy-gen \
                $(BINDIR)/conversion-gen \
                $(BINDIR)/client-gen \
                $(BINDIR)/lister-gen \
                $(BINDIR)/informer-gen
	touch $@

$(BINDIR)/defaulter-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/defaulter-gen

$(BINDIR)/deepcopy-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen

$(BINDIR)/conversion-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/conversion-gen

$(BINDIR)/client-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/client-gen

$(BINDIR)/lister-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/lister-gen

$(BINDIR)/informer-gen:
	go build -o $@ k8s.io/kubernetes/cmd/libs/go2idl/informer-gen
#################################################


# This target runs all required generators against our API types.
generate: .generate_exes $(TYPES_FILES) ## Generate files
	# Generate defaults
	$(BINDIR)/defaulter-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager/v1alpha1" \
		--extra-peer-dirs "$(PACKAGE_NAME)/pkg/apis/pager" \
		--extra-peer-dirs "$(PACKAGE_NAME)/pkg/apis/pager/v1alpha1" \
		--output-file-base "zz_generated.defaults"
	# Generate deep copies
	$(BINDIR)/deepcopy-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager/v1alpha1" \
		--output-file-base zz_generated.deepcopy
	# Generate conversions
	$(BINDIR)/conversion-gen \
		--v 1 --logtostderr \
		--go-header-file "$${GOPATH}/src/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager" \
		--input-dirs "$(PACKAGE_NAME)/pkg/apis/pager/v1alpha1" \
		--output-file-base zz_generated.conversion
	# generate all pkg/client contents
	$(HACK_DIR)/update-client-gen.sh

cacheBuilds: ## Make go build and go run faster
	go list -f '{{.Deps}}' ./... | tr "[" " " | tr "]" " " |   xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' |   xargs go install -a


build: ## build for any arch
	mkdir -p /tmp/commands

	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o ./k8s-test-controller-$(GOOS)-$(GOARCH) ./main.go
	tar -zcvf /tmp/commands/k8s-test-controller-$(GOOS)-$(GOARCH).tgz ./k8s-test-controller-$(GOOS)-$(GOARCH)

buildLinux: GOOS := linux
buildLinux: GOARCH := amd64
buildLinux: build

dockerBuild: ## Build docker container
	docker build -t srossross/k8s-test-controller:latest .

release: ## Create github release
	github-release release \
		--user $(USERNAME) \
		--repo $(REPONAME) \
		--tag $(TAG) \
		--name "Release $(TAG)" \
		--description "TODO: Description"

upload: ## Upload build artifacts to github

	github-release upload \
		--user $(USERNAME) \
		--repo $(REPONAME) \
		--tag $(TAG) \
		--name "k8s-test-controller-linux-amd64." \
		--file /tmp/commands/k8s-test-controller-linux-amd64.tgz



.PHONY: help

help: ## show this help and exit
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
