#!/usr/bin/env bash
set -e

# Upload release binaries from dist/ to Supabase Storage bucket
# Required environment variables:
# - SUPABASE_URL (e.g. https://your-project.supabase.co)
# - SUPABASE_SERVICE_ROLE_KEY or SUPABASE_ANON_KEY
# - SUPABASE_BUCKET (defaults to "releases")

SUPABASE_URL="${SUPABASE_URL:-}"
SUPABASE_KEY="${SUPABASE_SERVICE_ROLE_KEY:-${SUPABASE_ANON_KEY:-}}"
SUPABASE_BUCKET="${SUPABASE_BUCKET:-releases}"
DIST_DIR="${1:-dist}"

if [ -z "${SUPABASE_URL}" ] || [ -z "${SUPABASE_KEY}" ]; then
    echo "[!] Warning: SUPABASE_URL or SUPABASE_KEY not set. Skipping Supabase Storage upload."
    exit 0
fi

echo "==> Uploading release artifacts from ./${DIST_DIR}/ to Supabase Storage bucket '${SUPABASE_BUCKET}'..."

# Ensure bucket exists
curl -s -X POST "${SUPABASE_URL}/storage/v1/bucket" \
  -H "Authorization: Bearer ${SUPABASE_KEY}" \
  -H "apikey: ${SUPABASE_KEY}" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"${SUPABASE_BUCKET}\", \"name\": \"${SUPABASE_BUCKET}\", \"public\": true}" >/dev/null 2>&1 || true

for FILE in "${DIST_DIR}"/*; do
    if [ -f "${FILE}" ]; then
        FILENAME="$(basename "${FILE}")"
        echo "  --> Uploading ${FILENAME}..."
        
        UPLOAD_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
            "${SUPABASE_URL}/storage/v1/object/${SUPABASE_BUCKET}/${FILENAME}" \
            -H "Authorization: Bearer ${SUPABASE_KEY}" \
            -H "apikey: ${SUPABASE_KEY}" \
            -H "Content-Type: application/octet-stream" \
            -H "x-upsert: true" \
            --data-binary "@${FILE}")
        
        if [ "${UPLOAD_STATUS}" = "200" ] || [ "${UPLOAD_STATUS}" = "201" ]; then
            echo "      [+] Successfully uploaded ${FILENAME} (HTTP ${UPLOAD_STATUS})"
            echo "      CDN URL: ${SUPABASE_URL}/storage/v1/object/public/${SUPABASE_BUCKET}/${FILENAME}"
        else
            echo "      [!] Upload returned status ${UPLOAD_STATUS} for ${FILENAME}"
        fi
    fi
done

echo "==> Supabase Storage release upload complete."
