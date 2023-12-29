#!/usr/bin/env python3

import csv
import sys

import chess
import requests

resp = requests.get(
    "https://raw.githubusercontent.com/Destaq/chess-graph/cc81f477b55e1888b42df6e85204951428be5fb3/elo_reading/openings_sheet.csv"
)

input = csv.DictReader(resp.content.decode("utf-8").splitlines())
output = csv.DictWriter(sys.stdout, fieldnames=input.fieldnames)
output.writeheader()

for row in input:
    board = chess.Board()
    moves = row["moves"].split(" ")
    uci = []

    for move in moves:
        try:
            move = board.push_san(move).uci()
            uci.append(move)

        except Exception as exc:
            print(exc, row, file=sys.stderr)
            break

    row["moves"] = " ".join(uci)
    output.writerow(row)
