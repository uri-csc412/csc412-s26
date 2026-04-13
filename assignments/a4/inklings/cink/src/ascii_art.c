#include "ascii_art.h"
#include <stdio.h>

void clear_terminal(void) {
    // Simple ANSI clear.
    printf("\033[H\033[J");
}

void draw_grid_and_inklings(const SimState *state) {
    // TODO: draw grid and inklings using ASCII.
    (void)state;
}

void draw_state(const SimState *state) {
    // TODO: display ink tank levels, number of live threads, etc.
    (void)state;
}

void event_loop(SimState *state) {
    // TODO: periodically redraw terminal, process keyboard input, etc.
    (void)state;
}


static void set_text_color(TextColor color) {
    printf("\033[%dm", (int)color);
}

static void reset_text_color(void) {
    set_text_color(TEXT_DEFAULT);
}