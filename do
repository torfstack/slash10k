#! /bin/bash

build() {
  check_installed "go"

  gen
  build_server
  build_bot
}

build_server() {
  CGO_ENABLED=0 GOOS=linux go build -o bin/slash10k-server cmd/server/main.go
  version=$(cat version)
  echo "Building slash10k:$version"
  docker buildx build . -f Dockerfile -t ghcr.io/torfstack/slash10k:"$version"
  docker push ghcr.io/torfstack/slash10k:"$version"
}

build_bot() {
  CGO_ENABLED=0 GOOS=linux go build -o bin/slash10k-bot cmd/bot/main.go
  version=$(cat version)
  echo "Building slash10k-bot:$version"
  docker buildx build . -f Dockerfile-bot -t ghcr.io/torfstack/slash10k-bot:"$version"
  docker push ghcr.io/torfstack/slash10k-bot:"$version"
}

run() {
  check_installed "air"
  echo "Running..."
  DATABASE_CONNECTION_HOST=localhost \
    DATABASE_CONNECTION_PORT=5432 \
    DATABASE_CONNECTION_USER=postgres \
    DATABASE_CONNECTION_PASSWORD=mysecretpassword \
    DATABASE_CONNECTION_DBNAME=slash10k air
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
    run)
      run
      ;;
    *)
      echo "Usage: do [build|clean|db_apply|db_migrate|db_status|deploy|gen|run]"
      exit 1
      ;;
  esac
}

deploy() {
  check_installed "helm"
  echo "Deploying..."
  helm upgrade --install slash10k deployment --values deployment/values.yaml -f deployment/values.yaml -n default
}

check_installed() {
  if ! command -v "$1" &> /dev/null; then
    echo "$1 is not installed"
    exit 1
  fi
}

start "$@"