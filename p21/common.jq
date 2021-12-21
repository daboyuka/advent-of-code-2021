include "./helpers";

def parse:
  lines |
  map(last(scan("\\d+")) | tonumber)
;
