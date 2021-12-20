#!/usr/bin/env jq -s -R -f
include "p17/common";

parse as $box |
[ simulateall($box) | select(.hit).peak ] |
max
