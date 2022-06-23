#!/bin/bash

if [[ ! -f "./build/tgswitch" ]]; then
    echo "Cannot run test.."
    echo "Build does not exist"
    exit 1
fi
#find ./test-data/* -type d -print0 | while read -r -d $'\0' TEST_PATH;do ./build/tgswitch -c "${TEST_PATH}" || exit 1; done;

function commentStart(){
    echo "### Test start ###"
}

function commentComplete(){
    echo "### Test complete ###" 
    echo "" 
}

function runtestdir(){
    commentStart
    test_case=$1
    dir=$2
    expected_version=$3
    echo "Test $test_case file"
    ./build/tgswitch -c ./test-data/"$dir" || exit 1
    version=$(terragrunt -v | awk '{print $3}')
    echo "$version"
    sleep 1
    if [[ "$version" == "$expected_version" ]]; then
        echo "Switch successful"
    else
        echo "Switch failed"
        exit 1
    fi
    commentComplete
}


function runtestenv(){
    commentStart
    test_case=$1
    env=$2
    expected_version=$3

    echo "Test $test_case"
    export TG_VERSION="$env"
    ./build/tgswitch || exit 1
    version=$(terragrunt -v | awk '{print $3}')
    echo "$version"
    sleep 1
    if [[ "$version" == "$expected_version" ]]; then
        echo "Switch successful"
    else
        echo "Switch failed"
        exit 1
    fi
    commentComplete
}

function runtestarg(){
    commentStart
    test_case=$1
    arg=$2
    expected_version=$3
    echo "Test $test_case"
    ./build/tgswitch "$arg"|| exit 1
    version=$(terragrunt -v | awk '{print $3}')
    echo "$version"
    sleep 1
    if [[ "$version" == "$expected_version" ]]; then
        echo "Switch successful"
    else
        echo "Switch failed"
        exit 1
    fi
    commentComplete
}

runtestdir "terragrunt version" "test_terragrunt-version" "v0.36.0"
runtestdir "terragrunt hcl" "test_terragrunt_hcl" "v0.36.0"
runtestdir "tgswitchrc" "test_tgswitchrc" "v0.33.0"
runtestdir ".toml" "test_tgswitchtoml" "v0.34.0"
runtestenv "env variable" "0.37.1" "v0.37.1"
runtestarg "passing argument" "0.36.1" "v0.36.1"
