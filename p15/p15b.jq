#!/usr/bin/env jq -s -R -r -f
include "./helpers/grid";
include "p15/common";

def dupright($n):
  (.[0]|length) as $rowlen |
  reduce range(1; $n) as $i ( . ;
    map( . + ( .[-$rowlen:] | map(. % 9 + 1) ) )
  )
;

def dupdown($n):
  length as $collen |
  reduce range(1; $n) as $i ( . ;
    . + ( .[-$collen:] | map(map(. % 9 + 1)) )
  )
;

parsenumgrid |
dupright(5) |
dupdown(5) |
computerisk(bottomright) |
at(bottomright)
