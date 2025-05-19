#!/bin/bash

# Встановлюємо зелений фон і чорний прямокутник
curl -X POST -d "green
bgrect 0.1 0.1 0.9 0.9
update" http://localhost:17000