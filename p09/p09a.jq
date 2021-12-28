#!/usr/bin/env jq -s -R -f
include "./helpers/grid";

def lowpoints:
  scangrid as $pt |
  [ ( at($pt)//9 ) < ( at($pt | nbr4)//9 ) ] |  # neighbors are higher?
  if all then $pt else empty end
;

parsenumgrid |
[ at(lowpoints) + 1 ] |
add

