#!/usr/bin/env bash
set -euo pipefail

echo "==================================================="
echo "    Running CloudStoreX Full Verification Suite    "
echo "==================================================="

echo ""
echo "--> [1/2] Running Backend Go Test Suite..."
cd backend
go test -v -cover ./...
cd ..

echo ""
echo "--> [2/2] Running Frontend Vitest Unit Suite..."
cd frontend
npm test -- --run
cd ..

echo ""
echo "==================================================="
echo "          All CloudStoreX Tests Passed             "
echo "==================================================="
