#!/bin/bash

#################################################
source $(dirname "$0")/env.sh
cd $ROOT
dependencies
#################################################

lints && tests