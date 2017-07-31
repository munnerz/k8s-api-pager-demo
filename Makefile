PACKAGE_NAME := github.com/munnerz/k8s-api-pager-demo

# A temporary directory to store generator executors in
BINDIR ?= bin
GOPATH ?= $HOME/go
HACK_DIR ?= hack

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
generate: .generate_exes $(TYPES_FILES)
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
