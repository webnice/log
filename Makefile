DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)

default: lint test

test:
	clear
	mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. log.v2; true
	echo "mode: set" > coverage.log
	for pkg in `cat .testpackages`; do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log -coverpkg=$$pkg; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done

.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=coverage.log
.PHONY: test

lint:
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=15 \
	--disable=aligncheck \
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
