#include <iostream>

// Declared but never defined - linker error, not caught by LSP
void missing_function();

int main() {
    std::cout << "Hello World\n";
    missing_function();
    return 0;
}
