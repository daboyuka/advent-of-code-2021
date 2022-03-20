#!/usr/bin/env jq -f

[inputs] as $depths
| [$depths[:-3], $depths[3:]]
| transpose
| map(select (.[1] > .[0]))
| length
