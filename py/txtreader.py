import curses
import sys
from sys import stderr

if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

# The following constant is "temporary", ideally it should be calculated based:
# ADVANCE = 30
KEY_ESCAPE_CODE = 27
STATUSBAR_COLOR_PAIRCODE = 2

# book_page could probably a slice.


def book_chunk(lines, from_line, to_line, book_number_of_lines):
    return lines[from_line:to_line]


def print_page(stdscr, selected_row_idx, book_page):
    for idx, book_page_line in enumerate(book_page):
        if idx == selected_row_idx:
            stdscr.attron(curses.color_pair(1))
            stdscr.addstr(idx, 0, book_page_line)
            stdscr.attroff(curses.color_pair(1))
        else:
            stdscr.addstr(idx, 0, book_page_line)


def print_status_bar(stdscr, position, status_text):
    stdscr.attron(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))
    stdscr.addstr(position, 0, status_text)
    stdscr.attroff(curses.color_pair(1))


def main(stdscr):
    try:
        with open(filename, 'r') as f:
            lines = [line.rstrip('\n') for line in f.readlines()]
    except FileNotFoundError:
        sys.exit(f"error: file not found: {filename}\n")
    else:
        curses.curs_set(0)
        curses.init_pair(1, curses.COLOR_BLACK, curses.COLOR_WHITE)
        curses.init_pair(STATUSBAR_COLOR_PAIRCODE,
                         curses.COLOR_BLACK, curses.COLOR_GREEN)

        book_number_of_lines = len(lines)

        MAX_HEIGHT, _ = stdscr.getmaxyx()
        ADVANCE = MAX_HEIGHT
        from_line = 0
        to_line = ADVANCE
        current_row = 0
        line_number = 1

        book_page = book_chunk(lines, from_line, to_line, book_number_of_lines)
        print_page(stdscr, current_row, book_page)
        print_status_bar(
            stdscr, MAX_HEIGHT - 1, f"Current line: {line_number}, from_line: {from_line}, to_line: {to_line}")

        # Loop
        while True:
            key = stdscr.getch()
            stdscr.clear()

            if key in [KEY_ESCAPE_CODE]:                # exit
                stdscr.refresh()
                sys.exit(0)

            elif key == curses.KEY_UP:
                if line_number > 1:
                    line_number -= 1
                    current_row -= 1
                if from_line > 0:
                    from_line -= 1
                    to_line -= 1
                    line_number -= 1

            elif key == curses.KEY_DOWN:
                if line_number < book_number_of_lines:
                    if current_row >= (MAX_HEIGHT - 2):
                        from_line += 1
                        to_line += 1
                        line_number += 1
                        stdscr.clear()
                    else:
                        current_row += 1
                        line_number += 1

            stdscr.refresh()

            # Update status bar:
            book_page = book_chunk(
                lines, from_line, to_line-1, book_number_of_lines)
            print_page(stdscr, current_row, book_page)
            print_status_bar(
                stdscr, MAX_HEIGHT - 1, f"Current line: {line_number}, from_line: {from_line}, to_line: {to_line}")


curses.wrapper(main)
