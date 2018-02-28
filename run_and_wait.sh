#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

execute_task() {
  set -e -o pipefail
  local task_hash="$1"

  computes-cli task enqueue "$task_hash"
}

get_dataset_hash() {
  set -e -o pipefail
  local filename="$1"

  cat "$filename" \
  | "$DIR/split-map-reduce-to-task" \
  | head -n 3 \
  | tail -n 1 \
  | tr -d '\n'
}

get_task_hash() {
  set -e -o pipefail
  local filename="$1"

  cat "$filename" \
  | "$DIR/split-map-reduce-to-task" \
  | tail -n 1 \
  | tr -d '\n'
}

dumplatest() {
  local dataset_hash
  set -e -o pipefail

  dataset_hash="$1"
  computes-cli dataset dumplatest "$dataset_hash"
}

wait_for_completion() {
  set -e -o pipefail
  local dataset_hash latest

  dataset_hash="$1"
  latest="null"

  while [ "$latest" == "null" ]; do
    sleep 1
    latest="$(dumplatest "$dataset_hash" | jq '.reduce.results')"
  done
  echo -n "$latest"
}

execute_and_wait() {
  local task_hash dataset_hash result
  set -e -o pipefail

  task_hash="$1"
  dataset_hash="$2"

  execute_task "$task_hash"
  result="$(wait_for_completion "$dataset_hash")"
  echo ""
  echo "result: $result"
}

main() {
  local filename task_hash
  set -e -o pipefail

  filename="$1"
  task_hash="$(get_task_hash "$filename")"
  dataset_hash="$(get_dataset_hash "$filename")"

  time execute_and_wait "$task_hash" "$dataset_hash"
}
main "$@"
