CC      := gcc
CFLAGS  := -std=c11 -Wall -Wextra -O2
SAN     := -fsanitize=address,undefined -fno-omit-frame-pointer
LDFLAGS := $(SAN)
CFLAGS  += $(SAN)

all: clean build-example build-tests run-tests

arena.o: arena.c arena.h
	$(CC) $(CFLAGS) -c arena.c

arena_tests.o: arena_tests.c arena.h
	$(CC) $(CFLAGS) -c arena_tests.c

build-example:
	$(CC) $(CFLAGS) -o arena_example main.c arena.c

build-tests: arena_tests.o arena.o
	$(CC) $(CFLAGS) -o arena_tests $^ $(LDFLAGS)

run-tests:
	./arena_tests

run-example:
	./arena_example

clean:
	rm -f *.o arena_tests
	rm -f *.o arena_example
