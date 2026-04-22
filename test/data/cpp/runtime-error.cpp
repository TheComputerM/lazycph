#include <iostream>
#include <vector>
#include <stdexcept>

int main() {
    std::vector<int> v;
    // Throws std::out_of_range at runtime
    std::cout << v.at(10) << "\n";
    return 0;
}
