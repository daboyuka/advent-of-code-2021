#!/usr/bin/env jq -s -R -f
include "p14/common";

parse |
fuserules |
countpairs |
reduce range(10) as $_ ( . ; step ) |
countelems |
map(.) | max - min
