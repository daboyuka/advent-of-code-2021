#!/usr/bin/env jq -s -R -f
include "./p05/p05common";

parse | solve(false)  # with diagonals
