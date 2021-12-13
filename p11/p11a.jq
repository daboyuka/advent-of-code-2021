#!/usr/bin/env jq -s -R -f
include "p11/common";
include "./helpers/grid";

parsenumgrid |
reduce range(100) as $_ ( {g: ., flashes: 0} ; step ) |
debug |
.flashes
