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
  def pexplode($depth):  # input: node, output: {n:node, adds:[l,r]}
    def pexside($idx; $depth):  # input, output: {n:node, adds:[l,r]}
      if .[$idx] | isnum then ., 0
      else
        [ .[$idx] | pexplode($depth+1) ] as $child |  # explode the child node
        .[$idx] = $child[0] |  # replace child with updated child
        .[1-$idx] |= pincside($idx; $child[2-$idx]) |  # update other side child with the other side overflow
        ., $child[1+$idx]
      end
    ;
    [ pexside(0; $depth) ] as [$n, $addl] | $n |  # explode left subtree
    [ pexside(1; $depth) ] as [$n, $addr] | $n |  # explode right subtree
    if $depth >= 4 and (.[0]|isnum) and (.[1]|isnum) then  # explode self
      0, $addl + .[0], $addr + .[1]
    else
      ., $addl, $addr
    end
  ;
  first(pexplode(0))
;

def psplitfirst:  # input: node, output: node, bool (updated)
  if isnum then
    if . < 10 then ., false
    else [./2 | floor, ceil], true end
  else
    reduce range(2) as $i ( [., false] ;
      if .[1] then .
      else
        [ .[0][$i] | psplitfirst ] as $child |
        .[0][$i] = $child[0] |
        .[1] = $child[1]
      end
    ) |
    .[]
  end
;

def psimplify:
  [., true] |
  until(
    .[1] | not;
    .[0] | pexplode | [psplitfirst]
  )[0]
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
