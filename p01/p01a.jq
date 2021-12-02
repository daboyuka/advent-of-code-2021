#!/usr/local/bin/jq -s -f

reduce .[1:][] as $depth ( {prev: .[0], inc: 0} ;
  .inc += if $depth > .prev then 1 else 0 end |
  .prev = $depth
) |
.inc
