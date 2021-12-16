#!/usr/bin/env jq -s -R -f
include "./helpers";

def parse:  # output: { "from": ["to1", "to2", ...], ... }
  reduce lines[] as $edge ( {} ;
    ( $edge | split("-") ) as [$a, $b] |
    .[$a] += [$b] |
    .[$b] += [$a]
  )
;

def isucase: .[:1] | . == ascii_upcase;

def countpaths($mayDV):  # input: edge map, output: number of paths ($mayDV = may double visit)
  . as $edges |
  def traverse($v; $mayDV):  # input: current node, output: number of paths ($v = visited map)
    def visit:  # input: current node, output: number of paths
      if . == "start" then 0  # start = dead end
      elif . == "end" then 1  # end = 1 successful path
      elif isucase     then traverse($v; $mayDV)                # ucase = freely continue
      elif $v[.] | not then traverse($v + {(.): true}; $mayDV)  # unvisted lcase = mark visited, continue
      elif $mayDV      then traverse($v; false)                 # visited lcase but allow double = continue but disallow double
      else 0 end  # visited lcase and disallow double = dead end
    ;
    reduce $edges[.][]? as $next ( 0 ; . + ( $next | visit ) )
  ;
  "start" | traverse({}; $mayDV)
;
