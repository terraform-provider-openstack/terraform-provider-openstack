#!/bin/bash
#
set -x
set -o pipefail

failed=

# Ensure that workflows execute proper functional tests
# for the OpenStack provider.

TESTS_FILTER=$(grep FILTER .github/workflows/functional-*.yml | awk -F'"' '{print $(NF-1)}' | paste -sd "|")
EGREP_SKIP="terraform-provider-openstack|database|loadbalancer|TestUnit|TestAccProvider|Taas"

DIFF="$(diff -u <(go test ./openstack/ -list "(?i)(?:${TESTS_FILTER})" | egrep -vi "${EGREP_SKIP}") <(go test ./openstack/ -list "Test" | egrep -vi "${EGREP_SKIP}"))"

if [[ -n $DIFF ]]; then
  echo "The following tests are not covered by the functional tests:"
  echo "$DIFF"
  echo "Please update the functional test names to cover these tests."
  exit 1
fi

UNIT_TESTS=$(go test ./openstack/ -v -list 'Unit' | grep -i "Unit")
UNIT_TESTS=($UNIT_TESTS)


if [[ -z $UNIT_TESTS ]]; then
    echo "No unit tests to run"
    exit 0
fi


for unit_test in "${UNIT_TESTS[@]}"; do
  go test ./openstack -v -count=5 -run $(echo "$unit_test" | tr " " "|")
  # Check the error code after each suite, but do not exit early if a suite failed.
  if [[ $? != 0 ]]; then
    failed=1
  fi
done

# Run tests under openstack/internal
UNIT_TESTS=$(go test ./openstack/internal/pathorcontents -v -list 'Unit' | grep -i "Unit")
UNIT_TESTS=($UNIT_TESTS)

for unit_test in "${UNIT_TESTS[@]}"; do
  go test ./openstack/internal/pathorcontents  -v -count=5 -run $(echo "$unit_test" | tr " " "|")
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
