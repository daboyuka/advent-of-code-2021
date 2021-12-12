#!/usr/bin/env jq -s -R -f
include "./helpers";

def parse: lines | map(split("")|map(tonumber));

def scanpoints:  # input: grid, output: all coordinates
  ((range(length)|[.]) + (range(.[0]|length)|[.]))
;

def inbounds($p):
  $p[0] >= 0 and $p[0] < length and
  $p[1] >= 0 and $p[1] < (.[0]|length)
;

def at($p):
  if $p[0] < 0 or $p[0] >= length then 9
  else .[$p[0]] |
    if $p[1] < 0 or $p[1] >= length then 9
    else .[$p[1]] end
  end
;

def set($p; $to): .[$p[0]][$p[1]] = $to;

def nbr: .[0] += 1, .[0] += -1, .[1] += 1, .[1] += -1;
