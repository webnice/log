DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
PACKETS=$(shell cat .testpackages)

default: lint test

link:
.PHONY: link

## Generate code by go generate or other utilities
generate: link
	# GOPATH=${GOPATH} go generate
	# GOPATH=${GOPATH} easyjson -output_filename gelf/gelf_client_gen.go src/github.com/webnice/lv2/gelf/gelf_client.go
.PHONY: generate

## Dependence managers
dep: link
	# GOPATH=${GOPATH} glide install
.PHONY: dep

test: link
	echo "mode: set" > coverage.log
	for PACKET in $(PACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=coverage.log
.PHONY: cover

bench: link
	GOPATH=${GOPATH} go test -race -bench=. -benchmem ./...
.PHONY: bench

lint: link
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--disable=aligncheck \
	--disable=gotype \
	--disable=structcheck \
	--skip=src/vendor \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	./...
.PHONY: lint

clean:
	rm -rf ${DIR}/src; true
	rm -rf ${DIR}/bin/*; true
	rm -rf ${DIR}/pkg/*; true
	rm -rf ${DIR}/*.log; true
	rm -rf ${DIR}/*.lock; true
.PHONY: clean
