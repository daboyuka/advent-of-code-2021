#!/usr/bin/env jq -s -R -f
include "p09/common";

def lowpoints:
  scanpoints as $pt |
  if [ at($pt) < at($pt | nbr) ] | all then $pt else empty end
;

parse |
[ at(lowpoints) + 1 ] |
add

