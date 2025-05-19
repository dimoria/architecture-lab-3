#!/bin/bash

# Ініціалізація: білий фон і фігура по центру
curl -X POST -d "reset
white
figure 0.5 0.5
update" http://localhost:17000

# Анімація руху по діагоналі
for i in {1..100}; do
  curl -X POST -d "move 0.002 0.002
  update" http://localhost:17000
  sleep 1
done