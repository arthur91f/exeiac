#!/bin/bash

echo "INFO: we use module test for the moment" >&2
"$(dirname "$0")/module_test.sh" "$@"

