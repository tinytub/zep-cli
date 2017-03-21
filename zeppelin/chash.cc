#include <string>
#include <iostream>
#include <functional>
#include <stdint.h>

extern "C" uint64_t chash(const char*);
uint64_t chash(const char *str) {
  return std::hash<std::string>()(std::string(str));
}

//int main() {
 //   std::cout << chash("hello world") << std::endl;
//}
