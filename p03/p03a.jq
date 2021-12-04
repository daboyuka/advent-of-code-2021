#!/usr/bin/env jq -s -R -f

# Bits are maintained as one-char strings

# Protip: func arg without $ prefix is a full JQ filter (deferred execution), with $ is just a value
def bit($idx): .[$idx:$idx+1];

def lastbit: bit(length-1);

def inv: if . == "0" then "1" else "0" end;

def mostCommon($pos):  # takes array of bitmasks, emits most common bit (break tie to '1')
  map(bit($pos)) |
  {c: (map(tonumber)|add), l: length} |  # count 1 bits and total bits
  if 2*.c >= .l then "1" else "0" end
;

def bitsval:
  if length == 1 then tonumber
  else (lastbit | tonumber) + 2 * (.[:length-1] | bitsval)
  end;

split("\n") | map(select(. != "")) |  # each (non-empty) line as a string
(.[0] | length) as $nbits |           # get bitmask length
[ mostCommon(range($nbits)) ] |       # map array of bitmasks to array of most common bit for each pos
{ gamma: ., epsilon: map(inv) } |     # array of bits is gamma; its inverse is epsilon
map_values(join("") | bitsval) |      # convert bit arrays to bitmasks to int values
.gamma * .epsilon
