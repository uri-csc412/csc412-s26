#include "sim.h"
#include <stdlib.h>
#include <stdio.h>

int sim_init(SimState *state, const AppConfig *cfg) {
    // TODO: allocate grid, initialize inklings at valid positions,
    // initialize ink levels, mutexes, etc.

    state->config = *cfg;
    state->num_rows = cfg->simulation.default_rows;
    state->num_cols = cfg->simulation.default_cols;
    state->red_level = cfg->ink.initial_red_level;
    state->green_level = cfg->ink.initial_green_level;
    state->blue_level = cfg->ink.initial_blue_level;

    // TODO: allocate and zero grid
    // TODO: initialize inklings array and place them on the grid
    // TODO: initialize mutexes

    return 0;
}

void sim_shutdown(SimState *state) {
    // TODO: free grid, inklings, destroy mutexes, etc.
    (void)state;
}

void *inkling_thread(void *arg) {
    // TODO: move one inkling around, consume ink, update grid, etc.
    (void)arg;
    return NULL;
}

void *producer_thread_red(void *arg) {
    // TODO: refill red ink using locks and sleep based on config timing.
    (void)arg;
    return NULL;
}

void *producer_thread_green(void *arg) {
    // TODO: refill green ink.
    (void)arg;
    return NULL;
}

void *producer_thread_blue(void *arg) {
    // TODO: refill blue ink.
    (void)arg;
    return NULL;
}
