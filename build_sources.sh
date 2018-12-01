#!/bin/bash

echo "Building all sources and moving them to ../ScoringEngine/var/lib/gingertechengine/"

for d in */ ; do
  echo "$d"
  cd "$d" || exit 1
  go build
  mv ${d::-1} ../../ScoringEngine/var/lib/gingertechengine/
  cd ..
done
