#! /bin/bash

build() {
  gen
  CGO_ENABLED=0 GOOS=linux go build -o bin/scurvy10k-backend src/main.go
  docker buildx build . -t scurvy10k-backend
}

gen() {
  echo "Generating templ..."
  templ generate
  echo "Generating sql..."
  sqlc generate
}

clean() {
  echo "Cleaning up..."
  echo "templ/.go"
  rm templ/*.go &> /dev/null
  echo "bin"
  rm -r bin &> /dev/null
  echo "sqlc"
  rm -r sql/db &> /dev/null
}

start() {
  check_installed "templ"
  check_installed "go"

  case "$1" in
    build)
      build
      ;;
    gen)
      gen
      ;;
    clean)
      clean
      ;;
    *)
      echo "Usage: do [build|clean|gen]"
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