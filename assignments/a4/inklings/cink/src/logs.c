#include "config.h"
#include <stdio.h>
#include <stdlib.h>
#include <dirent.h>
#include <string.h>
#include <unistd.h>   // getcwd
#include <limits.h>   // PATH_MAX

/* logs program

    Read all inkling*.txt files from the logs dir
    Merge lines into a actions.txt sorted by timestamp
    Be idempotent (running logs multiple times does not duplicate entries)
    Include all log entries (no missing lines)
    
*/

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

int main(void) {
    char logs_path[PATH_MAX];
    find_relative_dir("logs", logs_path, sizeof logs_path);

    printf("logs_c: using logs directory: %s\n", logs_path);
    printf("logs_c: please implement merging %s/*.txt into actions.txt\n",
           logs_path);

    return 0;
}
