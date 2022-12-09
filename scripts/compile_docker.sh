#! /bin/bash

set -e

cd _ci/build/coding-challenge

docker build -t coding-challenge .
