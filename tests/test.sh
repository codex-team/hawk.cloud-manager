#!/usr/bin/env bash
set -xeuo pipefail

docker-compose -f tests/docker-compose-test.yaml up --build -d
test_status_code=0
docker-compose -f tests/docker-compose-test.yaml logs > logs.txt ;\
docker-compose -f tests/docker-compose-test.yaml run integration_tests go test || test_status_code=$? ;\
docker-compose -f tests/docker-compose-test.yaml down ;\
if grep -q fail logs.txt; then
  cat logs.txt | grep fail
  exit 1
fi
exit $test_status_code
