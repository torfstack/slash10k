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

db_migrate() {
  check_installed "atlas"
  if [ -z "$2" ]; then
    echo "Usage: do db_migrate <name of migration> "
    exit 1
  fi
  echo "generating migration if necessary"
  atlas migrate diff $2 --to "file://sql/schema.sql" --dev-url "docker://postgres?search_path=public" --format '{{ sql . "  " }}' --dir "file://sql/migrations"
}

db_apply() {
  check_installed "atlas"
  echo "applying migrations"
  pw=$(kubectl get secret scurvy10k -o jsonpath="{.data.db-password}" | base64 --decode)
  atlas migrate apply --url "postgres://scurvy10k:$pw@localhost:5432/scurvy10k?search_path=public&sslmode=disable" --dir "file://sql/migrations"
}

db_status() {
  check_installed "atlas"
  echo "checking migration status"
  pw=$(kubectl get secret scurvy10k -o jsonpath="{.data.db-password}" | base64 --decode)
  atlas migrate status --url "postgres://scurvy10k:$pw@localhost:5432/scurvy10k?search_path=public&sslmode=disable"
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
    db_apply)
      db_apply
      ;;
    db_migrate)
      db_migrate "$@"
      ;;
    db_status)
      db_status
      ;;
    deploy)
      deploy
      ;;
    *)
      echo "Usage: do [build|clean|db_apply|db_migrate|deploy|gen]"
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