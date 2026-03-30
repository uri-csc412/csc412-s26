#!/usr/bin/env bash

#
# map-reduce tests (tolerant + selectable)
#
# usage:
#   ./test-mapreduce.sh               # runs all tests
#   ./test-mapreduce.sh wc            # just wc
#   ./test-mapreduce.sh indexer       # just indexer
#   ./test-mapreduce.sh mappar        # just map parallelism
#   ./test-mapreduce.sh reducepar     # just reduce parallelism
#   ./test-mapreduce.sh jobcount      # just jobcount
#   ./test-mapreduce.sh early         # just early-exit
#   ./test-mapreduce.sh crash         # just crash test
#

set -euo pipefail

# which tests to run
if [ "$#" -eq 0 ]; then
  TESTS=(wc indexer mappar reducepar jobcount early crash)
else
  TESTS=("$@")
fi

# how fast the coordinator/worker should poll (default 500ms)
: "${POLL_MS:=500}"

# find a timeout command if we have one
TIMEOUT=timeout
if timeout 2s sleep 1 > /dev/null 2>&1; then
  :
else
  if gtimeout 2s sleep 1 > /dev/null 2>&1; then
    TIMEOUT=gtimeout
  else
    TIMEOUT=
    echo '*** Cannot find timeout command; proceeding without timeouts.'
  fi
fi
if [ -n "${TIMEOUT}" ]; then
  TIMEOUT+=" -k 2s 180s "
fi

failed_any=0

########################################
# helper: run coord + N workers and wait
run_cluster () {
  local jobname="$1"
  local nworkers="${2:-3}"

  POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreducecoordinator ../data/pg-*.txt &
  local coord_pid=$!

  # give the coordinator time to create socket
  sleep 0.5

  local pids=()
  for i in $(seq 1 "$nworkers"); do
    POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker "$jobname" &
    pids+=($!)
  done

  # wait for coordinator — might exit nonzero
  if ! wait "$coord_pid"; then
    echo "    (cluster for $jobname exited nonzero)"
  fi

  # kill leftover workers
  for p in "${pids[@]}"; do
    kill "$p" 2>/dev/null || true
  done
}

