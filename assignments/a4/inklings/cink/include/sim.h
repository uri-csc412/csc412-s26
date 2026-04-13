#ifndef INKLINGS_H
#define INKLINGS_H

#include "config.h"
#include <pthread.h>

// Travel direction
typedef enum {
    DIR_NORTH = 0,
    DIR_WEST,
    DIR_SOUTH,
    DIR_EAST,
    DIR_NUM_TRAVEL_DIRECTIONS
} TravelDirection;

typedef enum {
    TRAV_RED = 0,
    TRAV_GREEN,
    TRAV_BLUE,
    TRAV_NUM_TYPES
} InklingType;

// Inkling (traveler) state
typedef struct {
    InklingType type;
    int row;
    int col;
    TravelDirection dir;
    int is_live;        // treat as bool
} InklingInfo;

// global simulation state
typedef struct {
    AppConfig config;
    int **grid;              // 2D grid of cells
    int num_rows;
    int num_cols;

    InklingInfo *inklings;
    int num_inklings;

    int red_level;
    int green_level;
    int blue_level;

    pthread_mutex_t red_lock;
    pthread_mutex_t green_lock;
    pthread_mutex_t blue_lock;

    pthread_mutex_t grid_lock;
} SimState;

int  sim_init(SimState *state, const AppConfig *cfg);
void sim_shutdown(SimState *state);
void *inkling_thread(void *arg);
void *producer_thread_red(void *arg);
void *producer_thread_green(void *arg);
void *producer_thread_blue(void *arg);

#endif // INKLINGS_H
