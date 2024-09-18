#!/bin/sh

if [ $# -ne 1 ]; then
    echo "project root is expected"
fi

PROJECT_ROOT="$1"
TMP_DIR=$( mktemp -d -t dapr-client-gen-XXXXXXXX )

mkdir -p "${TMP_DIR}/client"
mkdir -p "${PROJECT_ROOT}/pkg/client"

echo "tmp dir: $TMP_DIR"

echo "applyconfiguration-gen"
"${PROJECT_ROOT}"/bin/applyconfiguration-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/applyconfiguration" \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

echo "client-gen"
"${PROJECT_ROOT}"/bin/client-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/clientset" \
  --input-base=github.com/lburgazzoli/k8s-controller-playground/api \
  --input=playground/v1alpha1 \
  --fake-clientset=false \
  --clientset-name "versioned"  \
  --apply-configuration-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset

echo "lister-gen"
"${PROJECT_ROOT}"/bin/lister-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/listers" \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/listers \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

echo "informer-gen"
"${PROJECT_ROOT}"/bin/informer-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/informers" \
  --versioned-clientset-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset/versioned \
  --listers-package=github.com/lburgazzoli/k8s-controller-playground/pkg/client/listers \
  --output-pkg=github.com/lburgazzoli/k8s-controller-playground/pkg/client/informers \
  github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1

#
# See: https://github.com/kubernetes/code-generator/issues/150
sed -i \
  's/WithAPIVersion(\"playground\/v1alpha1\")/WithAPIVersion(\"playground.lburgazzoli.github.io\/v1alpha1\")/g' \
  "${TMP_DIR}"/client/applyconfiguration/playground/v1alpha1/component.go

cp -r \
  "${TMP_DIR}"/client/* \
  "${PROJECT_ROOT}"/pkg/client

