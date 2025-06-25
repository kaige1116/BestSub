#!/bin/sh
set -e

chmod +x /app/bestsub

exec "cd /app && ./bestsub" "$@" 