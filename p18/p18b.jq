#!/usr/bin/env jq -s -R -f
include "p18/common";

parse |
[
  range(length) as $i |
  range(length - $i - 1) as $j |
  [ .[$i], .[$j] ] |
  padd |
  mag
] |
max
