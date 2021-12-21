include "./helpers";

# types:
# rot3: [[ra, rn], [ua, un], [fa, fn]]  # r, u, f = right, up, fwd; a, n = axis (int), negative (bool)
# pos3: [X, Y, Z]
# trans3: [rot3, pos3]
# scanner: [pos3, ...] (sorted)

def parse:  # input: text, output: array of (scanner = sorted array of (point = 3-array of numbers))
  linegroups |
  map(
    .[1:] |  # skip --- scanner X --- header
    map(split(",") | map(tonumber))
  )
;

# rot3 defines a rotational translation from a relative frame of reference back to the absolute
# rot3: [ [rightDim, rightNeg], [upDim, upNeg], [fwdDim, fwdNeg] ]

def mkrot3s:  # output: stream of all 24 rot3s
  null |
  range(3) as $right |
  (false, true) as $rightN |
  (range(3) | select(. != $right)) as $up |
  (false, true) as $upN |
  (3 - $right - $up) as $fwd |
  (($rightN != $upN) != ((3 + $right - $up) % 3 == 1)) as $fwdN |  # $rightN XOR $upN XOR ($right - $up == 1 mod 3)
  [[$right, $rightN], [$up, $upN], [$fwd, $fwdN]] as $uninv |  # [i] -> mapping from input dim i to output dim
  reduce range(3) as $i ( [null,null,null] ;                        # invert so that [i] -> mapping for output dim i from input dim
    .[$uninv[$i][0]] = [$i, $uninv[$i][1]]
  )
;

def applyrot3($rot):  # input, output: pos3
  . as $inpos |
  $rot | map($inpos[.[0]] * if .[1] then -1 else 1 end)
;

def negpos3: map(-.);  # input, output: pos3
def addpos3($other): reduce range(length) as $i ( .; .[$i] += $other[$i] );  # input, output: pos3
def subpos3($other): negpos3 | addpos3($other);  # input, output: pos3

def applytrans3($trans): applyrot3($trans[0]) | addpos3($trans[1]);  # input, output: pos3

def zerotrans3: [ [[0, false], [1, false], [2, false]], [0,0,0] ];

def countmatches($absscanner):  # input: rel. scanner, output: num matching points
  length - (. - $absscanner | length)
;

def absolutize($absscanner):  # input: [ rel scanner, null ], output: [ abs scanner, trans3 ] or null if no match
  .[0] |
  length as $rellen |
  ($absscanner | length) as $abslen |
  [mkrot3s] as $allrots |
  ( $allrots | length ) as $rotslen |
  12 as $thresh |

  def nextcombo:  # input, output: [rotidx, absidx, relidx]
    .[0] |= ( . + 1 ) % $rotslen |
    if .[0] != 0 then .
    else .[1] |= ( . + 1 ) % $abslen |
      if .[1] != 0 then .
      else .[2] |= ( . + 1 ) % $rellen |
        if .[2] != 0 then .
        else null
        end
      end
    end
  ;

  def iter($combo):  # input: relscanner, output: trans3 or null
    def check:  # input: relscanner, output: trans3 or null
      if $combo == null then [., null]
      else
        $combo as [ $rotidx, $absidx, $relidx ] |
        $allrots[$rotidx] as $rot |
        (
          .[$relidx] |
          applyrot3($allrots[$rotidx]) |
          negpos3 |
          addpos3($absscanner[$absidx])
        ) as $off |
        [ $rot, $off ] as $trans |
        map(applytrans3($trans)) |
        if countmatches($absscanner) >= $thresh then [., $trans]
        else null end
      end
    ;
    check // iter($combo | nextcombo)
  ;
  iter([0, 0, 0])
;

def findabs: to_entries | map(select(.value[1] != null).key);

# matches a single absolute scanner to all relative scanners, as possible
def applyabs($absidx):  # input, output: array of [ scanner (abs or rel), trans3 (or null) ]
  .[$absidx][0] as $absscanner |
  map(if .[1] then . else absolutize($absscanner) end)
;

def solve:  # input, output: array of [ scanner (abs or rel), trans3 (or null) ]
  def _solve($tocheck):
    debugtee("checking", $tocheck) |
    if map(.[1] != null) | all then .
    else
      findabs as $absbefore |
      reduce $tocheck[] as $i ( . ; applyabs($i) ) |
      _solve(findabs - $absbefore)
    end
  ;
  map([., null]) | .[0][1] = zerotrans3 | _solve([0])
;
