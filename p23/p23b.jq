#!/usr/bin/env jq -s -R -f
include "p23/common";

parse |
.[0][1:1] = [3,3] |
.[1][1:1] = [2,1] |
.[2][1:1] = [1,0] |
.[3][1:1] = [0,2] |
trimatrest |
solve
