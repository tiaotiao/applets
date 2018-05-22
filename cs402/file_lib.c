#include <stdio.h>

FILE* fileHandle = NULL;

int open_file(char* filename) {
    fileHandle = fopen(filename, "r");
    if (fileHandle == NULL) {
        return -1;
    }
    return 0;
}

int close_file() {
    // TODO check return value, return -1
    fclose(fileHandle);

    return 0;
}

int read_int(int* i) {
    int r = fscanf(fileHandle, "%d", i);
    
    // TODO End of file or bad format
    if (r ...) {
        return -1;
    }

    return 0;
}

int read_float(float* d);
int read_string(char* s);

