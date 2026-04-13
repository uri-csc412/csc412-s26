#include "config.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>
#include <dirent.h>
#include <unistd.h>   // getcwd
#include <limits.h>   // PATH_MAX

typedef enum {
    SEC_NONE = 0,
    SEC_LOGS,
    SEC_SIMULATION,
    SEC_INK,
    SEC_TIMING,
    SEC_EVENT_LOOP,
    SEC_PRODUCERS,
    SEC_IO,
    SEC_UI,
    SEC_UI_ICON,
    SEC_UI_MESSAGES
} Section;

// trim leading + trailing whitespace in-place
static char *trim(char *s) {
    if (!s) return s;
    // leading
    while (*s && isspace((unsigned char)*s)) {
        s++;
    }
    // trailing
    char *end = s + strlen(s);
    while (end > s && isspace((unsigned char)end[-1])) {
        end--;
    }
    *end = '\0';
    return s;
}

// convert "true"/"false" (case-insensitive) to 1/0, fallback: atoi
static int parse_bool_or_int(const char *v) {
    if (!v) return 0;
    if (strcasecmp(v, "true") == 0)  return 1;
    if (strcasecmp(v, "false") == 0) return 0;
    return atoi(v);
}

// remove surrounding quotes if present: "...." -> ....
static void strip_quotes(char *v) {
    size_t len = strlen(v);
    if (len >= 2 && v[0] == '"' && v[len - 1] == '"') {
        // shift left one, overwrite final quote with '\0'
        memmove(v, v + 1, len - 2);
        v[len - 2] = '\0';
    }
}

// Find a directory named `name` in {name, ../name, ../../name} from CWD.
// On success, writes the absolute path into out_buf and returns.
// On failure, prints an error and exits(1).
static void find_relative_dir(const char *name, char *out_buf, size_t out_sz) {
    char cwd[PATH_MAX];

    if (!getcwd(cwd, sizeof(cwd))) {
        perror("getcwd");
        exit(1);
    }

    const char *candidates[] = {
        "",        // ./<name>
        "..",      // ../<name>
        "../..",   // ../../<name>
    };
    size_t n = sizeof(candidates) / sizeof(candidates[0]);

    for (size_t i = 0; i < n; i++) {
        char candidate[PATH_MAX];
        if (candidates[i][0] == '\0') {
            // cwd/name
            snprintf(candidate, sizeof(candidate), "%s/%s", cwd, name);
        } else {
            // cwd/../name, cwd/../../name, etc.
            snprintf(candidate, sizeof(candidate), "%s/%s/%s", cwd, candidates[i], name);
        }

        DIR *d = opendir(candidate);
        if (d) {
            closedir(d);
            snprintf(out_buf, out_sz, "%s", candidate);
            return;
        }
    }

    fprintf(stderr,
            "ERROR: could not find '%s' directory "
            "(tried %s/%s, ../%s, ../../%s)\n",
            name, cwd, name, name, name);
    exit(1); // bc the program can't do nothing :/
}

