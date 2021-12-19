include "./helpers";

# Basic approach: parse input into lines into big array of points, then group
# by unique point and count groups with >1 instance.
# This is not fast, but fast enough. Converting points to strings and incrementing
# by key in an object might be a lot speedier.

# This file is shared between parts A and B (only difference is calling solve w/ true or false)

# types:
#   pt: [1,2]
#   vec: [[1,2],[3,4]]
#   vecset: [ [[1,2],[3,4]], ... ]

def parse:  # output: vecset
  htmlunescape |  # recover >'s
  lines |         # ["1,2 -> 3,4", ...]
  map(
    split(" -> ") | map( split(",") | map(tonumber) )
  )
;

def ptadd: transpose | map(.[1] + .[0]) ;   # input: vec
def ptdiff: transpose | map(.[1] - .[0]) ;  # input: vec

def isdiag: ptdiff | map(select(. != 0)) | length != 1 ;  # input: vec

def toline:  # input: vec, output: stream of pt
  ( ptdiff | map(fabs) | max ) as $linelen |  # compute line "length" (largest dimension diff)
  .[1] = (
    ptdiff |                     # vec to delta pt
    range($linelen + 1) as $i |  # iterate 0 to $linelen incl.
    map( . * $i / $linelen )     # scale delta pt from 0 to full, incrementing by 1 unit at a time
  ) |
  ptadd  # convert each start pt + delta pt into a pt on line
;

def solve($nodiag):  # input: vecset
    map(select($nodiag and isdiag | not)) |  # remove diagonal lines if needed
    map(toline) |  # line pts (produces some duplicate pts)
    group_by(.) |  # collect one group per unique pt containing all instances
    map(select(length > 1)) |  # convert group to group size, drop groups of size 1 (singleton pts)
    length  # count groups of size >1
;
