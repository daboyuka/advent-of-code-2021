# Min-heap
include "./helpers";

# Internal function

def heapparent: (.-1) / 2 | floor;
def heaplchild: .*2+1;

def heapswap($i; $j): .[$i] as $t | .[$i] = .[$j] | .[$j] = $t;

def heapcmp($i; $j; ex): (.[$i] | ex) < (.[$j] | ex);

def heapup($idx; ex):
  ( $idx | heapparent ) as $pidx |
  if $idx == 0 or heapcmp($pidx; $idx; ex) then .
  else heapswap($pidx; $idx) | heapup($pidx; ex) end
;

def heapdown($idx; ex):
  ( $idx | heaplchild ) as $lidx |
  ( $lidx + 1 ) as $ridx |
  ( $lidx < length and heapcmp($lidx; $idx; ex) ) as $lcmp |  # $lidx < $idx
  ( $ridx < length and heapcmp($ridx; $idx; ex) ) as $rcmp |  # $ridx < $idx
  if ($lcmp|not) and ($rcmp|not) then  # $idx < $lidx, $ridx
    .
  elif $lcmp and ($rcmp|not) then      # $lidx < $idx < $ridx
    heapswap($lidx; $idx) | heapdown($lidx; ex)
  elif ($lcmp|not) and $rcmp then      # $ridx < $idx < $lidx
    heapswap($ridx; $idx) | heapdown($ridx; ex)
  elif heapcmp($lidx; $ridx; ex) then      # $lidx < $ridx < $idx
    heapswap($lidx; $idx) | heapdown($lidx; ex)
  else                                 # $ridx < $lidx < $idx
    heapswap($ridx; $idx) | heapdown($ridx; ex)
  end
;

# User-facing functions

def heapinit(ex):  # input: array, output: heap
  reduce range(length-1|heapparent; -1; -1) as $idx (
    . ; heapdown($idx; ex)
  )
;
def heapinit: heapinit(.);

def heappeek: .[0];

def heappush($v; ex): . += [$v] | heapup(length-1; ex);
def heappush($v): heappush($v; .);

def heappop(ex): heapswap(0; length-1)[:-1] | heapdown(0; ex);
def heappop: heappop(.);
