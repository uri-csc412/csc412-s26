#ifndef ASCII_ART_H
#define ASCII_ART_H

#include "sim.h"

// Terminal UI / ASCII drawing.

// ANSII Colors
typedef enum {
    TEXT_DEFAULT = 39,
    TEXT_BLACK   = 30,
    TEXT_RED     = 31,
    TEXT_GREEN   = 32,
    TEXT_YELLOW  = 33,
    TEXT_BLUE    = 34,
    TEXT_MAGENTA = 35,
    TEXT_CYAN    = 36,
    TEXT_WHITE   = 37
} TextColor;

void clear_terminal(void);
void draw_grid_and_inklings(const SimState *state);
void draw_state(const SimState *state);
void event_loop(SimState *state);

#endif // ASCII_ART_H
