// Source - https://stackoverflow.com/q
// Posted by AraMod, modified by community. See post 'Timeline' for change
// history Retrieved 2025-11-29, License - CC BY-SA 3.0

#include <stdio.h>

int main() {
  while (1) {
    int ch = getc(stdin);
    fflush(stdout);
    if (ch == EOF)
      break;

    putc(ch, stdout);
  }
  return 0;
}
