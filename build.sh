#!/usr/bin/env bash
BIN=bin
TARGET=jgame
export BACKEND_DEV=1

backend() {
  if [[ $BACKEND_DEV = 1 ]]; then echo "Building backend in development mode...";
  else echo "Building backend in production mode..."; fi
  mkdir -p $BIN
  go build -o ./$BIN/$TARGET ./cmd/$TARGET
}

frontend() {
  cd website
  npm run build
}

all() {
  backend
  frontend
}

clean() {
  rm -rf $BIN
  rm -rf ./website/dist
}

case "$1" in
  run)
    if [ -z "$2" ]; then
      echo "No argument after *run* command."
      exit 1
    elif [ "$2" = "back" ]; then
      backend
      ./$BIN/$TARGET
    else
      echo "Unknown application *${2}*"
      exit 1
    fi
    ;;
  back)
    backend
    ;;
  front)
    echo "Compiling frontend into static files..."
    frontend
    ;;
  clean)
    echo "Removing build directories..."
    clean
    ;;
  *)
    if [ -z "$1" ]; then
      BACKEND_DEV=0
      all
    else
      echo "Unknown command *${1}*"
    fi
    ;;
esac
