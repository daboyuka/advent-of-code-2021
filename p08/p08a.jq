#!/usr/bin/env jq -s -R -f
include "./helpers";

def parse:
  lines |
  map(
    [ scan("\\S+") ] |
    index("|") as $barIdx |
    {samp: .[:$barIdx], input: .[$barIdx+1:]}
  )
;

parse |
map(
  .input |
  map(select(length | . == 2 or . == 3 or . == 4 or . == 7)) |
  length
) |
add
