#!/bin/sh
# DHC/Railway: ClickClack starten und täglich nach Cloudflare R2 sichern.
# R2-Zugang kommt aus RCLONE_CONFIG_R2_* (Railway-Variablen), Ziel aus BACKUP_BUCKET.
set -eu
DATA=/app/data

backup() {
  stamp=$(date -u +%Y-%m-%dT%H%M%SZ)
  file="/tmp/clickclack-$stamp.db"
  clickclack backup --data "$DATA" --out "$file"
  gzip "$file"
  rclone copyto "$file.gz" "r2:$BACKUP_BUCKET/db/clickclack-$stamp.db.gz"
  rm -f "$file.gz"
  # ponytail: copy statt sync, gelöschte Uploads bleiben im Backup; Löschwege bei Bedarf über den Bucket.
  if [ -d "$DATA/uploads" ]; then rclone copy "$DATA/uploads" "r2:$BACKUP_BUCKET/uploads"; fi
  echo "railway-backup ok $stamp"
}

if [ "${1:-}" = "backup-now" ]; then
  backup
  exit
fi

if [ -n "${BACKUP_BUCKET:-}" ]; then
  (
    sleep 300
    while true; do
      backup || echo "railway-backup FAILED" >&2
      sleep "${BACKUP_INTERVAL_SECONDS:-86400}"
    done
  ) &
fi

exec clickclack serve --addr ":${PORT:-8080}" --data "$DATA"
