#!/usr/bin/env sh

set -e

if [[ -z $LICHESS_API_TOKEN ]]; then
  echo "LICHESS_API_TOKEN is not set"
  exit 1
fi

sed -i "s#__LICHESS_API_TOKEN__#$LICHESS_API_TOKEN#g" /lichess-bot/config.yml

docker/copy_files.sh
python3 lichess-bot.py $OPTIONS --disable_auto_logging
