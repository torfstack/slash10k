#! /bin/bash

build() {
  templ generate
  CGO_ENABLED=0 GOOS=linux go build -o bin/scurvy10k-backend src/main.go
  docker buildx build . -t scurvy10k-backend
}

clean() {
  echo "Cleaning up..."
  echo "Removing bin/ and templ/.go"
  rm templ/*.go &> /dev/null
  rm -r bin &> /dev/null
}

start() {
  check_installed "templ"
  check_installed "go"

  case "$1" in
    build)
      build
      ;;
    clean)
      clean
      ;;
    *)
      echo "Usage: do [build|clean]"
      exit 1
      ;;
  esac
}

check_installed() {
  if ! command -v "$1" &> /dev/null; then
    echo "$1 is not installed"
    exit 1
  fi
}

start "$@"