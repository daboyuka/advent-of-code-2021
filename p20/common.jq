include "./helpers";
include "./helpers/grid";

def char2val: if . == "#" then 1 else 0 end;

def parse:
  lines |
  {
    enhance: ( .[0] | split("") | map(char2val)),
    grid: (.[1:] | map(split("") | map(char2val))),
  }
;

def windowrange: .[0] += (-1,0,1) | .[1] += (-1,0,1);

def windowval($pt; $bg):
  reduce ( at($pt | windowrange) | . // $bg ) as $c (
    0 ; 2*. + $c
  )
;

# output: new tile
def enhanceval($enh; $idx): $enh[$idx];

# output: new file for 3x3 of $bg
def enhancefill($enh; $bg): enhanceval($enh; 511 * $bg);

# input: grid, output: new tile
def enhancepos($enh; $pt; $bg): enhanceval($enh; windowval($pt; $bg));

# input: grid; output: enhanced grid (1 step)
def enhance($enh; $bg):
  [ range(-1;length+2) as $row |
    [ range(-1;.[0]|length+2) as $col |
      enhancepos($enh; [$row, $col]; $bg)
    ]
  ]
;
