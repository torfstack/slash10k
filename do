#! /bin/bash

build() {
  check_installed "go"

  gen
  CGO_ENABLED=0 GOOS=linux go build -o bin/slash10k cmd/bot/main.go
  version=$(cat version)
  echo "Building slash10k:$version"
  docker buildx build . -f Dockerfile -t ghcr.io/torfstack/slash10k:"$version"
  docker push ghcr.io/torfstack/slash10k:"$version"
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
    deploy_dev)
      deploy_dev
      ;;
    deploy_prod)
      deploy_prod
      ;;
    *)
      echo "Usage: do [build|clean|deploy_dev|deploy_prod|gen]"
      exit 1
      ;;
  esac
}

deploy_dev() {
  check_installed "helm"
  echo "Deploying DEV ..."
  helm upgrade --install slash10kdev deployment --values deployment/values-dev.yaml -f deployment/values-dev.yaml -n default
}

deploy_prod() {
  check_installed "helm"
  echo "Deploying PROD ..."
  helm upgrade --install slash10k deployment --values deployment/values-prod.yaml -f deployment/values-prod.yaml -n default
}

check_installed() {
  if ! command -v "$1" &> /dev/null; then
    echo "$1 is not installed"
    exit 1
  fi
}

start "$@"