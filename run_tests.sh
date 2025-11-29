#!/bin/bash
COVERAGE_FILE="coverage.out"

SPEC_FILE=".coverage-ignore.yaml"

echo "Running tests and creating raw coverage report..."
go test -v -coverprofile=coverage.out -covermode=atomic ./...

echo "Filtering coverage report using $SPEC_FILE..."
GO_COVER_IGNORE_SPEC_PATH="$SPEC_FILE" \
GO_COVER_IGNORE_COVER_PROFILE_PATH="$COVERAGE_FILE" \
go-cover-ignore

echo "Checking threshold on the filtered report..."
go-test-coverage --config=.testcoverage.yml --profile=$COVERAGE_FILE
