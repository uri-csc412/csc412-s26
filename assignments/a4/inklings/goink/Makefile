GOCMD       := go
GOBUILD     := $(GOCMD) build
GOCLEAN     := $(GOCMD) clean

BINDIR          := bin
BINARY_INKLINGS := inklings_go
BINARY_LOGS     := logs_go

RUN_ARGS ?= 8 8 4

all: $(BINARY_INKLINGS) $(BINARY_LOGS)

$(BINARY_INKLINGS):
	$(GOBUILD) -o $(BINDIR)/$(BINARY_INKLINGS) ./cmd/inklings

$(BINARY_LOGS):
	$(GOBUILD) -o $(BINDIR)/$(BINARY_LOGS) ./cmd/logs

run: $(BINARY_INKLINGS)
	./$(BINDIR)/$(BINARY_INKLINGS) $(RUN_ARGS)

run-logs: $(BINARY_LOGS)
	./$(BINDIR)/$(BINARY_LOGS)

clean:
	$(GOCLEAN)
	rm -f $(BINDIR)/$(BINARY_INKLINGS) $(BINDIR)/$(BINARY_LOGS)
	rm -rf ../logs/*

.PHONY: all clean run run-logs
