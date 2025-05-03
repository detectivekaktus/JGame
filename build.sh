#!/usr/bin/env bash
BIN=bin

backend() {
  mkdir -p $BIN
  go build -o ./$BIN/jgame ./cmd/jgame
}

frontend() {
  cd website
  npm run build
}

all() {
  backend
  frontend
}

case "$1" in
  back)
    backend
    ;;
  front)
    frontend
    ;;
  *)
    if [ -z "$1" ]; then
      all
    else
      echo "Unknown command ${1}"
    fi
    ;;
esac
