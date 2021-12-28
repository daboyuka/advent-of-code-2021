#!/usr/bin/env jq -s -R -r -f
include "./helpers";
include "./helpers/grid";

# floodfill fills a connected region (of cells < 9) with the given value
# input: {g: grid, tofill: [...]}, output: {g: updatedGrid, tofill: []}
def floodfill($fill):
  reduce .tofill[] as $pt ( {g, tofill: []};
    if .g | at($pt)//9 | . == $fill or . == 9 then .
    else
      .g |= set($pt; $fill) |
      .tofill += [$pt | nbr4]
    end
  ) |
  if .tofill | length == 0 then .
  else floodfill($fill) end
;

# fillscan scans the entire grid, and floodfills each region (of cells < 9)
# with incrementally higher values, starting at 10
def fillscan:  # input grid, output grid
  reduce scangrid as $pt ( {g: ., nextfill: 10};
    if .g | ( at($pt)//9 ) >= 9 then .
    else
      .g = ( .tofill = [$pt] | floodfill(.nextfill).g ) |
      .nextfill += 1
    end
  ) |
  .g
;

# basinsizes computes all basin sizes
def basinsizes:
  fillscan |     # fillscan to fill in connected regions with unique values (> 9)
  flatten |      # flatten grid into simple array of cell values
  group_by(.) |  # one group per unique cell value
  map(select(first > 9) | length)  # count instances of each cell value >9 (i.e., floodfill values)
;

parsenumgrid |
basinsizes |
sort |
.[-1] * .[-2] * .[-3]
