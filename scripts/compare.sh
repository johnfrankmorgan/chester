#!/usr/bin/env bash

set -e

old=$1
new=$2

if [ -z "$old" ]; then
  echo "Usage: $0 <old> [new]"
  exit 1
fi

rm -rf tmp/clone

git clone -q . tmp/clone
cd tmp/clone

if [ -z "$new" ]; then
  cd ../..
  go build -o tmp/bin/new .
  cd tmp/clone
else
  git checkout $new
  go build -o ../../tmp/bin/new .
fi

git checkout $old
go build -o ../../tmp/bin/old .

cd ../..

cutechess-cli \
  -engine name=old cmd=tmp/bin/old proto=uci \
  -engine name=new cmd=tmp/bin/new proto=uci \
  -concurrency ${CONCURRENCY:-8} \
  -each tc=${TC:-30+2} \
  -rounds ${ROUNDS:-10} \
  -sprt elo0=0 elo1=10 alpha=0.05 beta=0.05
