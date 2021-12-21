#!/usr/bin/env jq -s -R -f
include "p19/common";

parse |
solve |
map(.[0][]) | unique | length
