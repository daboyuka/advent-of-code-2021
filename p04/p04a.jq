#!/usr/local/bin/jq -s -R -f

include "./helpers";

def parse:
  linegroups |
  {
    numbers: ( .[0][0] | split(",") ),
    bingos: ( .[1:] |
      map(map([ scan("\\w+") ]))  # array of (bingo boards = array of (string -> array of words))
    )
  }
;

# takes bingo board, emits marked bingo board
def mark($num): map(map(if . == $num then "" else . end));

def wonrow: map(. == "") | all;    # array of items -> array of marked -> all marked?
def wonrows: map(wonrow) | any;    # array of rows -> array of all-marked? -> any all-marked?
def woncols: transpose | wonrows;
def won: wonrows or woncols;

def bingosum: [.[][] | select(. != "") | tonumber] | add;

def play:
  .numbers[0] as $n |         # capture current number
  .numbers |= .[1:] |         # advance number
  .bingos |= map(mark($n)) |  # advance all bingos
  (
    .bingos[] | select(won) |                # find all winner bingos
    { sum: bingosum, n: ( $n | tonumber ) }  # convert to sum and current number
  ) // play                                  # if there were no winners, play again
;

parse |
play |
debug |
.sum * .n
