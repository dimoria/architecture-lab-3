#!/bin/bash

curl -X POST -d "reset
white
figure 0.5 0.5
update" http://localhost:17000

radius=0.3
angle=0

for i in {1..100}; do
  x=$(echo "c($angle)*$radius + 0.5" | bc -l)
  y=$(echo "s($angle)*$radius + 0.5" | bc -l)
  
  curl -X POST -d "move $x $y
  update" http://localhost:17000
  
  angle=$(echo "$angle + 0.1" | bc -l)
  radius=$(echo "$radius - 0.003" | bc -l)
  sleep 0.5
done