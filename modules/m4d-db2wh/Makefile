.PHONY: all
all: helm-all

ROOT_DIR := ../..
DOCKER_NAME := m4d-db2wh
HELM_VALUES := \
	--values values.yaml.sample \
	--set image=${DOCKER_HOSTNAME}/${DOCKER_NAMESPACE}/mover:${DOCKER_TAGNAME}

include $(ROOT_DIR)/Makefile.env
include $(ROOT_DIR)/hack/make-rules/docker.mk
include $(ROOT_DIR)/hack/make-rules/helm.mk
include $(ROOT_DIR)/hack/make-rules/tools.mk

helm-crd-hook: crd-batch
.PHONY: crd-batch
crd-batch:
	$(ABSTOOLBIN)/kubectl apply -f $(ROOT_DIR)/manager/config/crd/bases/motion.m4d.ibm.com_batchtransfers.yaml
