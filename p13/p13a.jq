#!/usr/bin/env jq -s -R -r -f
include "p13/common";

parse |
.folds |= .[:1] |
folds |
unique |
length
