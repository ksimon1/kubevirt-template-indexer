all: template-indexer.yaml


template-indexer.yaml: template-indexer-deployment.yaml template-indexer-rbac.yaml template-indexer-service.yaml
	cat template-indexer-rbac.yaml template-indexer-deployment.yaml template-indexer-service.yaml > template-indexer.yaml

clean:
	rm template-indexer.yaml

.PHONY: clean
