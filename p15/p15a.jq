#!/usr/bin/env jq -s -R -f
include "./helpers/grid";
include "p15/common";

parsenumgrid |
computerisk(bottomright) |
at(bottomright)
