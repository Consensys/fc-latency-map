#!/bin/bash

trgt_coverage=${1:-"80"}

go test ./... -coverprofile cover.out.tmp -covermode=atomic
cat cover.out.tmp | grep -v "mocks.go" > coverage.txt
rm cover.out.tmp
curr_coverage=$(go tool cover -func coverage.txt | grep total | awk '{print $3}' | tr -d '%')

bash <(curl -s https://codecov.io/bash)

echo "Total: $curr_coverage%"

if awk "BEGIN {exit !("$curr_coverage" >= "$trgt_coverage")}"; then
  echo "Unit tests are passing $trgt_coverage% coverage"
  exit 0
else
  echo "Unit tests are not passing $trgt_coverage% coverage"
  exit 1
fi
