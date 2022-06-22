#!/bin/sh

find ./test-data/* -type d -print0 | while read -r -d $'\0' TEST_PATH;do ./build/tgswitch -c "${TEST_PATH}" || exit 1; done;

