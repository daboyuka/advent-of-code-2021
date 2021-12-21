#!/usr/bin/env jq -s -R -f
include "p21/common";
include "./helpers";

21 as $winscore |
($winscore + 10) as $maxscore |

def makegrid(dims):
  def _rep($ignore): .;
  reduce ( [dims] | reverse[] ) as $dim ( 0;
    [ _rep(range($dim)) ]
  )
;

def dirac: [3, 1], [4, 3], [5, 6], [6, 7], [7, 6], [8, 3], [9, 1];

def play:
  def at($score; $turn; $pos):
    .[$score[0]][$score[1]][$turn][$pos[0]][$pos[1]]
  ;

  def scorecompute:  # compute every score/turn
    def poscompute($score; $turn):  # compute every position for score
      def compute($pos):  # compute combinations for score/turn/pos
        def roll($rollcombos):  # compute number of universes end in score/turn/pos after this latest roll
          $rollcombos as [$roll, $combos] |
          ( 1 - $turn ) as $prevturn |
          ( $score | .[$prevturn] -= $pos[$prevturn] + 1 ) as $prevscore |
          ( $pos   | .[$prevturn] |= (. - $roll + 10) % 10 ) as $prevpos |
          if $prevscore | min < 0 or max >= $winscore then 0
          else
            $combos * at($prevscore; $prevturn; $prevpos)
          end
        ;

        ( [ roll(dirac) ] | add ) as $tcombos |
        if $tcombos > 0 then at($score; $turn; $pos) = $tcombos #| debugtee($score, $turn, $pos, "=", $tcombos)
        else . end
      ;

      reduce ( (range(10)|[.]) + (range(10)|[.]) ) as $pos ( .; compute($pos) )
    ;

    reduce ( (range(length)|[.]) + (range(1; length*2-1)|[.]) ) as [$diagidx, $diag] ( .;
      [ ($diag - $diagidx), $diagidx ] as $score |
      if $score | min < 0 or min >= $winscore or max > $maxscore then .
      else poscompute($score; 0) | poscompute($score; 1) end
    )
  ;

  . as $start |
  makegrid($maxscore+1, $maxscore+1, 2, 10, 10) |  # [p1 score 0..$maxscore][p2 score][turn 0..1][p1 pos 0..9][p2 pos]
  .[0][0][0][$start[0]-1][$start[1]-1] = 1 |
  scorecompute
;

def wins:
  [
    ( reduce .[range($winscore;$maxscore+1)][range($winscore)][0,1][range(10)][range(10)] as $wins ( 0; . + $wins ) ),
    ( reduce .[range($winscore)][range($winscore;$maxscore+1)][0,1][range(10)][range(10)] as $wins ( 0; . + $wins ) )
  ]
;

parse |
play |
wins | max
