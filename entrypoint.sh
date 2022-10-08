#!/bin/sh
set -e

export DB_HOST=betera_db
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=12345
export DB_NAME=postgres
export HOST=0.0.0.0


exec betera