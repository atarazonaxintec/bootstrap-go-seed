#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 || ! "$1" =~ ^[a-z][a-z0-9]*$ ]]; then
  echo "Uso: $0 <nombre-del-modulo> (solo letras minúsculas y números)" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
template="$script_dir/internal/health"
target="$script_dir/internal/$1"
module="$1"
module_type="$(printf '%s' "$module" | awk '{ print toupper(substr($0, 1, 1)) substr($0, 2) }')"

if [[ ! -d "$template" ]]; then
  echo "No se encontró la plantilla: $template" >&2
  exit 1
fi

if [[ -e "$target" ]]; then
  echo "El módulo ya existe: $target" >&2
  exit 1
fi

cp -R "$template" "$target"
find "$target" -type f -name '*.go' \
  -exec sed -i '' -e "s/Health/$module_type/g" -e "s/health/$module/g" {} +

while IFS= read -r -d '' file; do
  mv "$file" "${file//health/$module}"
done < <(find "$target" -type f -name '*health*' -print0)