all: binary

docker: binary
	docker build .

vendor:
	dep ensure

binary: vendor
	cd cmd/kubevirt-template-indexer && go build -v .

clean:
	rm -f cmd/kubevirt-template-indexer/kubevirt-template-indexer

.PHONY: all docker binary clean

