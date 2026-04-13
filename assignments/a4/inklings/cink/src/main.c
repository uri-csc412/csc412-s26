#include "sim.h"
#include "ascii_art.h"
#include "config.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char **argv) {
    AppConfig cfg;
    SimState state;

    if (load_config(&cfg) != 0) {
        fprintf(stderr, "ERROR: failed to load config\n");
        return 1;
    }

    // TODO: override cfg.simulation rows/cols/threads from argv[] if valid.
    (void)argc;
    (void)argv;

    if (sim_init(&state, &cfg) != 0) {
        fprintf(stderr, "ERROR: failed to initialize simulation\n");
        return 1;
    }

    // TODO: create inkling threads and producer threads.
    // TODO: enter event loop.
    event_loop(&state);

    // TODO: join threads and clean shutdown.
    sim_shutdown(&state);

    return 0;
}
