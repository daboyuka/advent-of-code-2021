#!/usr/bin/env jq -s -R -f
include "p11/common";
include "./helpers/grid";

parsenumgrid |
( length * (.[0]|length) ) as $gridsize |
{g: ., flashes: 0, prevflashes: 0, steps: 0} |
until(
  .flashes - .prevflashes == $gridsize;
  .steps += 1 | .prevflashes = .flashes | step
) |
debug |
.steps
