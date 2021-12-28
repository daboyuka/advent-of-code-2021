#!/usr/bin/env jq --raw-input -f

def parseCommand: 
	split(" ")
  | {type:.[0], amount:(.[1] | tonumber)};

[., inputs] 
| map (parseCommand)
| reduce .[] as $command (
	{x: 0, y: 0, aim: 0};
	if $command.type=="forward" then 
		.+{x: (.x + $command.amount),
		   y: (.y + (.aim*$command.amount))} 
	elif $command.type=="down" then
		.+{aim: (.aim + $command.amount)}
	elif $command.type=="up" then
		.+{aim: (.aim - $command.amount)}
	else
		.
	end
)
| .x * .y