########################################
# workspace for outputs
mkdir -p autograder-logs-apps
rm -rf autograder-logs-apps/* 2>/dev/null || true
cd autograder-logs-apps

########################################
run_wc () {
  echo '*** Starting wc test.'
  if ../bin/mapreducesequential wc ../data/pg-*.txt; then
    sort mapreduce-out-0 > mapreduce-correct-wc.txt
    rm -f mapreduce-out*
    run_cluster wc 3
    sort mapreduce-out* 2>/dev/null | grep . > mapreduce-wc-all || true
    if cmp -s mapreduce-wc-all mapreduce-correct-wc.txt; then
      echo '--- wc test: PASS'
    else
      echo '--- wc test: FAIL (output mismatch)'
      failed_any=1
    fi
  else
    echo '--- wc test: FAIL (sequential wc failed)'
    failed_any=1
  fi
}

run_indexer () {
  echo '*** Starting indexer test.'
  rm -f mapreduce-*
  if ../bin/mapreducesequential indexer ../data/pg-*.txt; then
    sort mapreduce-out-0 > mapreduce-correct-indexer.txt
    rm -f mapreduce-out*
    run_cluster indexer 3
    sort mapreduce-out* 2>/dev/null | grep . > mapreduce-indexer-all || true
    if cmp -s mapreduce-indexer-all mapreduce-correct-indexer.txt; then
      echo '--- indexer test: PASS'
    else
      echo '--- indexer test: FAIL (output mismatch)'
      failed_any=1
    fi
  else
    echo '--- indexer test: FAIL (sequential indexer failed)'
    failed_any=1
  fi
}

run_mappar () {
  echo '*** Starting map parallelism test.'
  rm -f mapreduce-*
  run_cluster mtiming 2
  local NT
  NT=$(cat mapreduce-out* 2>/dev/null | grep '^times-' | wc -l | sed 's/ //g')
  if [ "$NT" != "2" ]; then
    echo '--- map parallelism test: FAIL (saw '"$NT"' workers)'
    failed_any=1
  elif cat mapreduce-out* 2>/dev/null | grep '^parallel.* 2' > /dev/null; then
    echo '--- map parallelism test: PASS'
  else
    echo '--- map parallelism test: FAIL (did not run in parallel)'
    failed_any=1
  fi
}

run_reducepar () {
  echo '*** Starting reduce parallelism test.'
  rm -f mapreduce-*
  run_cluster rtiming 2
  local NT
  NT=$(cat mapreduce-out* 2>/dev/null | grep '^[a-z] 2' | wc -l | sed 's/ //g')
  if [ "$NT" -lt "2" ]; then
    echo '--- reduce parallelism test: FAIL (too few parallel reduces)'
    failed_any=1
  else
    echo '--- reduce parallelism test: PASS'
  fi
}

run_jobcount () {
  echo '*** Starting job count test.'
  rm -f mapreduce-*
  run_cluster jobcount 4
  local NT
  NT=$(cat mapreduce-out* 2>/dev/null | awk '{print $2}')
  if [ "$NT" = "8" ]; then
    echo '--- job count test: PASS'
  else
    echo '--- job count test: FAIL (got "'"$NT"'" instead of 8)'
    failed_any=1
  fi
}

run_early () {
  echo '*** Starting early exit test.'
  rm -f mapreduce-*
  local DF=anydone$$
  rm -f "$DF"

  (POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreducecoordinator ../data/pg-*.txt ; touch "$DF") &
  sleep 0.5
  (POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker early_exit ; touch "$DF") &
  (POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker early_exit ; touch "$DF") &
  (POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker early_exit ; touch "$DF") &

  # wait for first to finish
  while [ ! -e "$DF" ]; do
    sleep 0.2
  done
  rm -f "$DF"

  sort mapreduce-out* 2>/dev/null | grep . > mapreduce-wc-all-initial || true

  # wait for all background jobs
  wait || true

  sort mapreduce-out* 2>/dev/null | grep . > mapreduce-wc-all-final || true

  if cmp -s mapreduce-wc-all-final mapreduce-wc-all-initial; then
    echo '--- early exit test: PASS'
  else
    echo '--- early exit test: FAIL (output changed after a worker exited)'
    failed_any=1
  fi
  rm -f mapreduce-*
}

run_crash () {
  echo '*** Starting crash test.'
  rm -f mapreduce-* mapreduce-done

  # first: get the "golden" output using the non-crashing plugin
  if ../bin/mapreducesequential nocrash ../data/pg-*.txt; then
    sort mapreduce-out-0 > mapreduce-correct-crash.txt
    rm -f mapreduce-out*
  else
    echo '--- crash test: FAIL (sequential nocrash failed)'
    failed_any=1
    return
  fi

  # start coordinator in background; when it exits we touch mapreduce-done
  ( POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreducecoordinator ../data/pg-*.txt ; touch mapreduce-done ) &
  coord_pid=$!

  # give it a moment to create the socket
  sleep 0.5

  # this is the socket your coordinator uses; change if yours differs
  SOCKNAME="socket/412-mapreduce.sock"

  # launch one crashing worker right away
  POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker crash >/dev/null 2>&1 &

  # hard cap so we don't loop forever that would suuuuuuck
  MAX_SECONDS=12
  start_ts=$(date +%s)

  # keep feeding crash workers **only while**:
  #  - coordinator is alive
  #  - we haven't timed out
  while kill -0 "$coord_pid" 2>/dev/null; do
    now_ts=$(date +%s)
    elapsed=$(( now_ts - start_ts ))
    if [ $elapsed -ge $MAX_SECONDS ]; then
      echo "--- crash test: coordinator did not finish within ${MAX_SECONDS}s, marking FAIL"
      # kill coordinator, break out
      kill "$coord_pid" 2>/dev/null || true
      failed_any=1
      break
    fi

    # if the socket is still there, it means the coordinator is still serving
    if [ -e "$SOCKNAME" ]; then
      POLL_MS="$POLL_MS" $TIMEOUT ../bin/mapreduceworker crash >/dev/null 2>&1 || true
    fi

    sleep 1
  done

  # wait for background jobs to finish (coordinate + any last workers)
  wait 2>/dev/null || true

  # if we never got the done file, bail
  if [ ! -f mapreduce-done ]; then
    echo '--- crash test: FAIL (coordinator never completed)'
    failed_any=1
    return
  fi

  # compare outputs
  sort mapreduce-out* 2>/dev/null | grep . > mapreduce-crash-all || true
  if cmp -s mapreduce-crash-all mapreduce-correct-crash.txt; then
    echo '--- crash test: PASS'
  else
    echo '--- crash test: FAIL (output mismatch)'
    failed_any=1
  fi
}

########################################
# dispatch
for t in "${TESTS[@]}"; do
  case "$t" in
    wc)         run_wc ;;
    indexer)    run_indexer ;;
    mappar)     run_mappar ;;
    reducepar)  run_reducepar ;;
    jobcount)   run_jobcount ;;
    early)      run_early ;;
    crash)      run_crash ;;
    all)        # allow ./test-mapreduce.sh all
                run_wc
                run_indexer
                run_mappar
                run_reducepar
                run_jobcount
                run_early
                run_crash
                break
                ;;
    *)
      echo "*** unknown test '$t' (valid: wc indexer mappar reducepar jobcount early crash all)"
      failed_any=1
      ;;
  esac
done

########################################
if [ "$failed_any" -eq 0 ]; then
  echo '*** PASSED REQUESTED TESTS'
  exit 0
else
  echo '*** SOME REQUESTED TESTS FAILED'
  exit 1
fi
