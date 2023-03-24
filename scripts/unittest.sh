#!/bin/bash
#
set -x
set -o pipefail

failed=


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