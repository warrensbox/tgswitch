#!/bin/bash

if [[ ! -f "./build/tgswitch" ]]; then
    echo "Cannot run test.."
    echo "Build does not exist"
    exit 1
fi
#find ./test-data/* -type d -print0 | while read -r -d $'\0' TEST_PATH;do ./build/tgswitch -c "${TEST_PATH}" || exit 1; done;

function runtestdir(){
    test_case=$1
    dir=$2
    expected_version=$3
    echo "Test $test_case file"
    ./build/tgswitch -c ./test-data/"$dir" || exit 1
    version=$(terragrunt -v | awk '{print $3}')
    echo "$version"
    if [[ "$version" == "$expected_version" ]]; then
        echo "Switch successful"
    else
        echo "Switch failed"
        exit 1
    fi
}


function runtestenv(){
    test_case=$1
    env=$2
    expected_version=$3

    echo "Test $test_case"
    export TG_VERSION="$env"
    ./build/tgswitch || exit 1
    version=$(terragrunt -v | awk '{print $3}')
    echo "$version"
    if [[ "$version" == "$expected_version" ]]; then
        echo "Switch successful"
    else
        echo "Switch failed"
        exit 1
    fi
}

runtestdir "terragrunt version" "test_terragrunt-version" "v0.36.0"
runtestdir "terragrunt hcl" "test_terragrunt_hcl" "v0.37.4"
runtestdir "tfswitchrc" "test_tgswitchrc" "v0.33.0"
runtestdir ".toml" "test_tfswitchtoml" "v0.34.0"
runtestenv "env variable" "0.37.1" "v0.37.1"

