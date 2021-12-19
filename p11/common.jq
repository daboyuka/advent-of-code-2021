include "./helpers";
include "./helpers/grid";

# input, output: {g: grid, flashes: number}
def flash:
  .flashes as $startflashes |
  reduce (.g | scanpoints) as $pt ( . ;
    if .g | at($pt) <= 9 then .
    else
      .flashes += 1 |
      .g |= (
        set($pt; 0) |
        reduce ($pt | nbr8) as $pt2 ( . ;
          at($pt2) as $v |
          if $v != null and $v != 0 then set($pt2; $v+1)
          else . end
        )
      )
    end
  ) |
  if .flashes == $startflashes then .
  else flash end
;

# input, output: {g: grid, flashes: number}
def step:
  .g |= map(map(.+1)) |
  flash
;
