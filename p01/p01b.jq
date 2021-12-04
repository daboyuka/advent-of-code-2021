#!/usr/bin/env jq -s -f

reduce .[3:][] as $depth ( {prevwin: .[:3], inc: 0} ;  # keep .prevwin as a 3-element window
  ( .prevwin[1:] + [$depth] ) as $nextwin |
  .inc += if ($nextwin|add) > (.prevwin|add) then 1 else 0 end |
  .prevwin = $nextwin
) |
.inc
