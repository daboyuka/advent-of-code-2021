include "./helpers";

def parse: lines | map(fromjson);

def isnum: type == "number";
def vecadd($a; $b): reduce range($a|length) as $i ( [] ; . += [$a[$i] + $b[$i]] );

def checktree(tree; $l):
  (
    tree |
    if isnum then .
    elif length == 2 then .[] | checktree(.; $l)
    else error("\($l.file):\($l.line): bad tree node \(.)") end |
    empty
  ), .
;

def pincside($idx; $v):
  if $v == 0 then .
  elif isnum then . + $v
  else .[$idx] |= pincside($idx; $v)
  end
;

def pexplode:
  def pexplode($depth):  # input: node, output: {n:node, adds:[l,r], updated:bool}
    def pexside($idx; $depth):  # input, output: {n:node, adds:[l,r], updated:bool}
  
      if .n[$idx] | isnum then .
      else
        ( .n[$idx] | pexplode($depth+1) ) as $child |  # explode the child node
        .n[$idx] |= $child.n |  # replace child with updated child
        .n[1-$idx] |= pincside($idx; $child.adds[1-$idx]) |  # update other side child with the other side overflow
        .adds[$idx] = $child.adds[$idx] |  # set our overall overflow to the child's overflow on this side
        .updated = true
      end
    ;
    {n: ., adds: [0,0]} |  # start with current node and no overflow
    pexside(0; $depth) |   # explode left subtree
    pexside(1; $depth) |   # explode right subtree
    if $depth >= 4 and (.n[0]|isnum) and (.n[1]|isnum) then  # explode self
      {n: 0, adds: vecadd(.adds; .n), updated: true}
    else . end
  ;
  pexplode(0) | del(.adds)
;

def psplitfirst:
  def _psplitfirst:
    if .updated then .
    elif .n | isnum then
      if .n < 10 then .
      else {n: [.n/2 | floor, ceil], updated: true} end
    else
      reduce range(2) as $i ( . ;
        ( .n |= .[$i] | _psplitfirst ) as $child |
        .n[$i] = $child.n |
        .updated = $child.updated
      )
    end
  ;
  {n: .} | _psplitfirst
;

def psimplify:
  {n: ., updated:true} |
  until(
    .updated | not;
    .n | pexplode | .n | psplitfirst
  ) |
  .n
;

def padd:  # input: array of pairs
  reduce .[1:][] as $p (
    .[0] ;
    [., $p] | psimplify
  )
;

def mag:
  if isnum then .
  else 3 * (.[0] | mag) + 2 * (.[1] | mag) end
;
