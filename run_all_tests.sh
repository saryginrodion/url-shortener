#!/bin/sh

set -euo pipefail

echo "Starting unit tests"
./run_unit_tests.sh

echo "Starting integration tests"
./run_integration_tests.sh
