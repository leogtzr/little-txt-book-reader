import curses
from curses import textpad
import sys
from sys import stderr
from enum import Enum
import utils
import re
import book


if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

KEY_ESCAPE_CODE = 27
HIGHLIGHT_COLOR_PAIRCODE = 1
STATUSBAR_COLOR_PAIRCODE = 2
HELP_KEY_CODES = [ord('h'), ord('H')]
TOGGLE_STATUSBAR_KEY_CODE = ord('.')
SHOW_PERCENTAGE_POINTS_KEY_CODES = [ord('P'), ord('p')]
GOTO_KEY_CODES = [ord('g'), ord('G')]


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


def print_status_bar(stdscr, bookwnd_nav):
    if not bookwnd_nav.show_status_bar:
        return

    perc = utils.percent(bookwnd_nav.line_number,
                         bookwnd_nav.book_number_lines())

    if bookwnd_nav.show_percentage_points:
        lines_to_new_p_point = utils.lines_to_change_percentage_point(
            bookwnd_nav.line_number, bookwnd_nav.book_number_lines())
        status_text = f"{bookwnd_nav.line_number} of {bookwnd_nav.book_number_lines()}      (%{perc:.1f})  (> {lines_to_new_p_point})"
    else:
        status_text = f"{bookwnd_nav.line_number} of {bookwnd_nav.book_number_lines()}      (%{perc:.1f}) [cr: {bookwnd_nav.current_row}] - [wh: {bookwnd_nav.window_height}]: from: {bookwnd_nav.from_line}, to: {bookwnd_nav.to_line}"

    pos_height = bookwnd_nav.window_height - 1
    pos_width = bookwnd_nav.window_width
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
        'H       -> Show the Help Dialog',
        'P       -> Show Percentage Points'
    ]

    for idx, help_entry in enumerate(help_entries):
        stdscr.addstr(border_offset + idx + 1, border_offset+1, help_entry)


def show_goto_dialog(stdscr, bookwnd_nav):
    screen_height, screen_width = stdscr.getmaxyx()
    border_offset = 2
    box = [[border_offset, border_offset], [
        screen_height-border_offset, screen_width-border_offset]]
    textpad.rectangle(
        stdscr, box[0][0], box[0][1], box[1][0], box[1][1])
    stdscr.addstr(border_offset + 1, border_offset + 1, "Go To: ")
    curses.echo()
    input = stdscr.getstr(
        border_offset + 1, (border_offset + 1) + len('Go To: '), 20)
    input = input.strip()
    input = input.rstrip()

    curses.noecho()
    return re.sub('\D', '', input.decode("utf-8"))


def print_page(stdscr, lines, bookwnd_nav):
    book_page = book_chunk(lines, bookwnd_nav.from_line,
                           bookwnd_nav.to_line, bookwnd_nav.book_number_lines())
    print_page_section(stdscr, bookwnd_nav.current_row, book_page)
    print_status_bar(stdscr, bookwnd_nav)


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

        bookwnd_nav = book.BookWindowNavigation(
            book_number_of_lines, MAX_HEIGHT, MAX_WIDTH)

        print_page(stdscr, lines, bookwnd_nav)

        while True:
            key = stdscr.getch()

            if key in [KEY_ESCAPE_CODE]:
                if bookwnd_nav.window_mode == book.WindowMode.help:
                    bookwnd_nav.window_mode = book.WindowMode.reading
                elif bookwnd_nav.window_mode == book.WindowMode.reading:
                    stdscr.refresh()
                    sys.exit(0)
                elif bookwnd_nav.window_mode == book.WindowMode.goto:
                    bookwnd_nav.window_mode = book.WindowMode.reading

            elif key in HELP_KEY_CODES:
                stdscr.clear()
                bookwnd_nav.window_mode = book.WindowMode.help
                print_help_screen(stdscr)

            elif key in GOTO_KEY_CODES:
                stdscr.clear()
                bookwnd_nav.window_mode = book.WindowMode.goto
                input_goto = show_goto_dialog(stdscr, bookwnd_nav)
                if input_goto:
                    bookwnd_nav.from_line = int(input_goto)
                    bookwnd_nav.to_line = bookwnd_nav.from_line + bookwnd_nav.window_height
                    bookwnd_nav.line_number = bookwnd_nav.from_line
                    bookwnd_nav.current_row = 0
                bookwnd_nav.window_mode = book.WindowMode.reading

            elif key == curses.KEY_DOWN:
                if bookwnd_nav.current_row == (bookwnd_nav.window_height - 1):
                    # Reset:
                    bookwnd_nav.current_row = 0
                    bookwnd_nav.line_number += 1
                    bookwnd_nav.from_line += bookwnd_nav.window_height
                    bookwnd_nav.to_line = bookwnd_nav.from_line + bookwnd_nav.window_height
                else:
                    bookwnd_nav.current_row += 1
                    bookwnd_nav.line_number += 1

            elif key == curses.KEY_UP:
                if bookwnd_nav.current_row == 0:
                    # Do we have enough space to sub up?
                    if bookwnd_nav.line_number > bookwnd_nav.window_height:
                        bookwnd_nav.current_row = bookwnd_nav.window_height - 1
                        bookwnd_nav.line_number -= 1
                        bookwnd_nav.from_line -= bookwnd_nav.window_height
                        bookwnd_nav.to_line -= bookwnd_nav.window_height
                else:
                    bookwnd_nav.current_row -= 1
                    bookwnd_nav.line_number -= 1

            elif key == TOGGLE_STATUSBAR_KEY_CODE:
                bookwnd_nav.show_status_bar = not bookwnd_nav.show_status_bar

            elif key in SHOW_PERCENTAGE_POINTS_KEY_CODES:
                bookwnd_nav.show_percentage_points = not bookwnd_nav.show_percentage_points

            if bookwnd_nav.window_mode == book.WindowMode.reading:
                stdscr.clear()
                print_page(stdscr, lines, bookwnd_nav)


curses.wrapper(main)
