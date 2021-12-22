include "./helpers";

# types:
# pos3: [X,Y,Z]
# cube: [pos3,pos3] = [min,max]
# step: [true/false, cube] = [on?, cube]

def parse:
  lines | map(
    split(" ") |
    (.[0] == "on") as $on |
    .[1] | split(",") | map(  # ["x=1..2", ...]
      .[2:] | split("..") | map(tonumber)  # [[1,2], ...]
    ) |
    transpose |  # [ [min=x,y,z], [max=x,y,z] ]
    {
      "on": $on,
      "min": .[0],
      "max": (.[1] | map(.+1)),  # increment max to be exclusive
    }
  )
;

def splitcube($dim; $at):
  if $at < .min[$dim] or $at > .max[$dim] then error("bad splitcube") else
    ( .max[$dim] = $at ),
    ( .min[$dim] = $at )
  end
;

def zerocube: {"on":false,"min":[0,0,0],"max":[0,0,0]};

def isdegen: .min[0] >= .max[0] or .min[1] >= .max[1] or .min[2] >= .max[2];
def isvalid: isdegen | not;

def intersects($other):
  def idim($other; $dim): .max[$dim] > $other.min[$dim] and .min[$dim] < $other.max[$dim];
  idim($other; 0) and idim($other; 1) and idim($other; 2)
;

def partdim($a; $b; $dim):  # computes $a - $b; input: cube, output: [updated $a, updated $b, other splits ...]
  [$a, $b] |
  [ .[0,1].min[$dim] ] as [$mina, $minb] |
  if $mina < $minb then  # a extends lower than b; lop off lower a as split
    ( .[0] | [ splitcube($dim; $minb) ] ) as [ $low, $mid ] |
    . += [ $low ] |
    .[0] = $mid
  else  # b extends lower than a; delete lower b
    .[1] |= last(splitcube($dim; $mina))
  end |
  [ .[0,1].max[$dim] ] as [$maxa, $maxb] |
  if $maxa > $maxb then  # a extends higher than b; lop off higher a as split
    ( .[0] | [ splitcube($dim; $maxb) ] ) as [ $mid, $high ] |
    . += [ $high ] |
    .[0] = $mid
  else  # b extends higher than a; delete higher b
    .[1] |= first(splitcube($dim; $maxa))
  end
;

def cubeinter($b):
  if intersects($b) | not then null else
    reduce range(3) as $dim (
      [., $b];
      partdim(.[0]; .[1]; $dim)[:2]
    ) |
    .[0]
  end
;

def cubesub($b):
  if intersects($b) | not then . else
    foreach range(3) as $dim (
      [., $b];
      partdim(.[0]; .[1]; $dim);
      .[2:][]
    ) |
    select(isvalid)
  end
;

def cubevol: . as $c | reduce range(3) as $i ( 1; . * ( $c.max[$i] - $c.min[$i] ) );

def sumcubes:
  reduce .[] as $step ( [];
    if $step.on then
      map(cubesub($step)) + [$step]
    else
      map(cubesub($step))
    end
  )
;
