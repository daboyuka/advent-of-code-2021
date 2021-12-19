# Acts as debug, except values are filtered by "what" writing to stderr (values flow through unaffected).
def debugtee(what):
  (
    [what] |
    if length > 1 then map(tostring) | join(" ")
    else . end |
    debug |
    empty
  ), .
;

# Splits input as per -sR into non-empty lines
def lines: split("\n") | map(select(. != "")) ;

def htmlunescape: gsub("&gt;" ; ">") | gsub("&lt;" ; "<") ;

def linegroups:
  reduce split("\n")[] as $line ( [[]] ;
    if $line == "" then . + [[]]
    else                .[-1] += [$line]
    end
  ) |
  map(select(length > 0)) # only keep non-empty linegroups
;

def filterfirst(ffirst; f): [f] | (.[0] | ffirst), .[1:][];

def assert(f; err; $loc):
    (if $loc then "\($loc.file):\($loc.line): " else "" end) as $pre |
    ( f // error($pre + (err|tostring)) | empty ), .;
def assert(f; err): assert(f; err; null);

def asserttype($t; $loc): assert(type == $t; "expected \($t), got \(type)"; $loc);
def asserttype($t): asserttype($t; null);
