include "./helpers";

def parse:  # input: text, output: [ [minX, maxX], [minY, maxY] ]
  lines[0] |
  ltrimstr("target area: ") |
  split(", ") |
  map(.[2:] | split("..") | map(tonumber))
;

def ininterv($bound): . >= $bound[0] and . <= $bound[1];
def ltinterv($bound): . < $bound[0];
def gtinterv($bound): . > $bound[1];

def inbox($box): [ (0,1) as $i | .[$i] | ininterv($box[$i]) ] | all;
def pastbox($box): (.[0] | gtinterv($box[0])) or (.[1] | ltinterv($box[1]));

def sgn: if . < 0 then -1 elif . > 0 then 1 else 0 end;

def simulate($box):  # input: [xvel, yvel], output: { ... }
  {pos: [0,0], vel: ., peak: 0} |
  until(
    .pos | inbox($box) or pastbox($box);
    {
      pos: [ (0,1) as $i | .pos[$i] + .vel[$i] ],
      vel: [ ( .vel[0] | . - sgn ), .vel[1] - 1 ],
      peak: ( [.peak, .pos[1]] | max )
    }
  ) |
  .hit = ( .pos | inbox($box) )
;

def simulateall($box):
  ( $box[0][1] | fabs ) as $absmaxX |
  ( $box[1][0] | fabs ) as $absminY |
  ( range($absmaxX + 1) | [.] ) +
  ( range(-$absminY; $absminY+1) | [.] ) |  # all pairs between [0, -|minY|] and [|maxX|, |minY|] inclusive
  simulate($box)
;
