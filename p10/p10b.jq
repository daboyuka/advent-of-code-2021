#!/usr/bin/env jq -s -R -f
include "./helpers";
include "p10/common";

def score:
  {")":1, "]":2, "}":3, ">":4} as $bracescore |
  reduce split("")[] as $c ( 0 ; .*5 + $bracescore[$c] )
;

lines |
map(check | select(.corrupt|not).missingtail) |  # check each string, capture array of missing tails
map(score) |
debug |
sort[(length-1)/2]
