#!/usr/bin/env jq -s -R -f
include "./helpers";

# input: string, output: single string
def check:
  {"(":")", "[":"]", "{":"}", "<":">"} as $braces |
  split("") |  # string -> array of chars
  reduce .[] as $c (  # -> {corrupt: "X"} or {"missingtail": "XYZ"}
    {stack: []};
    if .corrupt then .  # once we've hit a corrupt symbol, stop processing
    elif $braces[$c]               then .stack += [$c]  # opening -> push
    elif $braces[.stack[-1]] != $c then {corrupt: $c}   # bad closing -> corrupt
    else .stack |= .[:-1]
    end
  ) |
  if .corrupt then .
  else
    { missingtail: (
      .stack |
      reverse |
      map($braces[.]) |
      join("")
    )}
  end
;
