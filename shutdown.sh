#!/usr/bin/env bash
# Hentikan proses pengembangan MATA SIDIK yang berjalan dari repositori ini.
# Script ini tidak mematikan sistem operasi, database eksternal, IDE, atau agent.
set -uo pipefail

ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
BACKEND="$ROOT/backend"
FRONTEND="$ROOT/frontend"
DEV_PORTS=(5173 8080 18084 18085 18086 18087 18088 18089 18090 18091 18092 18093 18094 18095)
PIDS=()

add_pid() {
  local pid="$1"
  [[ "$pid" =~ ^[0-9]+$ ]] || return
  [[ "$pid" == "$$" ]] && return
  for existing in "${PIDS[@]:-}"; do
    [[ "$existing" == "$pid" ]] && return
  done
  PIDS+=("$pid")
}

pid_cwd() {
  readlink -f "/proc/$1/cwd" 2>/dev/null || true
}

pid_cmd() {
  tr '\0' ' ' < "/proc/$1/cmdline" 2>/dev/null || true
}

# 1. Proses yang sedang listen di port pengembangan aplikasi.
for port in "${DEV_PORTS[@]}"; do
  while read -r pid; do
    add_pid "$pid"
  done < <(lsof -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null || true)
done

# 2. Proses Go/Vite/pnpm yang cwd-nya berada di repositori ini.
while read -r pid; do
  cwd="$(pid_cwd "$pid")"
  cmd="$(pid_cmd "$pid")"
  case "$cwd" in
    "$BACKEND"|"$BACKEND"/*)
      case "$cmd" in
        *"go run ./cmd/api"*|*"go run ./cmd/worker"*|*"/exe/api"*|*"/api "*) add_pid "$pid" ;;
      esac
      ;;
    "$FRONTEND"|"$FRONTEND"/*)
      case "$cmd" in
        *"vite"*|*"pnpm dev"*|*"pnpm-native dev"*|*"npm run dev"*|*"esbuild"*) add_pid "$pid" ;;
      esac
      ;;
  esac
done < <(pgrep -f 'go run|vite|pnpm|npm|esbuild|/api' 2>/dev/null || true)

if [[ ${#PIDS[@]} -eq 0 ]]; then
  echo "Tidak ada proses MATA SIDIK yang berjalan."
  exit 0
fi

echo "Menghentikan proses MATA SIDIK: ${PIDS[*]}"
kill -TERM "${PIDS[@]}" 2>/dev/null || true

for _ in $(seq 1 20); do
  alive=()
  for pid in "${PIDS[@]}"; do
    kill -0 "$pid" 2>/dev/null && alive+=("$pid")
  done
  [[ ${#alive[@]} -eq 0 ]] && break
  sleep 0.25
done

if [[ ${#alive[@]} -gt 0 ]]; then
  echo "Memaksa berhenti: ${alive[*]}"
  kill -KILL "${alive[@]}" 2>/dev/null || true
fi

echo "MATA SIDIK dihentikan."
