#! /bin/bash

build() {
  check_installed "go"

  gen
  build_server
  build_bot
}

build_server() {
  CGO_ENABLED=0 GOOS=linux go build -o bin/scurvy10k-server cmd/server/main.go
  version=$(cat version)
  echo "Building scurvy10k:$version"
  docker buildx build . -f Dockerfile -t ghcr.io/torfstack/scurvy10k:"$version"
  docker push ghcr.io/torfstack/scurvy10k:"$version"
}

build_bot() {
  CGO_ENABLED=0 GOOS=linux go build -o bin/scurvy10k-bot cmd/bot/main.go
  version=$(cat version)
  echo "Building scurvy10k-bot:$version"
  docker buildx build . -f Dockerfile-bot -t ghcr.io/torfstack/scurvy10k-bot:"$version"
  docker push ghcr.io/torfstack/scurvy10k-bot:"$version"
}

gen() {
  check_installed "templ"

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
    deploy)
      deploy
      ;;
    *)
      echo "Usage: do [build|clean|deploy|gen]"
      exit 1
      ;;
  esac
}

deploy() {
  check_installed "helm"
  echo "Deploying..."
  helm upgrade --install scurvy10k deployment --values deployment/values.yaml -f deployment/values.yaml -n default
}

check_installed() {
  if ! command -v "$1" &> /dev/null; then
    echo "$1 is not installed"
    exit 1
  fi
}

start "$@"