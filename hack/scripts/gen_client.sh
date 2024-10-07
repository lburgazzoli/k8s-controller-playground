#!/usr/bin/env bash

if [ $# -ne 1 ]; then
    echo "project root is expected"
fi

PROJECT_ROOT="$1"
TMP_DIR=$( mktemp -d -t dapr-client-gen-XXXXXXXX )

mkdir -p "${TMP_DIR}/client"
mkdir -p "${PROJECT_ROOT}/pkg/client"

echo "tmp dir: $TMP_DIR"

echo "Generating openapi schema"
go run k8s.io/kube-openapi/cmd/openapi-gen \
  --output-file zz_generated.openapi.go \
  --output-dir "pkg/generated/openapi" \
  --output-pkg "github.com/lburgazzoli/k8s-controller-playground/pkg/generated/openapi" \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1 \
  k8s.io/apimachinery/pkg/apis/meta/v1 \
  k8s.io/apimachinery/pkg/runtime \
  k8s.io/apimachinery/pkg/version

echo "Generate ApplyConfiguration"
go run k8s.io/code-generator/cmd/applyconfiguration-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --openapi-schema <(go run ${PROJECT_ROOT}/cmd/main.go modelschema) \
  --output-dir="${TMP_DIR}/client/applyconfiguration" \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

echo "Generate client"
go run k8s.io/code-generator/cmd/client-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/clientset" \
  --input-base=github.com/lburgazzoli/k8s-controller-playground/api \
  --input=playground/v1alpha1 \
  --fake-clientset=false \
  --clientset-name "versioned"  \
  --apply-configuration-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset

echo "Generate lister"
go run k8s.io/code-generator/cmd/lister-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/listers" \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/listers \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

echo "Generate informer"
go run k8s.io/code-generator/cmd/informer-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/informers" \
  --versioned-clientset-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset/versioned \
  --listers-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/listers \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/informers \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

#
# See: https://github.com/kubernetes/code-generator/issues/150
#sed -i \
#  's/WithAPIVersion(\"playground\/v1alpha1\")/WithAPIVersion(\"playground.lburgazzoli.github.io\/v1alpha1\")/g' \
#  "${TMP_DIR}"/client/applyconfiguration/playground/v1alpha1/component.go

cp -r \
  "${TMP_DIR}"/client/* \
  "${PROJECT_ROOT}"/pkg/client

