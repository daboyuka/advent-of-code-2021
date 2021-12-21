#!/usr/bin/env jq -s -R -f
include "p21/common";
include "./helpers";

def play:
  {
    pos: map(.-1), # convert to 0-based
    score: [0, 0],
    turn: 0,
    nextRoll: 0, # 0-based
    totalRolls: 0,
  } |
  until(
    .score | max >= 1000 ;
    ( [.nextRoll | (., .+1, .+2) | .%100+1] | add ) as $troll |
    .nextRoll |= (.+3)%100 |
    .totalRolls += 3 |
    .pos[.turn] |= (. + $troll) % 10 |
    .score[.turn] += .pos[.turn]+1 |
    .turn |= 1-.
  ) |
  .pos |= map(.+1) # convert back to 1-based
;

parse |
play |
(.score | min) * .totalRolls
