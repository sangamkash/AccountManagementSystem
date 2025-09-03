#!/bin/sh
set -e

echo "=================Checking goose..."
which goose || echo "goose not found"
goose --version || echo "goose version check failed"
echo "==================Running migrations..."
goose up

echo "===================Starting app..."
exec "$@"