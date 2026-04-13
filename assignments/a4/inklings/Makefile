# select implementation: c or go (this is what the autograder will set)
IMPL ?= go

C_DIR := cink
GO_DIR := goink

ifeq ($(IMPL),c)
  SUBDIR    := $(C_DIR)
  BIN_NAME  := inklings_c
  LOGS_BIN  := logs_c
else ifeq ($(IMPL),go)
  SUBDIR    := $(GO_DIR)
  BIN_NAME  := inklings_go
  LOGS_BIN  := logs_go
else
  $(error Unknown IMPL '$(IMPL)'; must be 'c' or 'go')
endif

RUN_ARGS ?= 8 8 4

.PHONY: all inklings run logs run-logs clean remake

all: inklings

remake: clean inklings run

inklings:
	$(MAKE) -C $(SUBDIR)

run: inklings
	./$(SUBDIR)/bin/$(BIN_NAME) $(RUN_ARGS)

run-logs:
	./$(SUBDIR)/bin/$(LOGS_BIN)

clean:
	-$(MAKE) -C $(C_DIR) clean
	-$(MAKE) -C $(GO_DIR) clean
