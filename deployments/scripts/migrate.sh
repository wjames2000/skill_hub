#!/bin/bash
set -euo pipefail

# Database Migration Script for SkillHub Pro
# Usage: ./migrate.sh [environment]
#   environment: local | staging | production (default: local)

ENV="${1:-local}"
CONFIG_FILE="../backend/configs/config.${ENV}.yaml"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "Error: Config file not found: $CONFIG_FILE"
  exit 1
fi

echo "Running database migrations for environment: $ENV"

# Extract DSN from config using grep (requires yq for proper parsing)
if command -v yq &> /dev/null; then
  DSN=$(yq eval '.db.dsn' "$CONFIG_FILE")
elif command -v python3 &> /dev/null; then
  DSN=$(python3 -c "
import yaml
with open('$CONFIG_FILE') as f:
    cfg = yaml.safe_load(f)
print(cfg['db']['dsn'])
" 2>/dev/null || echo "")
fi

if [ -z "${DSN:-}" ]; then
  echo "Warning: Could not parse DSN, using default DSN from config"
  echo "Please ensure database is accessible and run migrations manually:"
fi

MIGRATIONS_DIR="../backend/migrations"

echo "Migration files found:"
ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null || echo "  (no migration files found)"

echo ""
echo "To execute migrations, run:"
for f in "$MIGRATIONS_DIR"/*.sql; do
  [ -f "$f" ] || continue
  filename=$(basename "$f")
  echo "  psql \"\$DSN\" -f \"$f\""
done

echo ""
echo "Example (using Docker):"
echo "  docker exec -i skill-hub-postgres psql -U skillhub -d skillhub < $MIGRATIONS_DIR/001_init_sync_tables.sql"
echo "  docker exec -i skill-hub-postgres psql -U skillhub -d skillhub < $MIGRATIONS_DIR/002_vector_router_tables.sql"

# Auto-execute for local environment via Docker
if [ "$ENV" = "local" ]; then
  echo ""
  echo "Auto-running migrations for local environment..."
  for f in "$MIGRATIONS_DIR"/*.sql; do
    [ -f "$f" ] || continue
    filename=$(basename "$f")
    echo "  Executing $filename..."
    docker exec -i skill-hub-postgres-1 psql -U skillhub -d skillhub < "$f" 2>/dev/null && \
      echo "  ✓ $filename applied" || \
      echo "  ⚠ Could not apply $filename (container may not be running)"
  done
fi

echo ""
echo "Migration check completed."
