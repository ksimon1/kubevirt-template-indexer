FROM fedora:28

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

COPY cmd/kubevirt-template-indexer/kubevirt-template-indexer /usr/sbin/kubevirt-template-indexer
COPY cluster/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
