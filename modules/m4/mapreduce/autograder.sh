#!/usr/bin/env bash
set -euo pipefail

LOG_DIR="autograder-logs"
rm -rf "$LOG_DIR" 2>/dev/null || true
mkdir -p "$LOG_DIR"

# permissions to execute
chmod +x test-mapreduce.sh
chmod +x test-mapreduce-many.sh 

echo "=== Part [1/4] Build ==="
if ! make build >"$LOG_DIR/build.log" 2>&1; then
    echo "[FAIL] build failed. See $LOG_DIR/build.log"
    exit 1
fi
echo "[OK] build"

echo "=== Part [2/4] Functional tests (test-mapreduce.sh) ==="
TEST_FAILED=0

# run the main test suite with a faster coordinator poll
if command -v timeout >/dev/null 2>&1; then
    if ! POLL_MS=200 timeout -k 2s 900s ./test-mapreduce.sh all >"$LOG_DIR/test.log" 2>&1; then
        TEST_FAILED=1
        echo "[FAIL] test-mapreduce.sh exited nonzero. See $LOG_DIR/test.log"
    else
        echo "[OK] test-mapreduce.sh"
    fi
else
    if ! POLL_MS=200 ./test-mapreduce.sh >"$LOG_DIR/test.log" 2>&1; then
        TEST_FAILED=1
        echo "[FAIL] test-mapreduce.sh exited nonzero. See $LOG_DIR/test.log"
    else
        echo "[OK] test-mapreduce.sh"
    fi
fi

echo "=== Part [3/4] Stability (test-mapreduce-many.sh) ==="
MANY_FAILED=0
if [ -x ./test-mapreduce-many.sh ]; then
    if command -v timeout >/dev/null 2>&1; then
        if ! POLL_MS=200 timeout -k 2s 900s ./test-mapreduce-many.sh 3 >"$LOG_DIR/test-many.log" 2>&1; then
            MANY_FAILED=1
            echo "[WARN] test-mapreduce-many.sh reported failure. See $LOG_DIR/test-many.log"
        else
            echo "[OK] test-mapreduce-many.sh"
        fi
    else
        if ! POLL_MS=200 ./test-mapreduce-many.sh 3 >"$LOG_DIR/test-many.log" 2>&1; then
            MANY_FAILED=1
            echo "[WARN] test-mapreduce-many.sh reported failure. See $LOG_DIR/test-many.log"
        else
            echo "[OK] test-mapreduce-many.sh"
        fi
    fi
else
    echo "[SKIP] test-mapreduce-many.sh not present/executable"
fi

echo "=== Part [4/4] Log analysis ==="
WC_PASS=$(grep -c 'wc test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
INDEX_PASS=$(grep -c 'indexer test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
MAP_PAR_PASS=$(grep -c 'map parallelism test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
RED_PAR_PASS=$(grep -c 'reduce parallelism test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
EARLY_PASS=$(grep -c 'early exit test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
CRASH_PASS=$(grep -c 'crash test: PASS' "$LOG_DIR/test.log" 2>/dev/null || true)
DATA_RACES=$(grep -c 'DATA RACE' "$LOG_DIR/test.log" 2>/dev/null || true)

echo "testing summary:"
echo "  wc:                 $WC_PASS"
echo "  indexer:            $INDEX_PASS"
echo "  map parallelism:    $MAP_PAR_PASS"
echo "  reduce parallelism: $RED_PAR_PASS"
echo "  early exit:         $EARLY_PASS"
echo "  crash:              $CRASH_PASS"
echo "  data races:         $DATA_RACES"
echo "  many-tests failed:  $MANY_FAILED"

EXIT_CODE=0

# main test failed / fail
if [ "$TEST_FAILED" -ne 0 ]; then
    EXIT_CODE=1
fi

# core tests didn’t pass / fail
if [ "$WC_PASS" -eq 0 ] || [ "$INDEX_PASS" -eq 0 ] || [ "$CRASH_PASS" -eq 0 ]; then
    EXIT_CODE=1
fi

# make data races fatal if you want too:
# if [ "$DATA_RACES" -gt 0 ]; then
#     EXIT_CODE=1
# fi

exit $EXIT_CODE
