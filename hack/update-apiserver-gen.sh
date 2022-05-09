#!/bin/bash

# Copyright (c) 2022 Red Hat, Inc.


# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")"/..)
BINDIR="${REPO_ROOT}"/_output
SC_PKG='github.com/clyang82/hub-of-hubs-apis'

# generate openapi for servicecatalog and settings group
"${BINDIR}"/openapi-gen "$@" \
	--v 1 --logtostderr \
	--go-header-file "${REPO_ROOT}"/hack/custom-boilerplate.go.txt \
	--input-dirs "open-cluster-management.io/governance-policy-propagator/api/v1" \
	--output-package "${SC_PKG}/pkg/server/apis/openapi" \
  --report-filename ".api_violation.report"
