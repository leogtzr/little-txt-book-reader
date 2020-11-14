# TODO: finish the help screen.

import curses
from curses import textpad
import sys
from sys import stderr
from enum import Enum


class WindowMode(Enum):
    reading = 1
    help = 2


class BookWindowNavigation:
    '''This class will contain everything related with the object navigation
        Current page number, navigation mode (help, reading), etc
    '''

    def __init__(self, book_number_lines, window_height, window_width):
        self._book_number_of_lines = book_number_lines
        self.window_height = window_height
        self.window_width = window_width
        self.from_line = 0
        self.to_line = window_height
        self.current_row = 0
        self.line_number = 1
        self.window_mode = WindowMode.reading

    def book_number_lines(self):
        return self._book_number_of_lines


if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

KEY_ESCAPE_CODE = 27
HIGHLIGHT_COLOR_PAIRCODE = 1
STATUSBAR_COLOR_PAIRCODE = 2
HELP_KEY_CODES = [ord('h'), ord('H')]


def book_chunk(lines, from_line, to_line, book_number_of_lines):
    return lines[from_line:to_line]


def print_page_section(stdscr, selected_row_idx, book_page):
    for idx, book_page_line in enumerate(book_page):
        if idx == selected_row_idx:
            stdscr.attron(curses.color_pair(HIGHLIGHT_COLOR_PAIRCODE))
            stdscr.addstr(idx, 0, book_page_line)
            stdscr.attroff(curses.color_pair(HIGHLIGHT_COLOR_PAIRCODE))
        else:
            stdscr.addstr(idx, 0, book_page_line)


def print_status_bar(stdscr, pos_height, pos_width, status_text):
    stdscr.attron(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))
    stdscr.addstr(pos_height, pos_width//2, status_text)
    stdscr.attroff(curses.color_pair(1))


def print_help_screen(stdscr):
    screen_height, screen_width = stdscr.getmaxyx()
    border_offset = 3
    box = [[border_offset, border_offset], [
        screen_height-border_offset, screen_width-border_offset]]
    textpad.rectangle(
        stdscr, box[0][0], box[0][1], box[1][0], box[1][1])

    help_entries = [
        'Down    -> Go Down',
        'Up      -> Go Up',
        'G       -> Go To',
        '.       -> Toggle Status Bar',
        'ESC     -> Closes the program/Dialogs',
        'S       -> Save Progress',
        'H       -> Show the Help Dialog'
    ]

    for idx, help_entry in enumerate(help_entries):
        stdscr.addstr(border_offset + idx + 1, border_offset+1, help_entry)


def print_page(stdscr, lines, bookwnd_nav):
    book_page = book_chunk(lines, bookwnd_nav.from_line,
                           bookwnd_nav.to_line, bookwnd_nav.book_number_lines())
    print_page_section(stdscr, bookwnd_nav.current_row, book_page)
    print_status_bar(
        stdscr, bookwnd_nav.window_height - 1, bookwnd_nav.window_width, f"Current line: {bookwnd_nav.line_number}")


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

        bookwnd_nav = BookWindowNavigation(
            book_number_of_lines, MAX_HEIGHT, MAX_WIDTH)

        print_page(stdscr, lines, bookwnd_nav)

        while True:
            key = stdscr.getch()

            if key in [KEY_ESCAPE_CODE]:
                if bookwnd_nav.window_mode == WindowMode.help:
                    bookwnd_nav.window_mode = WindowMode.reading
                elif bookwnd_nav.window_mode == WindowMode.reading:
                    stdscr.refresh()
                    sys.exit(0)

            elif key in HELP_KEY_CODES:
                stdscr.clear()
                bookwnd_nav.window_mode = WindowMode.help
                print_help_screen(stdscr)

            elif key == curses.KEY_UP:
                if bookwnd_nav.line_number > 1:
                    bookwnd_nav.current_row -= 1
                    bookwnd_nav.line_number -= 1
                if bookwnd_nav.from_line > 0:
                    bookwnd_nav.from_line -= 1
                    bookwnd_nav.to_line -= 1
                    bookwnd_nav.line_number -= 1

            elif key == curses.KEY_DOWN:
                if bookwnd_nav.line_number < book_number_of_lines:
                    if bookwnd_nav.current_row >= (bookwnd_nav.window_height - 2):
                        bookwnd_nav.from_line += 1
                        bookwnd_nav.to_line += 1
                        stdscr.clear()
                    else:
                        bookwnd_nav.current_row += 1
                    bookwnd_nav.line_number += 1

            if bookwnd_nav.window_mode == WindowMode.reading:
                stdscr.clear()
                print_page(stdscr, lines, bookwnd_nav)


curses.wrapper(main)
