#!/usr/bin/env bash
set -euo pipefail

echo "This is a placeholder test runner."

# from the handout - run in the /inklings directory
(time timeout -k 1 15 bash -c "make run 20 20 40 >> stdout.txt") 2>&1 | tee -a all.txt