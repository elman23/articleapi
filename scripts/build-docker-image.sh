#!/bin/bash

version="latest"

if [ $# -ne 0 ]; then
    version="$1"
fi

scriptsdir="scripts"

lastdir=$(pwd | sed 's#.*/##')

builddir="."
if [ $lastdir = $scriptsdir ]; then
    builddir=".."
fi

echo "Build dir:" $builddir
echo "Building articleapi version" $version "."

docker build -t articleapi:$version $builddir