include "./helpers";

def parsegrid: lines | map(split(""));
def parsenumgrid: lines | map(split("")|map(tonumber));

def rendergrid: map(join("")) | join("\n");

def griddims: [., .[0] | length];  # returns [rows, cols]

# output: all coordinates between $startpt (incl.) and $endpt (excl.)
def scanpoints($startpt; $endpt):
  (range($startpt[0]; $endpt[0])|[.]) + (range($startpt[1]; $endpt[1])|[.])
;
# input: grid, output: all coordinates
def scangrid: scanpoints([0, 0]; griddims);

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
