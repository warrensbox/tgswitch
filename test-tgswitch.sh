#!/bin/sh

if [[ ! -f "./build/tgswitch" ]]; then
    echo "Cannot run test.."
    echo "Build does not exist"
    exit 1
fi
#find ./test-data/* -type d -print0 | while read -r -d $'\0' TEST_PATH;do ./build/tgswitch -c "${TEST_PATH}" || exit 1; done;

echo "Test terragrunt-version file"
./build/tgswitch -c ./test-data/test_terragrunt-version


echo "Test terragrunt hcl file"


echo "Test tfswitchrc file"


echo "Test tfswitchtoml file"

