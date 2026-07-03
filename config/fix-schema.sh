#!/bin/sh
# Rewrites OpenAPI 3.1 "exclusiveMinimum": <number> to OpenAPI 3.0 compatible
# "minimum": 1 so ogen can parse the schema without choking on the int-vs-bool mismatch.
set -e
FILE="$1"
jq 'walk(if type == "object" and (.exclusiveMinimum // null | type) == "number" then del(.exclusiveMinimum) | .minimum = 1 else . end)' "$FILE" > "${FILE}.tmp" && mv "${FILE}.tmp" "$FILE"
