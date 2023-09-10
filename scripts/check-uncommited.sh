#!/bin/bash

uncommited=$(git ls-files -d -m -o --exclude-standard --directory --no-empty-directory)
if [[ -z "$uncommited" ]];
then
  exit 0
else
  echo "::error Uncommited changes detected"
  echo "$uncommited"
  exit 1
fi
