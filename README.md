kubevirt-template-indexer
=========================

Template index service for [kubevirt](http://kubevirt.io) built using [the controller-runtime project](https://github.com/kubernetes-sigs/controller-runtime)

License: APACHE v2
Copyright: Red Hat Inc

Description
-----------

TODO

Usage
-----

TODO

Build
-----

The Makefile automates the build process. To set up the dependencies, run
```
make vendor
```

To build the binary, run
```
make binary
```
This stage also automatically takes care of set up the dependencies, so you don't _need_ to do it explicitely

To build the docker image, run
```
make docker
```

Should you want to remove the binary, just run
```
make clean
```

Contribute
----------
Just fork this repo and send a PR. If you are sending in a code change, make sure your contribution is covered by some automated test (either existing ones or new ones)


Run it in a Kubernetes cluster
------------------------------

To run the server in a kubernetes cluster, just run
```
kubectl create -f template-indexer.yaml
```

the `template-indexer.yaml` file is an amalgamation of the manifests which define the account/rbac settings, the deployment and the service.
If you need for whatever reason to (re)create it, do
```
make -C cluster
```

If you want to install the server step by step:
first, set up the accounts and RBAC:
```
kubectl create -f cluster/template-indexer-rbac.yaml
```

Then deploy the server:
```
kubectl create -f cluster/template-indexer-deployment.yaml
```

Last, create the service:
```
kubectl create -f cluster/template-indexer-service.yaml
```

Run it outside a Kubernetes cluster
-----------------------------------

Partially implemented and not suppored yet.


TODO
----
- code docs
- routes package uses globals
- error responses are not handled (and unspecified)
- the /template filtering is not properly tested
- check why sometimes the sync doesn't happen (just timing?)
- functional tests?
- integration tests
