include "./helpers";

def hextobits: {
    "0": [0,0,0,0],
    "1": [0,0,0,1],
    "2": [0,0,1,0],
    "3": [0,0,1,1],
    "4": [0,1,0,0],
    "5": [0,1,0,1],
    "6": [0,1,1,0],
    "7": [0,1,1,1],
    "8": [1,0,0,0],
    "9": [1,0,0,1],
    "A": [1,0,1,0],
    "B": [1,0,1,1],
    "C": [1,1,0,0],
    "D": [1,1,0,1],
    "E": [1,1,1,0],
    "F": [1,1,1,1],
  }[.]
;

def bitstonum: reduce .[] as $b ( 0 ; 2*. + $b );

# packet functions:
# input: {bits: [...], <packet fields>, output: {bits: [...<fewer>...], <more packet fields>}
#
# Packet fields: {
#   ver: number,
#   typ: number,
#   val: number,  # literal val
#   pks: array, packets  # sub-packets
# }

def pp:
  .ver = (.bits[0:3] | bitstonum) |
  .typ = (.bits[3:6] | bitstonum) |
  .bits |= .[6:] |
  if .ver == 4 then pplit else ppop end
;

def pplit:
  .bits[0] == 0 as $islast |
  .val = (.val // 0)*16 + ( .bits[1:5] | bitstonum ) |
  .bits |= .[5:] |
  if $islast then .
  else pplit end  # combine in another lit word
;

def ppop:
  def ppframed($cbits; pframed): ( .bits[:15] | bitstonum ) as $v |
  if .[0] == 0 then
    ( .bits[:15] | bitstonum ) as $nbits |
    .bits |= .[15:] | ppsbits($nbits)
  else
    ( .bits[:11] | bitstonum ) as $nps |
    .bits |= .[11:] | ppsbits($nps)
  end
;

def ppsbits($nbits):
  ( .bits | length ) - $nbits as $endbits |
  until(
    .bits | length == $endbits ;
    pp as $p |
    .pks |= ( . // [] ) + [ $p | del(.bits) ] |
    .bits = $p.bits
  )
  if $nbits == 0 then return .
  else
    pp as $p |
    $nbits - (.bits|length) + ($p.bits|length) as $nbitsleft |
    .pks |= ( . // [] ) + [ $p | del(.bits) ] |
    .bits = $p.bits |
    ppsbits($nbitsleft)  # read another sub-packet
  end

def ppscount($p):

def parsebits: lines[0] | [ split("")[] | hextobits[] ];

def parse:
  parsebits | pp;
