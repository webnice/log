DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
PACKETS=$(shell cat .testpackages)

default: lint test

generate:
	#GOPATH=${GOPATH} go generate
	#GOPATH=${GOPATH} easyjson -output_filename configuration.go src/gopkg.in/webnice/web.v1/types.go
.PHONY: generate

test:
	clear
	mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. log.v2; true
	echo "mode: set" > coverage.log
	for PACKET in $(PACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log $$PACKET; \
		if [ ! $$? == 0 ]; then exit $$?; fi; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=coverage.log
.PHONY: cover

lint:
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
.PHONY: clean
