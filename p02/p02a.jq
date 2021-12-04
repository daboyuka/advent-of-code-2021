#!/usr/bin/env jq -s -R -f

reduce (                                # "as $move" block: parse input text into {dir:str, amt:num} pairs
  split("\n") | map(select(. != "")) |  # each (non-empty) line as a string
  map(
    split(" ") |                      # each line as an array of words...
    {dir: .[0], amt: (.[1]|tonumber)} # ...converted to a dir/amt pair object
  )
)[] as $move ( {x:0,d:0} ;
  if $move.dir == "forward" then .x += $move.amt
  elif $move.dir == "up" then    .d -= $move.amt
  elif $move.dir == "down" then  .d += $move.amt
  else "bad dir " + $move.dir|error
  end
) |
debug # log the actual x,d
.x * .d
