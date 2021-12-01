#!/bin/bash

# Run from main dir
function 2d() { printf "%02d" $1; }
for ((i=1;i<=25;i++)); do
  j=$(2d $i);
  mkdir -p p$j
  sed "s/XXX/$j/" setup/template.go.txt > p$j/p$j.go
done
