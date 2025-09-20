#!/usr/bin/env bash
# Copies each *.txt file in the same directory as this script to the clipboard,
# pausing for a key press between files. Press 'q' to quit.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Determine clipboard command
clip_cmd=""
if command -v pbcopy >/dev/null 2>&1; then
  clip_cmd="pbcopy"
elif command -v wl-copy >/dev/null 2>&1; then
  clip_cmd="wl-copy"
elif command -v xclip >/dev/null 2>&1; then
  clip_cmd="xclip -selection clipboard"
elif command -v xsel >/dev/null 2>&1; then
  clip_cmd="xsel --clipboard --input"
else
  echo "No clipboard utility found. Install one of: pbcopy, wl-copy, xclip, xsel" >&2
  exit 1
fi

# Collect non-recursive .txt files in script directory
mapfile -t files < <(find "$script_dir" -maxdepth 1 -type f -name "*.txt" | sort)

if [[ ${#files[@]} -eq 0 ]]; then
  echo "No .txt files found in: $script_dir"
  exit 0
fi

echo "Found ${#files[@]} .txt files in: $script_dir"
echo "Press any key to copy next file to clipboard, or 'q' to quit."

for f in "${files[@]}"; do
  if ! cat "$f" | eval "$clip_cmd"; then
    echo "Failed to copy: $f" >&2
    continue
  fi
  size=$(wc -c < "$f" | tr -d ' ')
  base=$(basename "$f")
  echo
  echo "Copied to clipboard: $base (${size} bytes)"
  echo -n "[Any key = next, q = quit] "

  IFS= read -rsn1 key
  echo
  if [[ "${key:-}" == "q" || "${key:-}" == "Q" ]]; then
    echo "Quitting."
    exit 0
  fi
done

echo "Done."