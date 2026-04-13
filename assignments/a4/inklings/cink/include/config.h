#ifndef INKLINGS_CONFIG_H
#define INKLINGS_CONFIG_H

// Basic configuration structs that mirror config/inklings_config.yaml.

typedef struct {
    char dir[256];     // directory for log files
    char actions[256]; // filename for actions log
} LogsConfig;

typedef struct {
    int default_rows;
    int default_cols;
    int min_rows;
    int min_cols;
    int min_threads;
    int default_threads;
    int draw_colored_traveler_heads;  // treat as bool
} SimulationConfig;

typedef struct {
    int max_level;
    int max_add_ink;
    int refill_ink;
    int initial_red_level;
    int initial_green_level;
    int initial_blue_level;
} InkConfig;

typedef struct {
    int min_producer_sleep_us;
    int producer_sleep_us;
    int inkling_sleep_us;
    int inkling_sleep_step_us;
    int inkling_sleep_min_us;
    int inkling_sleep_max_us;
    int arg_warning_sleep_s;
    int shutdown_delay_s;
    int event_loop_interval_ms;
    int event_loop_sleep_s;
} TimingConfig;

typedef struct {
    int max_recursive_depth;
} EventLoopConfig;

typedef struct {
    int speedup_factor_num;
    int speedup_factor_den;
    int slowdown_factor_num;
    int slowdown_factor_den;
} ProducersConfig;

typedef struct {
    char pipe_path[256];
} IOConfig;

typedef struct {
    char north[8];
    char west[8];
    char south[8];
    char east[8];
    char none[8];
} IconDirectionsConfig;

typedef struct {
    char ink_tank_header[128];
    char live_threads_label[64];
    char no_args_message[256];
    char quit_no_threads_message[256];
    char quit_esc_message[256];
    char quit_generic_message[256];
} MessagesConfig;

typedef struct {
    IconDirectionsConfig icon_directions;
    MessagesConfig       messages;
} UIConfig;

typedef struct {
    LogsConfig       logs;
    SimulationConfig simulation;
    InkConfig        ink;
    TimingConfig     timing;
    EventLoopConfig  event_loop;
    ProducersConfig  producers;
    IOConfig         io;
    UIConfig         ui;
} AppConfig;

// Load configuration from config/inklings_config.yaml into the provided struct.
int load_config(AppConfig *cfg);

#endif // INKLINGS_CONFIG_H
