#!/usr/bin/env jq -s -R -f
include "p22/common";

{min:[-50,-50,-50], max:[51,51,51]} as $bound |

parse |
map(cubeinter($bound)) |
sumcubes |
map(cubevol) | add
