#!/usr/bin/env jq -s -R -r -f
include "p20/common";
include "helpers";
include "helpers/grid";

parse |
.enhance as $enh |
[ 0, enhancefill($enh; 0) ] as [ $bg1, $bg2 ] |
.grid |
iterf(50/2;
  enhance($enh; $bg1) | enhance($enh; $bg2) |
  debugtee("completed iter")
) |
map(.[]|select(. == 1)) | length
