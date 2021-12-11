#!/usr/bin/env jq -s -R -f
include "./helpers";

# Attempt at the "easy way" solution: don't discover the mappings between wires and segments,
# Just discover which set of wires indicates which digit.
#
# Start with the simple digits: 1, 4, 7, 8
# Then the following rules can discover the additional digits:
# 0 / 6 / 9: all have 6 segs
#   9: contains 4
#   6: not contains 1
#   0: otherwise
# 2 / 3 / 5:
#   3: if contains 1 or 7
#   5: if contains (4 - 1)
#   2: otherwise
#

def tosegs: split("") | sort;

def union($a; $b): error ; # ( $a | explode ) + ( $b | explode ) | unique | implode;
def inter($a; $b): error ; # ( $a | explode ) + ( $b | explode )
def diff($a; $b): ($a | explode) - (b | explode) | implode;

def parse:
  lines |
  map(
    split(" | ") |
    map( [ scan("\\S+") | tosegs ] ) |
    {samp: .[0], input: .[1]}
  )
;

def initstate:
  .known = [range(10) | null]
;

# Discover mapping signal -> digit
def discover:
  if .known | map(select(. == null)) | all then .
  else
    .known = reduce .samp[] as $seg ( .known ;
      ( $seg | length ) as $l |
      if $l ==  2 then   .[1] = $seg
      elif $l ==  3 then .[7] = $seg
      elif $l ==  4 then .[4] = $seg
      elif $l ==  7 then .[8] = $seg
      elif $l ==  5 then  # -> 2 / 3 / 5
        [.[1], .[7], (.[4]//[]) - (.[1]//[]) | . // [] | inside($seg)] as [$has1, $has7, $has4m1] |
        if (.[1] and $has1) or (.[7] and $has7) then              .[3] = $seg
        elif .[1] and .[4] and $has4m1 then                       .[5] = $seg
        elif .[1] and .[4] and ($has1|not) and ($has4m1|not) then .[2] = $seg
        else . end
      else  # $l ==  6 -> 0 / 6 / 9
        [.[1], .[4] | . // [] | inside($seg)] as [$has1, $has4] |
        if .[4] and $has4 then         .[9] = $seg
        elif .[1] and ($has1|not) then .[6] = $seg
        elif .[1] and .[4] and $has1 and ($has4|not) then .[0] = $seg
        else . end
      end
    ) |
    .samp = .samp - .known |
    discover
  end
;

def tomapping:  # converts known array to a mapping object
  to_entries |
  map({ (.value | join("")): (.key | tostring) }) |
  add
;

# input: array of segments, output: single number
def applymapping($mapping): .input | map($mapping[join("")]) | join("") | tonumber;

parse | map(
  initstate |
  discover |
  applymapping(.known | tomapping)
) |
add
