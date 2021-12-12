#!/bin/bash

# Quick CLI commands for AOC playing

DAY=$1

function 2d() { printf "%02d" "$1"; }

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

  for f in inputs/$DAY?.sample*.input; do
    sed -E -i '' -e 's|</?[^>]+>||g ; s|&lt;|<|g ; s|&gt;|>|g' "$f"
  done
}

function samp() {
  local prob=$DAY$1
  local samp=${2:-0}
  runfile "$prob" inputs/$DAY?.sample$samp.input
}

function run() {
  local prob=$DAY$1
  runfile "$prob" inputs/$DAY.input
}

function runfile() {
  local prob="$1"
  local f="$2"
  echo "$f"
  <"$f" go run . $prob | tee >(tail -n1 | pbcopy)
}

function sampjq() {
  local part="$1"
  local samp=${2:-0}
  runjqfile "$part" inputs/$DAY?.sample$samp.input
}

function runjq() {
  local part="$1"
  runjqfile "$part" inputs/$DAY.input
}

function runjqfile() {
  local day2d=$(2d $DAY)
  local part="$1"
  local f="$2"
  echo "$f"
  <"$f" p$day2d/p$day2d$part.jq | tee >(tail -n1 | pbcopy)
}
