#!/usr/bin/env bash
# Console first-boot entry: launch the Equate setup wizard on tty1 until complete.
set -euo pipefail

DEPLOY_FILE="/etc/equate/deploy-dir"
if [[ ! -f "${DEPLOY_FILE}" ]]; then
  echo "Equate deploy directory is not configured (${DEPLOY_FILE} missing)." >&2
  exit 1
fi

DEPLOY_DIR="$(tr -d '[:space:]' < "${DEPLOY_FILE}")"
if [[ -z "${DEPLOY_DIR}" || ! -d "${DEPLOY_DIR}" ]]; then
  echo "Equate deploy directory is invalid: ${DEPLOY_DIR}" >&2
  exit 1
fi

if [[ -f "${DEPLOY_DIR}/.setup-complete" ]]; then
  exit 0
fi

cd "${DEPLOY_DIR}"

if command -v equate >/dev/null 2>&1; then
  exec equate configure
fi

echo "equate CLI not found on PATH; install the appliance release first." >&2
exit 1
