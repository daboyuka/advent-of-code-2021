include "./helpers";
include "./helpers/grid";
include "./helpers/heap";

def bottomright: [length-1, (.[0]|length-1)];

def computerisk($forpt):  # input: grid, output: riskgrid
  . as $grid |
  {
    riskgrid: map(map(null)),
    heap: [[0,0,0]]  # heap elem: [risk, row, col]
  } |
  until(
    .riskgrid | at($forpt) != null ;
    ( .heap[0] | [.[0], .[1:]] ) as [ $risk, $pt ] |  # get next unvisited
    .heap |= heappop |
    if .riskgrid | at($pt) then .  # already locked in risk; do nothing else
    else
      .riskgrid |= set($pt; $risk) |  # lock in risk
      .heap |= reduce ( $pt | nbr4 | select(inbounds($grid)) ) as $nextpt (
        . ; heappush([$grid | $risk + at($nextpt)] + $nextpt)
      )
    end
  ) |
  .riskgrid
;
