# Acts as debug, except values are filtered by "what" writing to stderr (values flow through unaffected).
def debugtee(what): (what | debug | empty), . ;

# Splits input as per -sR into non-empty lines
def lines: split("\n") | map(select(. != "")) ;

def linegroups:
  reduce split("\n")[] as $line ( [[]] ;
    if $line == "" then . + [[]]
    else                .[-1] += [$line]
    end
  ) |
  map(select(length > 0)) # only keep non-empty linegroups
;
