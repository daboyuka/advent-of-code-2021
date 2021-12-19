include "./helpers";

# This file is shared between parts A and B (only difference is calling
# nevolve w/ 80 or 256).

def parse:
  lines[0] | split(",") | map(tonumber) |             # array of ages
  reduce .[] as $age ( [range(9)|0] ; .[$age] += 1 )  # array of counts by age
;

def evolve:
  .[7] += .[0] |  # 0-age fish become 7-age fish (about to be downshifted to 6)
  .[1:] + [.[0]]  # downshift all ages, and add 8-age fish equal to 0-age fish
;

def nevolve($n): if $n == 0 then . else evolve | nevolve($n-1) end;
