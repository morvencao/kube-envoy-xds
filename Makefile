.DEFAULT_GOAL	:= build

#------------------------------------------------------------------------------
# Variables
#------------------------------------------------------------------------------

SHELL 	:= /bin/bash
BINDIR	:= bin



#-----------------
#-- code generaion
#-----------------

generate: $(BINDIR)/gogofast $(BINDIR)/validate
	@echo "--> generating pb.go files"
	$(SHELL) generate_proto.sh

#------------------
#-- dependencies
#------------------
.PHONY: depend.update depend.install

depend.update: tools.glide
	@echo "--> updating dependencies from glide.yaml"
	@glide update

depend.install: tools.glide
	@echo "--> installing dependencies from glide.lock "
	@glide install

vendor:
	@echo "--> installing dependencies from glide.lock "
	@glide install

$(BINDIR):
	@mkdir -p $(BINDIR)

#---------------
#-- tools
#---------------
.PHONY: tools tools.glide tools.goimports tools.golint tools.govet

tools: tools.glide tools.goimports tools.golint tools.govet

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing goimports"; \
		go get golang.org/x/tools/cmd/goimports; \
	fi

tools.govet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		echo "--> installing govet"; \
		go get golang.org/x/tools/cmd/vet; \
	fi

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get -u golang.org/x/lint/golint; \
	fi

tools.glide:
	@command -v glide >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing glide"; \
		curl https://glide.sh/get | sh; \
	fi

$(BINDIR)/gogofast: vendor
	@echo "--> building $@"
	@go build -o $@ vendor/github.com/gogo/protobuf/protoc-gen-gogofast/main.go

$(BINDIR)/validate: vendor
	@echo "--> building $@"
	@go build -o $@ vendor/github.com/lyft/protoc-gen-validate/main.go
