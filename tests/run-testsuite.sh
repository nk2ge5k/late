#!/bin/bash

export ROOT_DIR="$1"
PYTHONPATH="$ROOT_DIR:$PYTHONPATH"

export TZ=Europe/Moscow
export PYTHONPATH

TEST_DIR=$(realpath "$ROOT_DIR/tests")
$PYTHON -m pytest -vvv "$PYTEST_ARGS" "$TEST_DIR" "$TEST_DIR"
