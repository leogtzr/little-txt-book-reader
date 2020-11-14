import curses
from curses import textpad
import sys
from sys import stderr
from enum import Enum


class WindowMode(Enum):
    reading = 1
    help = 2


if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

# The following constant is "temporary", ideally it should be calculated based:
KEY_ESCAPE_CODE = 27
HIGHLIGHT_COLOR_PAIRCODE = 1
STATUSBAR_COLOR_PAIRCODE = 2


def book_chunk(lines, from_line, to_line, book_number_of_lines):
    return lines[from_line:to_line]


def print_page(stdscr, selected_row_idx, book_page):
    for idx, book_page_line in enumerate(book_page):
        if idx == selected_row_idx:
            stdscr.attron(curses.color_pair(HIGHLIGHT_COLOR_PAIRCODE))
            stdscr.addstr(idx, 0, book_page_line)
            stdscr.attroff(curses.color_pair(HIGHLIGHT_COLOR_PAIRCODE))
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
        curses.init_pair(HIGHLIGHT_COLOR_PAIRCODE,
                         curses.COLOR_BLACK, curses.COLOR_WHITE)
        curses.init_pair(STATUSBAR_COLOR_PAIRCODE,
                         curses.COLOR_BLACK, curses.COLOR_GREEN)
        book_number_of_lines = len(lines)
        MAX_HEIGHT, MAX_WIDTH = stdscr.getmaxyx()
        from_line = 0
        to_line = MAX_HEIGHT
        current_row = 0
        line_number = 1
        window_mode = WindowMode.reading

        book_page = book_chunk(lines, from_line, to_line, book_number_of_lines)
        print_page(stdscr, current_row, book_page)
        print_status_bar(
            stdscr, MAX_HEIGHT - 1, f"Current line: {line_number}")

        # Loop
        while True:
            key = stdscr.getch()
            stdscr.clear()

            if key in [KEY_ESCAPE_CODE]:
                if window_mode == WindowMode.help:
                    window_mode = WindowMode.reading
                elif window_mode == WindowMode.reading:
                    stdscr.refresh()
                    sys.exit(0)

            elif key in [72, 104]:
                window_mode = WindowMode.help
                stdscr.clear()

                sh, sw = stdscr.getmaxyx()
                box = [[3, 3], [sh-3, sw-3]]
                textpad.rectangle(
                    stdscr, box[0][0], box[0][1], box[1][0], box[1][1])

            elif key == curses.KEY_UP:
                if line_number > 1:
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
                        stdscr.clear()
                    else:
                        current_row += 1
                    line_number += 1

            stdscr.refresh()

            if window_mode == WindowMode.reading:
                book_page = book_chunk(
                    lines, from_line, to_line-1, book_number_of_lines)
                print_page(stdscr, current_row, book_page)
                print_status_bar(
                    stdscr, MAX_HEIGHT - 1, f"Current line: {line_number}")


curses.wrapper(main)
