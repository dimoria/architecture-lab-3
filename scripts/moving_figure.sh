#!/bin/bash
curl -X GET "localhost:17000?cmd=white,figure%200.5%200.5,update"

for i in {1..100}; do
    x=$(awk "BEGIN{printf \"%.2f\", 0.01*$i}")
    y=$(awk "BEGIN{printf \"%.2f\", 0.01*$i}")
    curl -X GET "localhost:17000?cmd=move%20${x}%20${y},update"
    sleep 1
done