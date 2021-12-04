#!/usr/bin/env jq -s -R -f

reduce (                                # "as $move" block: parse input text into {dir:str, amt:num} pairs
  split("\n") | map(select(. != "")) |  # each (non-empty) line as a string
  map(
    split(" ") |                      # each line as an array of words...
    {dir: .[0], amt: (.[1]|tonumber)} # ...converted to a dir/amt pair object
  )
)[] as $move ( {x:0,d:0,a:0} ;
  if $move.dir == "forward" then .x += $move.amt | .d += $move.amt * .a
  elif $move.dir == "up" then    .a -= $move.amt
  elif $move.dir == "down" then  .a += $move.amt
  else "bad dir " + $move.dir|error
  end
) |
debug # log the actual x,d,a
.x * .d
