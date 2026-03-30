# Makefile for M4 MapReduce solution CSC412

# where to put built binaries
BIN := bin

# go build flags (add -race or -v if you want)
GOFLAGS := -race

.PHONY: all build clean test test-many coordinator worker sequential run

all: build

build: $(BIN) coordinator worker sequential

$(BIN):
	mkdir -p $(BIN)

coordinator: $(BIN)
	go build $(GOFLAGS) -o $(BIN)/mapreducecoordinator ./coordinator

worker: $(BIN)
	go build $(GOFLAGS) -o $(BIN)/mapreduceworker ./worker

sequential: $(BIN)
	go build $(GOFLAGS) -o $(BIN)/mapreducesequential ./sequential

# run the original test script(s)
test: build
	./test-mapreduce.sh

test-many: build
	./test-mapreduce-many.sh

# usage: make run N=3
run: build
	@N=$(N); \
	if [ -z "$$N" ]; then \
		echo "Please provide number of workers, e.g. make run N=3"; \
		exit 1; \
	fi; \
	mkdir -p autograder-logs-apps; \
	cd autograder-logs-apps && { \
		echo "Starting coordinator..."; \
		../bin/mapreducecoordinator ../data/pg-*.txt & \
		COORD_PID=$$!; \
		echo "Starting $$N workers..."; \
		PIDS=""; \
		for i in `seq 1 $$N`; do \
			../bin/mapreduceworker wc & \
			PIDS="$$PIDS $$!"; \
		done; \
		# wait for coordinator to finish
		wait $$COORD_PID; \
		echo "Coordinator done, stopping workers..."; \
		kill $$PIDS 2>/dev/null || true; \
	}


clean:
	go clean ./...
	rm -rf $(BIN)/*

test-wc:
	./bin/mapreducesequential wc ./data/pg-*.txt
	cat mapreduce-out-0 | head
