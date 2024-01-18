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

# Prepare environment and add env vars to openrc
source `dirname $0`/stackenv.sh

# Source credentials as admin
source $DEVSTACK_PATH/openrc admin admin

for acceptance_test in "${ACCEPTANCE_TESTS[@]}"; do
  OS_DEBUG=1 TF_LOG=DEBUG TF_ACC=1 go test ./openstack -v -timeout 120m -run $(echo "$acceptance_test" | tr " " "|") |& tee -a ${LOG_DIR}/acceptance_tests.log
  # Check the error code after each suite, but do not exit early if a suite failed.
  if [[ $? != 0 ]]; then
    failed=1
  fi
done

# Source credentials as user (demo)
source $DEVSTACK_PATH/openrc demo demo

for acceptance_test in "${ACCEPTANCE_TESTS[@]}"; do
  OS_DEBUG=1 TF_LOG=DEBUG TF_ACC=1 go test ./openstack -v -timeout 120m -run $(echo "$acceptance_test" | tr " " "|") |& tee -a ${LOG_DIR}/acceptance_tests.log
  # Check the error code after each suite, but do not exit early if a suite failed.
  if [[ $? != 0 ]]; then
    failed=1
  fi
done

# Source credentials as admin and enable system scope
source $DEVSTACK_PATH/openrc admin admin
export OS_SYSTEM_SCOPE=true
unset OS_PROJECT_NAME
unset OS_TENANT_NAME
unset OS_PROJECT_DOMAIN_ID


for acceptance_test in "${ACCEPTANCE_TESTS[@]}"; do
  OS_DEBUG=1 TF_LOG=DEBUG TF_ACC=1 go test ./openstack -v -timeout 120m -run $(echo "$acceptance_test" | tr " " "|") |& tee -a ${LOG_DIR}/acceptance_tests.log
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
