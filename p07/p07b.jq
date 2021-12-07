#!/usr/bin/env jq -s -R -f
include "p07/common";

parse |
minfuel(. * (.+1) / 2)
