#!/usr/bin/env bash

set -e

depth=$1
position=$2
moves=$3

if [ -z "$depth" ] || [ -z "$position" ]; then
  echo "Usage: $0 <depth> <position> [moves]"
  exit 1
fi

if [ "$position" == "startpos" ]; then
  position="rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
fi

command="position fen $position"

if [ -n "$moves" ]; then
  command="$command moves $moves"
fi

command="$command\ngo perft $depth"

chester=$(echo -e $command | go run . uci | grep .)
chester_nodes=$(echo "$chester" | grep 'nodes:' | awk '{print $2}')
chester_moves=$(echo "$chester" | grep -v 'nodes:' | sed 's/^ //' | sed 's/://')

echo -e "$chester_moves\n\n$chester_nodes"
