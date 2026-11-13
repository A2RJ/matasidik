#!/usr/bin/env bash
# Regenerate the stack-diagram HTML sources and single WebP exports used by README.md and the docs.
#
# Requires:
#   - Python 3.9+ dan Pillow (konversi PNG sementara -> WebP lossless)
#   - stack-diagram clone:  STACK_DIR=/path/to/stack-diagram/skills/stack-diagram
#   - Node + Playwright:    PW_NODE_PATH=/path/to/node_modules (contains playwright)
#   - Google Chrome (the renderer uses channel "chrome")
set -euo pipefail
cd "$(dirname "$0")"

STACK_DIR="${STACK_DIR:-/tmp/opencode/stack-diagram/skills/stack-diagram}"
ICON_DIR="${ICON_DIR:-$PWD/icons}"
PW_NODE_PATH="${PW_NODE_PATH:-/tmp/opencode/pw/node_modules}"

if [[ ! -f "$STACK_DIR/assets/template.html" ]]; then
  echo "stack-diagram not found at: $STACK_DIR" >&2
  echo "Set STACK_DIR to .../stack-diagram/skills/stack-diagram" >&2
  exit 1
fi

python3 "$STACK_DIR/scripts/fetch_icons.py" "$ICON_DIR" \
  "input:lucide/clipboard-list" "verif:lucide/search-check" "valid:lucide/badge-check" \
  "engine:lucide/database" "dash:lucide/layout-dashboard" "keputusan:lucide/landmark" \
  "tl:lucide/list-checks" "draf:lucide/file-pen-line" "kirim:lucide/send" \
  "final:lucide/lock" "ubah:lucide/key-round" "pkm:lucide/hospital" \
  "dinas:lucide/building-2" "kabid:lucide/user-check" "pimpinan:lucide/crown" \
  "perbaiki:lucide/file-pen" "kembali:lucide/undo-2" "faktor:lucide/message-square" \
  "audit:lucide/scroll-text" "target:lucide/target" "table:lucide/table-2" \
  "periode:lucide/calendar" "wilayah:lucide/map" "role:lucide/shield" \
  "laporan:lucide/file-spreadsheet" "user:lucide/user" \
  "upload:lucide/file-up" "klasifikasi:lucide/scan-search" "workflow:lucide/workflow" \
  "calc:lucide/calculator" "check:lucide/circle-check" "loading:lucide/loader-circle" \
  "warning:lucide/triangle-alert" "clock:lucide/clock" "tree:lucide/list-tree" \
  "disk:lucide/hard-drive" "container:lucide/container" "server:lucide/server" \
  "code:lucide/braces" "arrow:lucide/arrow-right" "panduan:lucide/book-open"

export STACK_DIAGRAM_TEMPLATE="$STACK_DIR/assets/template.html"
python3 gen.py
python3 gen-end-to-end.py
python3 gen-teknis.py
python3 gen-schema.py
python3 gen-login.py
python3 gen-panduan.py

build_one() {
  local name="$1" width="$2" height="${3:-}"
  python3 "$STACK_DIR/scripts/build.py" "drafts/$name.html" "$ICON_DIR" "$name.html" \
    --fonts "$STACK_DIR/assets/fonts"
  NODE_PATH="$PW_NODE_PATH" node render.js "$PWD/$name.html" "$PWD/drafts/$name.png" "$width" "$height"
  python3 to-webp.py "$PWD/drafts/$name.png" "$PWD/$name.webp"
  rm -f "$PWD/drafts/$name.png"
}

build_canvas() {
  local name="$1" width="$2" height="$3"
  python3 "$STACK_DIR/scripts/build.py" "drafts/$name.html" "$ICON_DIR" "$name.html" \
    --fonts "$STACK_DIR/assets/fonts"
  NODE_PATH="$PW_NODE_PATH" node render-canvas.js "$PWD/$name.html" "$PWD/drafts/$name.png" "$width" "$height"
  python3 to-webp.py "$PWD/drafts/$name.png" "$PWD/$name.webp"
  rm -f "$PWD/drafts/$name.png"
}

build_one alur-bisnis 1340
build_one status-data 1240
build_one alur-peran 1360
build_one model-data 1320
build_one alur-data-iku 1860
build_one alur-upload 1280
build_one alur-workflow 1360
build_one alur-provenance 1400
build_one alur-trace-spm 1860
build_one arsitektur-aplikasi 1440
build_one schema-db 1920
build_one alur-login 1200 220
build_canvas alur-panduan 1680 440

mkdir -p ../../frontend/static/user-guide
cp alur-panduan.webp ../../frontend/static/user-guide/alur-panduan.webp

# Normalize generated HTML so trailing whitespace does not leak into diffs.
python3 - <<'PY'
from pathlib import Path
for p in sorted(Path('.').glob('*.html')):
    text = p.read_text(encoding='utf-8')
    normalized = '\n'.join(line.rstrip() for line in text.splitlines()) + '\n'
    if normalized != text:
        p.write_text(normalized, encoding='utf-8')
        print(f'normalized {p}')
PY

echo "Done. HTML sources and WebP exports refreshed in $PWD"
