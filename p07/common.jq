include "./helpers";

# This file is included by both parts; only the distmetric differs.
# distmetric: filter to convert distance to fuel cost
# poscounts: array of pairs (2-arrays) of [pos, count]

def parse:
  lines[0] | split(",") | map(tonumber) |  # string -> array of numbers
  group_by(.) | map([first, length])       #        -> poscounts array
;

# input: poscounts array
def fuel($at; distmetric): map((.[0] - $at | fabs | distmetric) * .[1]) | add;

# input: poscounts array
def minfuel(distmetric):
    ( map(.[0]) | max ) as $maxpos |
    [ fuel(range($maxpos+1); distmetric) ] |
    min
;
