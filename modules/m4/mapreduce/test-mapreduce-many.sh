#!/usr/bin/env bash

#
# map-reduce stability tests
#

if [ $# -ne 1 ]; then
    echo "Usage: $0 numTrials"
    exit 1
fi

runs=$1
chmod +x test-mapreduce.sh

for i in $(seq 1 "$runs"); do
    echo "*** Trial $i of $runs"

    # ensure clean environment between trials
    rm -rf autograder-logs-apps socket || true
    mkdir -p autograder-logs-apps socket

    # run one full test suite, forcing fast coordinator polling
    if command -v timeout >/dev/null 2>&1; then
        POLL_MS=200 timeout -k 2s 600s ./test-mapreduce.sh
    else
        POLL_MS=200 ./test-mapreduce.sh
    fi

    if [ $? -ne 0 ]; then
        echo "*** FAILED TESTS IN TRIAL $i"
        exit 1
    fi

    echo "--- Trial $i complete"
done

echo "*** PASSED ALL $runs TESTING TRIALS"
