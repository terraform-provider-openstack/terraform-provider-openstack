#!/bin/bash
#
set -x
set -o pipefail

timeout="60m"
failed=

LOG_DIR=/tmp/devstack-logs
mkdir -p ${LOG_DIR}

if [[ -z "${ACCEPTANCE_TESTS_FILTER}" ]]; then
    ACCEPTANCE_TESTS=($(python <<< "print(' '.join($ACCEPTANCE_TESTS))"))
else
    ACCEPTANCE_TESTS=$(go test ./openstack/ -v -list 'Acc' | grep -i -E "$ACCEPTANCE_TESTS_FILTER")
    ACCEPTANCE_TESTS=($ACCEPTANCE_TESTS)
fi

if [[ -z $ACCEPTANCE_TESTS ]]; then
    echo "No acceptance tests to run"
    exit 0
fi

# Source credentials as admin
source `dirname $0`/stackenv.sh admin

for acceptance_test in "${ACCEPTANCE_TESTS[@]}"; do
  OS_DEBUG=1 TF_LOG=DEBUG TF_ACC=1 TF_ACC_TERRAFORM_VERSION=1.2.9 go test ./openstack -v -timeout 120m -run $(echo "$acceptance_test" | tr " " "|") |& tee -a ${LOG_DIR}/acceptance_tests.log
  # Check the error code after each suite, but do not exit early if a suite failed.
  if [[ $? != 0 ]]; then
    failed=1
  fi
done

# Source credentials as user (demo)
source `dirname $0`/stackenv.sh demo

for acceptance_test in "${ACCEPTANCE_TESTS[@]}"; do
  OS_DEBUG=1 TF_LOG=DEBUG TF_ACC=1 TF_ACC_TERRAFORM_VERSION=1.2.9 go test ./openstack -v -timeout 120m -run $(echo "$acceptance_test" | tr " " "|") |& tee -a ${LOG_DIR}/acceptance_tests.log
  # Check the error code after each suite, but do not exit early if a suite failed.
  if [[ $? != 0 ]]; then
    failed=1
  fi
done

# If any of the test suites failed, exit 1
if [[ -n $failed ]]; then
  exit 1
fi

exit 0
