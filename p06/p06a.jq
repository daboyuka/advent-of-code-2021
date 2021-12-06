#!/usr/bin/env jq -s -R -f
include "./p06/common";

parse |
nevolve(80) |
add
