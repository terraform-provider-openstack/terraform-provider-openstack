#!/bin/bash
#
set -x
set -o pipefail

# set tests timeout to 60m if not set
: "${ACCEPTANCE_TESTS_TIMEOUT:=60m}"
# set tests parallelism to 1 if not set
: "${ACCEPTANCE_TESTS_PARALLELISM:=1}"

failed=

LOG_DIR=/tmp/devstack-logs
mkdir -p "${LOG_DIR}"

if [[ -n ${ACCEPTANCE_TESTS_FILTER} ]]; then
    ACCEPTANCE_TESTS="(?i)(?:${ACCEPTANCE_TESTS_FILTER})"
fi

if [[ -z ${ACCEPTANCE_TESTS} ]]; then
    echo "No acceptance tests to run"
    exit 0
fi

# Source credentials as admin
source `dirname $0`/stackenv.sh admin

export OS_DEBUG=1
export TF_LOG=DEBUG
export TF_ACC=1

go test ./openstack -v -timeout "${ACCEPTANCE_TESTS_TIMEOUT}" -parallel "${ACCEPTANCE_TESTS_PARALLELISM}" -run "${ACCEPTANCE_TESTS}" |& tee -a "${LOG_DIR}/acceptance_tests.log"
# Check the error code, but do not exit early.
if [[ $? != 0 ]]; then
  failed=1
fi

# Source credentials as user (demo)
source `dirname $0`/stackenv.sh demo

go test ./openstack -v -timeout "${ACCEPTANCE_TESTS_TIMEOUT}" -parallel "${ACCEPTANCE_TESTS_PARALLELISM}" -run "${ACCEPTANCE_TESTS}" |& tee -a "${LOG_DIR}/acceptance_tests.log"
# Check the error code, but do not exit early.
if [[ $? != 0 ]]; then
  failed=1
fi

# If any of the test suites failed, exit 1
if [[ -n $failed ]]; then
  exit 1
fi

# Check if there were no tests to run
if grep -q 'no tests to run' "${LOG_DIR}/acceptance_tests.log"; then
  exit 1
fi

exit 0
