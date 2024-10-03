#!/usr/bin/env bash

set -e

old=$1
new=$2

if [ -z "$old" ]; then
  echo "Usage: $0 <old> [new]"
  exit 1
fi

mkdir -p tmp

rm -rf tmp/clone
git diff HEAD >tmp/clone.diff

git clone -q . tmp/clone
cd tmp/clone

if [ -z $new ]; then
  git diff HEAD >../clone.diff
  git tag new
  new=new
fi

git checkout $old
go build -o ../../tmp/bin/old .

git checkout $new
if [ -f ../clone.diff ]; then
  git apply --allow-empty ../clone.diff
fi
go build -o ../../tmp/bin/new .

cd ../..

cutechess-cli \
  -engine name=old cmd=tmp/bin/old proto=uci \
  -engine name=new cmd=tmp/bin/new proto=uci \
  -each tc=10+0.1 \
  -rounds 100 \
  -sprt elo0=0 elo1=10 alpha=0.05 beta=0.05
