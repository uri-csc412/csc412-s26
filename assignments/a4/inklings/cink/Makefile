CC      = gcc
CFLAGS  = -D _DEFAULT_SOURCE  -Wall -Wextra -pedantic -std=c11 -g -O3 -Iinclude
LDFLAGS = -pthread

BINDIR          = bin
BINARY_INKLINGS = inklings_c
BINARY_LOGS     = logs_c

SRC_INK  = src/main.c src/config.c src/sim.c src/ascii_art.c
SRC_LOGS = src/logs.c src/config.c

OBJ_INK  = $(SRC_INK:.c=.o)
OBJ_LOGS = $(SRC_LOGS:.c=.o)

RUN_ARGS ?= 8 8 4

all: $(BINARY_INKLINGS) $(BINARY_LOGS)

$(BINARY_INKLINGS): $(OBJ_INK)
	$(CC) $(CFLAGS) -o $(BINDIR)/$@ $^ $(LDFLAGS)

$(BINARY_LOGS): $(OBJ_LOGS)
	$(CC) $(CFLAGS) -o $(BINDIR)/$@ $^ $(LDFLAGS)

run: $(BINARY_INKLINGS)
	./$(BINDIR)/$(BINARY_INKLINGS) $(RUN_ARGS)

run-logs: $(BINARY_LOGS)
	./$(BINDIR)/$(BINARY_LOGS)

clean:
	rm -f $(OBJ_INK) $(OBJ_LOGS) \
	      $(BINDIR)/$(BINARY_INKLINGS) $(BINDIR)/$(BINARY_LOGS)
	rm -rf ../logs/*

.PHONY: all run clean run-logs
