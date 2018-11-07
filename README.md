kubevirt-template-indexer
=========================

Template index service for [kubevirt](http://kubevirt.io) built using [the controller-runtime project](https://github.com/kubernetes-sigs/controller-runtime)

License: APACHE v2
Copyright: Red Hat Inc

Real documentation coming soon.

TODO
----
- fix this README
- code docs
- routes package uses globals
- error responses are not handled (and unspecified)
- expose the service outside the cluster?
- /template filtering is untested
- lack of unit tests for "templateindex" package
- check why sometimes the sync doesn't happen (just timing?)
- functional tests?
