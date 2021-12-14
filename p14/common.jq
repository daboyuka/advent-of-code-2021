include "./helpers";

def parse:
  lines | {
    start: .[0],
    rules: (
      .[1:] | map(
        split(" -> ") |
        {(.[0]): (.[1])}
      ) |
      add
    )
  }
;

# input: {"AB": "C"}, output: {"AB": ["AC", "CB"]}
def fuserules:
  .rules |= with_entries(
    .value = [.key[:1] + .value, .value + .key[1:]]
  )
;

# input: {start: string, ...}, output: add .paircounts = {"AB": 123, ...}
def countpairs:
  .paircounts = reduce ( .start | .[range(length-1):] | .[:2] ) as $pair (
    {};
    .[$pair] += 1
  )
;

def step:
  .rules as $rules |
  .paircounts |= reduce to_entries[] as $pair (
    {};
    .[$rules[$pair.key][]?] += $pair.value
  )
;

def countelems:
  reduce (.paircounts | to_entries[]) as $pair (
    # start with 1 count for start and end element (everything else will
    # be double-counted, so this allows us to divide by 2 at the end)
    {(.start[:1]): 1, (.start[-1:]): 1} ;
    .[$pair.key | split("")[]] += $pair.value
  ) |
  map_values(./2)
;
