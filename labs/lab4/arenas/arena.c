#include "arena.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdalign.h>   // super important to keep our bytes algined per OS - alignof(max_align_t)
#include <stdint.h>     // this too for alignof(max_align_t

// TODO: all functions - handle if a bad arena is passed :/

/*
* initiales our memory arena with a set size
* returns nothing or exits if our OS is out of mem
*/
void arena_init(Arena *arena, size_t buffer_size) {
    // TODO: how should we use the buffer_size in an arena?
    arena->buff = malloc(NULL); // using NULL here is silly?
    if (!arena->buff) {
        fprintf(stderr, "arena_init: out of memory - oh darn\n");
        exit(1);
    }
    // anything else to do here?
    return NULL; // look at the code comments, should we be explicitly returning anything?
}

/*
* uses our arena to grab a piece of heap memory
* returns a pointer containing the new heap memory
*/
void *arena_alloc(Arena *arena, size_t bytes) {
    // TODO: what is the number bytes for alignment based on our OS?
    // based on what you can find in tests.c?
    int bytes_align = 31;
    size_t aligned = (bytes + bytes_align) & ~((size_t)bytes_align);

    if (arena->offset + aligned > arena->size) {
        fprintf(stderr, "arena_alloc: overflow (%zu > %zu) - nooooes!\n",
                arena->offset + aligned, arena->size);
        exit(1);
    }
    // TODO: how should you use and set the offset here?
    return NULL; // what do we need to return here?
}

/*
* resets arena for reuse
* returns nothing
*/
void arena_reset(Arena *arena) {
    // TODO: reset the arena but do not free anything
}

/*
* clean up arean: free and reset everything related to the arena
* returns nothing
* 
* if not implemented correctly you will get at runtime...
* LeakSanitizer: detected memory leaks
* 
*/
void arena_free(Arena *arena) {
    // TODO: clean up everything, as if the arena and its parts never existed
}