#!/usr/bin/env bash
set -euo pipefail

# Sync provider docs to docs/vendor/<provider>/<YYYY-MM-DD>/
# Requirements: curl, pandoc (for HTML->GFM). Leaves raw HTML + markdown.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$ROOT/docs/sources"
OUT_DIR="$ROOT/docs/vendor"
DATE="$(date -u +%Y-%m-%d)"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    return 1
  fi
}

need_cmd curl
if ! command -v pandoc >/dev/null 2>&1; then
  echo "Warning: pandoc not found. Will save HTML only." >&2
  PANDOC=false
else
  PANDOC=true
fi

mkdir -p "$OUT_DIR"

slug() {
  echo "$1" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//'
}

for file in "$SRC_DIR"/*.txt; do
  [ -e "$file" ] || continue
  provider="$(basename "$file" .txt)"
  target_dir="$OUT_DIR/$provider/$DATE"
  mkdir -p "$target_dir"

  mapfile -t urls < <(grep -Ev '^\s*(#|$)' "$file")
  meta="$target_dir/_sources.txt"
  printf "provider=%s\nfetched_at=%s\n" "$provider" "$DATE" >"$meta"
  printf "urls:\n" >>"$meta"
  for u in "${urls[@]}"; do
    printf "  - %s\n" "$u" >>"$meta"
  done

  for url in "${urls[@]}"; do
    fname="$(slug "$url")"
    html="$target_dir/$fname.html"
    md="$target_dir/$fname.md"
    echo "Fetching $url -> $html"
    curl -L --fail --compressed "$url" -o "$html"
    if [ "$PANDOC" = true ]; then
      pandoc "$html" -f html -t gfm --wrap=none -o "$md"
    fi
  done

  # Update latest symlink for the provider
  latest="$OUT_DIR/$provider/latest"
  rm -f "$latest"
  ln -s "$DATE" "$latest"
done

echo "Done. Docs stored under $OUT_DIR/<provider>/<date>/"
