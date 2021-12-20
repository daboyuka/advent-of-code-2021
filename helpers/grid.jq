include "./helpers";

def parsegrid: lines | map(split(""));
def parsenumgrid: lines | map(split("")|map(tonumber));

def rendergrid: map(join("")) | join("\n");

def scanpoints:  # input: grid, output: all coordinates
  ((range(length)|[.]) + (range(.[0]|length)|[.]))
;

def inbounds($grid):
  .[0] >= 0 and .[0] < ($grid|length) and
  .[1] >= 0 and .[1] < ($grid[0]|length)
;

def at($p):
  if $p[0] < 0 or $p[0] >= length then null
  else .[$p[0]] |
    if $p[1] < 0 or $p[1] >= length then null
    else .[$p[1]] end
  end
;

def set($p; $to): .[$p[0]][$p[1]] = $to;

def nbr4: .[0] += (-1, 1), .[1] += (-1, 1);
def nbr8: nbr4, ( .[0] += (-1, 1) | .[1] += (-1, 1) );
