#!/usr/bin/env jq -s -R -f
include "p19/common";

parse |
solve |
map(.[1][1]) |  # array of translations ([scannerIdx][dimIdx])
. as $alltrans |
map(
  subpos3($alltrans[]) |
  map(fabs) |
  add
) | max
