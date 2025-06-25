#!/bin/sh
set -e

chmod +x /app/bestsub

cd /app && exec ./bestsub "$@"