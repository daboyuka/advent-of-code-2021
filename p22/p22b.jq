#!/usr/bin/env jq -s -R -f
include "p22/common";

parse |
sumcubes |
map(cubevol) | add
