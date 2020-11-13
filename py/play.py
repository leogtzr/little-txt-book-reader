import sys
from sys import stderr

if len(sys.argv) != 2:
    sys.exit(f"error: missing file")

filename = sys.argv[1]

# The following constant is "temporary", ideally it should be calculated based:
ADVANCE = 30


def book_chunk(lines, from_line, to_line, book_number_of_lines):
    return lines[from_line:to_line]

# book_page could probably a slice.


try:
    with open(filename, 'r') as f:
        lines = [line.rstrip('\n') for line in f.readlines()]
except FileNotFoundError:
    sys.exit(f"error: file not found: {filename}\n")
else:

    book_number_of_lines = len(lines)

    from_line = 0
    to_line = ADVANCE
    print(f"from:{from_line}, to={to_line}")
    book_page = book_chunk(lines, from_line, to_line, book_number_of_lines)
#     book_page = lines[from_line:30_00 0]
    for idx, book_line in enumerate(book_page):
        print(book_line)

# # 1
# # 2
# # 3
# # 4
# # 5
# # 6
# # 7
# # 8
# # 9
# # 10
# # 11
# # 12
# # 13
# # 14
