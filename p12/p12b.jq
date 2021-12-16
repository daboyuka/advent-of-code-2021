#!/usr/bin/env jq -s -R -f
include "p12/common";

parse |
countpaths(true)
