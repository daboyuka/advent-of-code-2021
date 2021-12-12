#!/usr/bin/env jq -s -R -f
include "./helpers";
include "p10/common";

{")":3, "]":57, "}":1197, ">":25137} as $score |
lines |
map(check | .corrupt | select(.)) |  # check each string, capture array of corrupt braces
map($score[.]) | add
