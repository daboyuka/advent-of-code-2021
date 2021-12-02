#!/bin/bash

# Quick CLI commands for AOC playing

DAY=$1

function get() {
  while ! curl --fail -b ~/.aoc_cookies https://adventofcode.com/2021/day/$DAY/input -o inputs/$DAY.input; do
    >&2 echo "failed to get day input; retrying"
    sleep 1
  done

  curl -b ~/.aoc_cookies https://adventofcode.com/2021/day/$DAY -o - | \
    awk '
      BEGIN{num=0;part="a";}
      /<h2 id="part2">/{part="b";}
      /<pre><code>/{gsub("<pre><code>",""); on=1;}
      /<\/code><\/pre>/{on=0;num++;}
      on{file = "inputs/'$DAY'" part ".sample" num ".input"; print > file;}
    '
}

function samp() {
  local prob=$DAY$1
  for f in inputs/$DAY?.sample*.input; do
    echo "$f"
    <"$f" go run . $prob || true
  done
}

function run() {
  local prob=$DAY$1
  <"inputs/$DAY.input" go run . $prob | tee >(tail -n1 | pbcopy)
}
