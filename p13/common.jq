include "./helpers";

def parse:
  linegroups |
  {
    points: (
      .[0] | map(
        split(",") | map(tonumber)
      )
    ),
    folds: (
      .[1] | map(
        ltrimstr("fold along ") |
        split("=") |
        {
          axis:  ( .[0] | if . == "x" then 0 else 1 end ),
          coord: ( .[1] | tonumber )
        }
      )
    )
  }
;

def fold($fold):  # input, output: grid
  map(
    if .[$fold.axis] < $fold.coord then .
    else .[$fold.axis] |= 2*$fold.coord - . end
  )
;

def folds:
  reduce .folds[] as $fold (
    .points ;
    fold($fold)
  )
;

def render:
  [ ( 0, 1 ) as $idx | map(.[$idx]) | min, max ] as [$minX, $maxX, $minY, $maxY] |
  reduce .[] as $pt (
    [range($maxY - $minY + 1) | [range($maxX - $minX + 1) | " "]] ;
    .[$pt[1]][$pt[0]] = "#"
  ) |
  map(join("")) | join("\n")
;