// NOTE: This is a simple, targeted parser for inklings_config.yaml.
// It assumes a fairly stable structure and "key: value" lines under
// known section headers.
int load_config(AppConfig *cfg) {
    char config_dir[PATH_MAX];
    find_relative_dir("config", config_dir, sizeof config_dir);

    char config_path[PATH_MAX];
    snprintf(config_path, sizeof config_path, "%s/%s", config_dir, "inklings_config.yaml");
    const char *path = config_path;

    FILE *f = fopen(path, "r");
    if (!f) {
        fprintf(stderr, "ERROR: load_config: failed to open config file '%s'\n", path);
        return -1;
    }


    memset(cfg, 0, sizeof(*cfg));

    char line[512];
    Section section = SEC_NONE;

    while (fgets(line, sizeof(line), f)) {
        char *p = trim(line);

        // skip empty or comment lines
        if (*p == '\0' || *p == '#') {
            continue;
        }

        // top-level sections
        if (strncmp(p, "logs:", 5) == 0) {
            section = SEC_LOGS;
            continue;
        } else if (strncmp(p, "simulation:", 11) == 0) {
            section = SEC_SIMULATION;
            continue;
        } else if (strncmp(p, "ink:", 4) == 0) {
            section = SEC_INK;
            continue;
        } else if (strncmp(p, "timing:", 7) == 0) {
            section = SEC_TIMING;
            continue;
        } else if (strncmp(p, "event_loop:", 11) == 0) {
            section = SEC_EVENT_LOOP;
            continue;
        } else if (strncmp(p, "producers:", 10) == 0) {
            section = SEC_PRODUCERS;
            continue;
        } else if (strncmp(p, "io:", 3) == 0) {
            section = SEC_IO;
            continue;
        } else if (strncmp(p, "ui:", 3) == 0) {
            section = SEC_UI;
            continue;
        }

        // nested UI subsections aren't these fun :/
        if (section == SEC_UI && strncmp(p, "icon_directions:", 16) == 0) {
            section = SEC_UI_ICON;
            continue;
        } else if (section == SEC_UI && strncmp(p, "messages:", 9) == 0) {
            section = SEC_UI_MESSAGES;
            continue;
        }

        // key: value lines
        char *colon = strchr(p, ':');
        if (!colon) {
            continue; // not a key: value line we care about
        }

        *colon = '\0';
        char *key   = trim(p);
        char *value = trim(colon + 1);

        // strip quotes for string values
        strip_quotes(value);

        // dispatch based on section + key
        switch (section) {
            
        case SEC_LOGS:
            if (strcmp(key, "dir") == 0) {
                snprintf(cfg->logs.dir, sizeof(cfg->logs.dir), "%s", value);
            } else if (strcmp(key, "actions") == 0) {
                snprintf(cfg->logs.actions, sizeof(cfg->logs.actions), "%s", value);
            }
            break;

        case SEC_SIMULATION:
            if (strcmp(key, "default_rows") == 0)
                cfg->simulation.default_rows = atoi(value);
            else if (strcmp(key, "default_cols") == 0)
                cfg->simulation.default_cols = atoi(value);
            else if (strcmp(key, "min_rows") == 0)
                cfg->simulation.min_rows = atoi(value);
            else if (strcmp(key, "min_cols") == 0)
                cfg->simulation.min_cols = atoi(value);
            else if (strcmp(key, "min_threads") == 0)
                cfg->simulation.min_threads = atoi(value);
            else if (strcmp(key, "default_threads") == 0)
                cfg->simulation.default_threads = atoi(value);
            else if (strcmp(key, "draw_colored_traveler_heads") == 0)
                cfg->simulation.draw_colored_traveler_heads = parse_bool_or_int(value);
            break;

        case SEC_INK:
            if (strcmp(key, "max_level") == 0)
                cfg->ink.max_level = atoi(value);
            else if (strcmp(key, "max_add_ink") == 0)
                cfg->ink.max_add_ink = atoi(value);
            else if (strcmp(key, "refill_ink") == 0)
                cfg->ink.refill_ink = atoi(value);
            else if (strcmp(key, "initial_red_level") == 0)
                cfg->ink.initial_red_level = atoi(value);
            else if (strcmp(key, "initial_green_level") == 0)
                cfg->ink.initial_green_level = atoi(value);
            else if (strcmp(key, "initial_blue_level") == 0)
                cfg->ink.initial_blue_level = atoi(value);
            break;

        case SEC_TIMING:
            if (strcmp(key, "min_producer_sleep_us") == 0)
                cfg->timing.min_producer_sleep_us = atoi(value);
            else if (strcmp(key, "producer_sleep_us") == 0)
                cfg->timing.producer_sleep_us = atoi(value);
            else if (strcmp(key, "inkling_sleep_us") == 0)
                cfg->timing.inkling_sleep_us = atoi(value);
            else if (strcmp(key, "inkling_sleep_step_us") == 0)
                cfg->timing.inkling_sleep_step_us = atoi(value);
            else if (strcmp(key, "inkling_sleep_min_us") == 0)
                cfg->timing.inkling_sleep_min_us = atoi(value);
            else if (strcmp(key, "inkling_sleep_max_us") == 0)
                cfg->timing.inkling_sleep_max_us = atoi(value);
            else if (strcmp(key, "arg_warning_sleep_s") == 0)
                cfg->timing.arg_warning_sleep_s = atoi(value);
            else if (strcmp(key, "shutdown_delay_s") == 0)
                cfg->timing.shutdown_delay_s = atoi(value);
            else if (strcmp(key, "event_loop_interval_ms") == 0)
                cfg->timing.event_loop_interval_ms = atoi(value);
            else if (strcmp(key, "event_loop_sleep_s") == 0)
                cfg->timing.event_loop_sleep_s = atoi(value);
            break;

        case SEC_EVENT_LOOP:
            if (strcmp(key, "max_recursive_depth") == 0)
                cfg->event_loop.max_recursive_depth = atoi(value);
            break;

        case SEC_PRODUCERS:
            if (strcmp(key, "speedup_factor_num") == 0)
                cfg->producers.speedup_factor_num = atoi(value);
            else if (strcmp(key, "speedup_factor_den") == 0)
                cfg->producers.speedup_factor_den = atoi(value);
            else if (strcmp(key, "slowdown_factor_num") == 0)
                cfg->producers.slowdown_factor_num = atoi(value);
            else if (strcmp(key, "slowdown_factor_den") == 0)
                cfg->producers.slowdown_factor_den = atoi(value);
            break;

        case SEC_IO:
            if (strcmp(key, "pipe_path") == 0) {
                snprintf(cfg->io.pipe_path, sizeof(cfg->io.pipe_path), "%s", value);
            }
            break;

        case SEC_UI_ICON:
            if (strcmp(key, "north") == 0)
                snprintf(cfg->ui.icon_directions.north,
                         sizeof(cfg->ui.icon_directions.north), "%s", value);
            else if (strcmp(key, "west") == 0)
                snprintf(cfg->ui.icon_directions.west,
                         sizeof(cfg->ui.icon_directions.west), "%s", value);
            else if (strcmp(key, "south") == 0)
                snprintf(cfg->ui.icon_directions.south,
                         sizeof(cfg->ui.icon_directions.south), "%s", value);
            else if (strcmp(key, "east") == 0)
                snprintf(cfg->ui.icon_directions.east,
                         sizeof(cfg->ui.icon_directions.east), "%s", value);
            else if (strcmp(key, "none") == 0)
                snprintf(cfg->ui.icon_directions.none,
                         sizeof(cfg->ui.icon_directions.none), "%s", value);
            break;

        case SEC_UI_MESSAGES:
            if (strcmp(key, "ink_tank_header") == 0)
                snprintf(cfg->ui.messages.ink_tank_header,
                         sizeof(cfg->ui.messages.ink_tank_header), "%s", value);
            else if (strcmp(key, "live_threads_label") == 0)
                snprintf(cfg->ui.messages.live_threads_label,
                         sizeof(cfg->ui.messages.live_threads_label), "%s", value);
            else if (strcmp(key, "no_args_message") == 0)
                snprintf(cfg->ui.messages.no_args_message,
                         sizeof(cfg->ui.messages.no_args_message), "%s", value);
            else if (strcmp(key, "quit_no_threads_message") == 0)
                snprintf(cfg->ui.messages.quit_no_threads_message,
                         sizeof(cfg->ui.messages.quit_no_threads_message), "%s", value);
            else if (strcmp(key, "quit_esc_message") == 0)
                snprintf(cfg->ui.messages.quit_esc_message,
                         sizeof(cfg->ui.messages.quit_esc_message), "%s", value);
            else if (strcmp(key, "quit_generic_message") == 0)
                snprintf(cfg->ui.messages.quit_generic_message,
                         sizeof(cfg->ui.messages.quit_generic_message), "%s", value);
            break;

        default:
            // ignore lines outside known sections
            break;
        }
    }

    fclose(f);
    return 0;
}
