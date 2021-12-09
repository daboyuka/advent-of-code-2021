#!/usr/bin/env jq -s -R -f
include "./helpers";

# This solution does it the "hard way" by solving for all wire->segment mappings,
# first using "exclusion rules" (e.g., if you have 5 segments active, it's one of
# (2, 3, 5), all of which have segments adg on, so the two wires that are off
# cannot map to any of those segments), the second using process of elimination
# (some wires may have multiple possibilities, but other wires do not, and their
# unique segment can be removed as a possibility from all others).
#
# After doing it this way, though, I learned there are much better approaches.
# I may try implementing one in JQ later...

# type signals, segments: sorted array of letters
# type options: maps wire -> possible segments
# type problem: { samp: [ signal string, ... ], input: [ signal string, ... ] }

def tosegs: split("") | sort;

def parse:
  lines |
  map(
    [ scan("\\S+") ] |
    index("|") as $barIdx |
    {samp: ( .[:$barIdx] | map(tosegs) ), input: ( .[$barIdx+1:] | map(tosegs) )}
  )
;

def compl: ( "abcdefg" | split("") ) as $all | $all - . ;

def fulloptions:  # output: full table {"a":["a",...,"g"], ...}
  ( "abcdefg" | split("") ) as $all |
  $all | map({(.): $all}) | add
;

# pruneoptions takes an options set and recursively prunes it by
# removing known segments from multi-possibility wires until all
# wires are uniquely known.
#
# It also converts the values in the option tables from arrays of
# strings to just strings.
def pruneoptions:
  [ .[] | select(length == 1)[] ] as $certain |
  if length > ( $certain | length ) then  # if any wires still have multiple possibilities
    map_values(
      if length > 1 then . - $certain
      else .
      end
    ) |
    pruneoptions
  else map_values(.[]) # array of single string -> string
  end
;

# input: problem, output: mappings
def solvemappings:
  # First pass: use ndigits -> exclusion rules to build initial options set
  ({
    "2": {"defOn": "cf", "defOff": "abdeg"},
    "3": {"defOn": "acf", "defOff": "bdeg"},
    "4": {"defOn": "bcdf", "defOff":  "aeg"},
    "5": {"defOn": "adg", "defOff": ""},
    "6": {"defOn": "abfg", "defOff": ""},
    "7": {"defOn": "abcdefg", "defOff": ""},
  } | map_values(map_values(split("")))) as $excltable |

  # First pass: ndigits -> exclusion rules
  reduce ( .samp[] ) as $signal (
    fulloptions ;
    ($signal | length | tostring) as $ndig |
    .[$signal[]] -= $excltable[$ndig].defOff |
    .[$signal | compl[]] -= $excltable[$ndig].defOn
  ) |

  # Second pass: pruning
  pruneoptions
;

def segtodig:
  {
    "abcefg":  "0",
    "cf":      "1",
    "acdeg":   "2",
    "acdfg":   "3",
    "bcdf":    "4",
    "abdfg":   "5",
    "abdefg":  "6",
    "acf":     "7",
    "abcdefg": "8",
    "abcdfg":  "9",
  }[.]
;

# input: array of signals (array of sorted array of wire letters), output: numeric output
def applymappings($mappings):
  map(  # convert each signals into...
    map($mappings[.]) | sort | # sorted array of (correct) segments
    join("") |                 # single segments string
    segtodig                   # digit string
  ) |
  join("") | tonumber  # combine digits into single number
;

def solve:
  solvemappings as $mappings |
  .input | applymappings($mappings);

parse |
map(solve | debug) |
add
