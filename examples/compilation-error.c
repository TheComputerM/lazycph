int main(void) {
    #ifdef __GNUC__
    __asm__("this_will_cause_a_compilation_errors");
    #endif
    return 0;
}
