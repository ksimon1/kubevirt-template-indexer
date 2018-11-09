kubevirt-template-indexer
=========================

Template index service for [kubevirt](http://kubevirt.io) built using [the controller-runtime project](https://github.com/kubernetes-sigs/controller-runtime)

License: APACHE v2
Copyright: Red Hat Inc

Description
-----------

`kubevirt-template-indexer` provides indexing services for the [kubevirt VM templates](https://github.com/kubevirt/common-templates/) through HTTP endpoints.
A client can connect to the server to get a filtered view of the templates deployed in the cluster, and to have summaries about them.

Usage
-----

The server exposes four HTTP endpoints, providing answers in JSON.

`/oses` returns a collection of all the OS of the templates deployed in the cluster. Example response:
```json
[
    {
        "id": "centos7.0",
        "name": ""
    },
    {
        "id": "opensuse15.0",
        "name": ""
    },
    {
        "id": "fedora28",
        "name": ""
    }
]
```

`/workloads` returns a collection of all the Workloads of the templates. Example response:
```json
[
    {
        "id": "generic",
        "name": ""
    },
    {
        "id": "highperformance",
        "name": ""
    }
]
```

`/size` returns a collection of all the size (flavors) of the templates. Example response:
```json
[
    {
        "id": "tiny",
        "name": ""
    },
    {
        "id": "medium",
        "name": ""
    },
    {
        "id": "large",
        "name": ""
    },
    {
        "id": "small",
        "name": ""
    }
]
```

Notice that the "name" field of the returned values above is empty. This is because the server cannot figure out where to look for
the description of the summarized itemd. You can supply configuration files to fill this information.
A configuration file must be a JSON file sitting in the directory specified by the `-C` flag. You need a file for each `os`, `workload`, `size`.
For example, to set the names for the `size`s, using the defaults, you can drop the file under `/etc/template-index/size.json`

The content should be a JSON map whose keys are the known sizes, and whose values are the names you want to set. Look under `examples/` for examples.

`/templates` return a summary of all the templates. Example response:
```json
[
    {
        "id": "centos7-generic-large",
        "name": "CentOS 7.0+ VM",
        "description": "This template can be used to create a VM suitable for CentOS 7 and newer. The template assumes that a PVC is available which is providing the necessary CentOS disk image.",
        "icon-id": "icon-centos",
        "osid": "centos7.0",
        "workload": "generic",
        "size": "large"
    },
    {
        "id": "centos7-generic-small",
        "name": "CentOS 7.0+ VM",
        "description": "This template can be used to create a VM suitable for CentOS 7 and newer. The template assumes that a PVC is available which is providing the necessary CentOS disk image.",
        "icon-id": "icon-centos",
        "osid": "centos7.0",
        "workload": "generic",
        "size": "small"
    },
    {
        "id": "centos7-generic-medium",
        "name": "CentOS 7.0+ VM",
        "description": "This template can be used to create a VM suitable for CentOS 7 and newer. The template assumes that a PVC is available which is providing the necessary CentOS disk image.",
        "icon-id": "icon-centos",
        "osid": "centos7.0",
        "workload": "generic",
        "size": "medium"
    },
    {
        "id": "centos7-generic-tiny",
        "name": "CentOS 7.0+ VM",
        "description": "This template can be used to create a VM suitable for CentOS 7 and newer. The template assumes that a PVC is available which is providing the necessary CentOS disk image.",
        "icon-id": "icon-centos",
        "osid": "centos7.0",
        "workload": "generic",
        "size": "tiny"
    }
]
```

You can filter the output using the query parameters. Example:
```json
[
    {
        "id": "fedora-highperformance-large",
        "name": "Fedora 23+ VM",
        "description": "This template can be used to create a VM suitable for Fedora 23 and newer. The template assumes that a PVC is available which is providing the necessary Fedora disk image.\nRecommended disk image (needs to be converted to raw) https://download.fedoraproject.org/pub/fedora/linux/releases/28/Cloud/x86_64/images/Fedora-Cloud-Base-28-1.1.x86_64.qcow2",
        "icon-id": "icon-fedora",
        "osid": "fedora28",
        "workload": "highperformance",
        "size": "large"
    },
    {
        "id": "fedora-highperformance-small",
        "name": "Fedora 23+ VM",
        "description": "This template can be used to create a VM suitable for Fedora 23 and newer. The template assumes that a PVC is available which is providing the necessary Fedora disk image.\nRecommended disk image (needs to be converted to raw) https://download.fedoraproject.org/pub/fedora/linux/releases/28/Cloud/x86_64/images/Fedora-Cloud-Base-28-1.1.x86_64.qcow2",
        "icon-id": "icon-fedora",
        "osid": "fedora28",
        "workload": "highperformance",
        "size": "small"
    },
    {
        "id": "fedora-highperformance-tiny",
        "name": "Fedora 23+ VM",
        "description": "This template can be used to create a VM suitable for Fedora 23 and newer. The template assumes that a PVC is available which is providing the necessary Fedora disk image.\nRecommended disk image (needs to be converted to raw) https://download.fedoraproject.org/pub/fedora/linux/releases/28/Cloud/x86_64/images/Fedora-Cloud-Base-28-1.1.x86_64.qcow2",
        "icon-id": "icon-fedora",
        "osid": "fedora28",
        "workload": "highperformance",
        "size": "tiny"
    },
    {
        "id": "fedora-highperformance-medium",
        "name": "Fedora 23+ VM",
        "description": "This template can be used to create a VM suitable for Fedora 23 and newer. The template assumes that a PVC is available which is providing the necessary Fedora disk image.\nRecommended disk image (needs to be converted to raw) https://download.fedoraproject.org/pub/fedora/linux/releases/28/Cloud/x86_64/images/Fedora-Cloud-Base-28-1.1.x86_64.qcow2",
        "icon-id": "icon-fedora",
        "osid": "fedora28",
        "workload": "highperformance",
        "size": "medium"
    }
]
```
Possible filter parameters are `size`, `os`, `workload`.


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

the `template-indexer.yaml` file is an amalgamation of the manifests which define the account/RBAC settings, the deployment and the service.
If you need for whatever reason to (re)create it, do
```
make -C cluster
```

If you want to install the server step by step:
first, set up the accounts and RBAC:
```
kubectl create -f cluster/template-indexer-RBAC.yaml
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

Partially implemented and not supported yet.


TODO
----
- code docs
- routes package uses globals
- error responses are not handled (and unspecified)
- functional tests?
- integration tests
