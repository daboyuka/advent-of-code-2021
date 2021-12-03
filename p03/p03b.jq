#!/usr/local/bin/jq -s -R -f

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

# takes array of bitmaps; outputs the single element if length == 1, else a
# subset of array with bit at $pos equal to the most (crit = .) or
# least (crit = inv) common bit at $pos
def filterbitmaps(crit; $pos):
  if length == 1 then .[0] else
    ( mostCommon($pos) | crit ) as $tbit |  # find most common bit at $pos and modulate with crit
    map(select( bit($pos) == $tbit )) |     # keep only bitmaps matching the target bit
    filterbitmaps(crit; $pos+1)             # do another filter pass
  end;

split("\n") | map(select(. != "")) |  # each (non-empty) line as a string
{
  o2:  filterbitmaps(.; 0),    # o2 uses most common bit as criterion
  co2: filterbitmaps(inv; 0),  # co2 uses least common bit (most common inverted) as criterion
} |
map_values(bitsval) |  # convert bitmaps to int values
.o2 * .co2
