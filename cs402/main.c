#include <stdio.h>

#include "define.h"
#include "file_lib.h"
#include "input.h"
#include "menu.h"

#define N 10000

struct Employee Emp[N];


void main(int numOfArgs, char** args) {
    int ret;

    // check arguments
    if (numOfArgs != 2) {
        printf("invalid argument: [filename]\n");
        return;
    }

    char* filename = argv[1];

    ret = input(filename);
}
