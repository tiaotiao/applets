#include "input.h"
#include "define.h"


int input(char* filename, struct Empolyee* emp) {
    // TODO
    int i = 0;
    int ret;

    while (1) {

        ret = read_ID(&emp[i].ID);
        // TODO check if End of File, break;

        // TODO check return value
        ret = read_Name(emp[i].FirstName);
        // TODO read 
    }
    
    return 0;
}


int read_ID(int* id_ptr) {
    int ret = read_int(id_ptr);
    if (ret != 0) {
        printf("read invalid id format\n");
        return -1;
    }
    
    int id = *id_ptr;
    // check valid
    if (id < 100000 || 999999 < id) {
        printf("id out of range %d\n", id);
        return -2;
    }
    return 0;
}

int read_Name(char* name) {
    // TODO
}

int read_Salary(float* salary_ptr) {
    // TODO
}


