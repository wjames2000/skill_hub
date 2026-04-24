#!/bin/bash
set -euo pipefail

# Deployment Verification Script for SkillHub Pro
# Usage: ./verify-deployment.sh [namespace]
NAMESPACE="${1:-skill-hub}"

echo "=========================================="
echo "  SkillHub Pro - Deployment Verification"
echo "=========================================="
echo ""

# 1. Check namespace exists
if kubectl get namespace "$NAMESPACE" &>/dev/null; then
  echo "✓ Namespace '$NAMESPACE' exists"
else
  echo "✗ Namespace '$NAMESPACE' not found"
  exit 1
fi

# 2. Check all pods are running
echo ""
echo "--- Pod Status ---"
PODS=$(kubectl get pods -n "$NAMESPACE" -o json)
if echo "$PODS" | jq -e '.items | length > 0' &>/dev/null; then
  NOT_READY=$(echo "$PODS" | jq -r '[.items[] | select(.status.phase != "Running" or ([.status.conditions[] | select(.type == "Ready" and .status == "True")] | length) == 0)] | length')
  TOTAL=$(echo "$PODS" | jq '.items | length')
  kubectl get pods -n "$NAMESPACE"
  if [ "$NOT_READY" -gt 0 ]; then
    echo "⚠ $NOT_READY/$TOTAL pods are not ready"
  else
    echo "✓ All $TOTAL pods are running and ready"
  fi
else
  echo "✗ No pods found in namespace '$NAMESPACE'"
  exit 1
fi

# 3. Check services
echo ""
echo "--- Service Status ---"
kubectl get svc -n "$NAMESPACE" -o name | while read -r svc; do
  ENDPOINTS=$(kubectl get "$svc" -n "$NAMESPACE" -o json | jq -r '.spec.clusterIP // "None"')
  echo "  ✓ $svc (ClusterIP: $ENDPOINTS)"
done

# 4. Check ingress
echo ""
echo "--- Ingress Status ---"
if kubectl get ingress -n "$NAMESPACE" &>/dev/null; then
  kubectl get ingress -n "$NAMESPACE"
else
  echo "  ⚠ No ingress configured"
fi

# 5. Health check endpoints
echo ""
echo "--- Health Checks ---"
SERVICES=("router-api" "admin-api")
for svc in "${SERVICES[@]}"; do
  POD=$(kubectl get pod -n "$NAMESPACE" -l "app=$svc" -o jsonpath="{.items[0].metadata.name}" 2>/dev/null || true)
  if [ -n "$POD" ]; then
    echo "  → Checking $svc ($POD)..."
    kubectl exec "$POD" -n "$NAMESPACE" -- wget -q -O- http://localhost:8080/health 2>/dev/null && \
      echo "  ✓ $svc health check passed" || \
      echo "  ✗ $svc health check failed"
  else
    echo "  ⚠ No pod found for $svc"
  fi
done

# 6. Database connectivity
echo ""
echo "--- Database Check ---"
DB_POD=$(kubectl get pod -n "$NAMESPACE" -l "app=postgres" -o jsonpath="{.items[0].metadata.name}" 2>/dev/null || true)
if [ -n "$DB_POD" ]; then
  kubectl exec "$DB_POD" -n "$NAMESPACE" -- pg_isready -U skillhub -d skillhub &>/dev/null && \
    echo "  ✓ PostgreSQL is ready" || \
    echo "  ✗ PostgreSQL is not ready"
else
  echo "  ⚠ No PostgreSQL pod found in namespace"
fi

# 7. Redis connectivity
echo ""
echo "--- Redis Check ---"
REDIS_POD=$(kubectl get pod -n "$NAMESPACE" -l "app=redis" -o jsonpath="{.items[0].metadata.name}" 2>/dev/null || true)
if [ -n "$REDIS_POD" ]; then
  kubectl exec "$REDIS_POD" -n "$NAMESPACE" -- redis-cli ping &>/dev/null && \
    echo "  ✓ Redis is ready" || \
    echo "  ✗ Redis is not ready"
else
  echo "  ⚠ No Redis pod found in namespace"
fi

echo ""
echo "=========================================="
echo "  Verification Complete"
echo "=========================================="
