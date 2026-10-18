#!/usr/bin/env bash
# Demo deploy: proves the pipeline wiring without any cloud. A real one would call your platform's CLI.
set -euo pipefail
: "${DEPLOY_ENVIRONMENT:?set by the FlowForge deploy job}"
if [ -z "${DEMO_TOKEN:-}" ]; then
  echo "::error::DEMO_TOKEN secret is not set for environment $DEPLOY_ENVIRONMENT (set it from Go-live)"
  exit 1
fi
npm run build
echo "Deployed storefront ${GITHUB_SHA:-local} to $DEPLOY_ENVIRONMENT (token length ${#DEMO_TOKEN})"
