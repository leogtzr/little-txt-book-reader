import curses
from curses import textpad
import sys
from sys import stderr

if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

try:
    with open(filename, 'r') as f:
        lines = [line.rstrip('\n') for line in f.readlines()]
except FileNotFoundError:
    sys.exit(f"error: file not found: {filename}\n")
else:
    print(lines)
