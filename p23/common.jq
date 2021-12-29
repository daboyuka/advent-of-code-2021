include "./helpers";

# types:
# kind: 0–3 (= A–D)
# prob: 2D array of kinds: [room][depth] = kind

def parse:
  lines |
  .[2:-1] | # strip constant top and bottom
  map(explode | [ .[3, 5, 7, 9] - 65 ]) |  # [depth][room]
  transpose
;

def trimatrest:
  reduce range(length) as $room ( .;
    .[$room] |= until(.[-1] != $room; .[:-1])  # pop last until != room
  )
;

def kindenergy($kind): pow(10; $kind);

def roomdist($a; $b): $a - $b | 2*fabs;

def ridx2hidx: .+1; # input: room idx, output: hall idx left of room
def hidx2pos: [0, 1, 3, 5, 7, 9, 10][.];
def ridx2pos: 2*.+2;

# minenergy returns the theoretical minimum energy to move each non-at-rest amphipod
# out of its start room and into its final room, ignoring all collisions.
def minenergy:  # input: prob, output: num
  . as $prob |
  reduce range(length) as $room ( 0;  # energy from moving out + over to (but not into) target rooms
    . + reduce ($prob[$room] | range(length)) as $depth ( 0;
      $prob[$room][$depth] as $kind |
      . + ( $depth + 1 + roomdist($room; $kind) ) * kindenergy($kind)
    )
  ) +
  reduce range(length) as $kind ( 0;  # energy from moving into target rooms
    . + (
      ( $prob[$kind] | length ) |        # num. depths in room (= num. amphis of kind)
      . * (.+1) / 2 * kindenergy($kind)  # (1 + 2 + ... +  depth) * kindenergy
    )
  )
;

# Checks if hall clear from $hidx (excl.) to the corresponding side of room $ridx (incl.)
# input: hall, output: bool
def pathclear($ridx; $hidx):
  ($ridx | ridx2hidx) as $thidx |  # target hidx (left side of room)
  if $hidx <= $thidx then
    .[$hidx+1:$thidx+1]  # $idx (excl.) to $tidx (lside, incl.)
  else
    .[$thidx+1:$hidx]  # $tidx+1 (rside, incl.) to $idx (excl.)
  end |
  map(. == -1) | all # check if all hall slots are clear
;

def canpark($idx; $roomsopen):  # $roomsopen: [kind] -> kind's room allows parking?, input: hall, output: bool
  .[$idx] as $kind |
  if $idx < 0 or $idx >= length then false # $idx out of bounds; failure
  elif $kind == -1 then true               # [$idx] is empty; success (trivally; nothing to do)
  elif $roomsopen[$kind] | not then false  # $kind is not parkable; failure
  else pathclear($kind; $idx)              # success if hall is clear (not including $idx itself)
  end
;

def parkall($hidx; $roomsopen):  # input, output: hall
  def _parkall($lhidx; $rhidx; $roomsopen):
    if   canpark($lhidx; $roomsopen) then .[$lhidx] = -1 | _parkall($lhidx-1; $rhidx; $roomsopen)
    elif canpark($rhidx; $roomsopen) then .[$rhidx] = -1 | _parkall($lhidx; $rhidx+1; $roomsopen)
    else . end
  ;
  _parkall($hidx | ridx2hidx; $hidx | ridx2hidx+1; $roomsopen)
;

def lhall($hidx):  # input: hall, output: rightmost occupied hall idx <= $hidx (or -1 if none)
  if $hidx == -1 or .[$hidx] != -1 then $hidx
  else lhall($hidx-1) end
;
def rhall($hidx):  # input: hall, output: leftmost occupied hall idx >= $hidx (or length if none)
  if $hidx == length or .[$hidx] != -1 then $hidx
  else rhall($hidx+1) end
;
def hallrange($room):  # input: hall, output: hall idxs accessible from $room by increasing distance
  ( $room | ridx2hidx ) as $lhpos |  # hpos on left of room
  [range(
    lhall($lhpos) + 1;
    rhall($lhpos + 1)
  )] |
  sort_by(hidx2pos - ( $room | ridx2pos ) | fabs)[]
;

def excessdist($room; $kind; $viahidx):  # excess dist = room->viahidx + viahidx->kind - room->kind
  ($room | ridx2pos) as $rpos |
  ($kind | ridx2pos) as $trpos |
  ($viahidx | hidx2pos) as $viahpos |
  ( $rpos - $viahpos | fabs ) +
  ( $trpos - $viahpos | fabs ) -
  ( $trpos - $rpos | fabs )
;

# input: state, output: bool
def roomsopen: map(length == 0);  # [room] -> room is open?
def roomsclear: roomsopen | all;  # all rooms clear of non-at-rest
def hallclear:  isempty(.[] | select(. != -1));  # hall clear

def nextrooms:  # input: prob, output: non-empty room idxs, fewer non-at-rest first
  . as $prob |
  [ range(length) ] |
  sort_by(-($prob[.][0] // -1))[]
;

def solve:  # input: state, output: best extra energy
  def _solve($hall; $energy; $bestEnergy):
    . as $prob |
    if $energy >= $bestEnergy then null
    elif $prob | roomsclear then
      if $hall | hallclear then $energy | debug else null end
    else
      reduce ( $prob | nextrooms ) as $room ( $bestEnergy;
        if $prob[$room] | length == 0 then . else
          $prob[$room][0] as $kind |             # get kind of top amphi in room
          kindenergy($kind) as $kindenergy |
          ( $prob | .[$room] |= .[1:] ) as $prob |  # remove top amphi from room
          ( $prob | roomsopen ) as $roomsopen    |  # capture list of open rooms
          reduce ( $hall | hallrange($room) ) as $hidx ( .;
            . as $bestEnergy |
            $prob | _solve(
              $hall | .[$hidx] = $kind | parkall($room; $roomsopen);
              $energy + excessdist($room; $kind; $hidx) * $kindenergy;
              $bestEnergy
            ) // $bestEnergy
          )
        end
      )
    end
  ;
  minenergy + _solve([-1,-1,-1,-1,-1,-1,-1]; 0; infinite)
;
