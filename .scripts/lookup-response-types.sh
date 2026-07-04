#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <InterfaceName>"
  echo "example: $0 FileListRes"
  exit 1
fi

iface="$1"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
interfaces_file="$root_dir/pkg/client/oas_interfaces_gen.go"
schemas_file="$root_dir/pkg/client/oas_schemas_gen.go"
decoders_file="$root_dir/pkg/client/oas_response_decoders_gen.go"

if [[ ! -f "$interfaces_file" || ! -f "$schemas_file" ]]; then
  echo "generated client files not found under $root_dir/pkg/client"
  exit 1
fi

marker_method="$(
  awk -v iface="$iface" '
    $0 ~ "^type " iface " interface \\{" { in_iface = 1; next }
    in_iface && $0 ~ "^}" { in_iface = 0 }
    in_iface && $0 ~ /\(\)/ {
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", $0)
      sub(/\(\).*/, "", $0)
      print $0
      exit
    }
  ' "$interfaces_file"
)"

if [[ -z "$marker_method" ]]; then
  echo "interface '$iface' not found, or marker method missing in oas_interfaces_gen.go"
  exit 1
fi

echo "Interface: $iface"
echo "Marker method: $marker_method()"
echo
echo "Concrete implementors:"

impl_lines="$(
  rg "func \(\*[^)]*\) ${marker_method}\(\)" "$schemas_file" || true
)"

if [[ -z "$impl_lines" ]]; then
  echo "  (none found)"
  exit 0
fi

while IFS= read -r line; do
  type_name="$(sed -E 's/.*func \(\*([^)]*)\).*/\1/' <<<"$line")"
  if [[ "$type_name" == *Error ]]; then
    echo "  - $type_name (error variant)"
  else
    echo "  - $type_name (non-error variant)"
  fi
done <<<"$impl_lines"

op_name="${iface%Res}"
echo
echo "Decoder entrypoint:"
rg "func decode${op_name}Response" "$decoders_file" || true
