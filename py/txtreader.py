# TODO: word building mode.
import curses
from curses import textpad
import sys
from enum import Enum
import utils
import re
import book
import progress
import os
from pathlib import Path
import subprocess
import shutil
import clipboard

if len(sys.argv) != 2:
    sys.exit(1)

filename = sys.argv[1]

# Note: the following might change:
PROGRAM_PATH_DIR = os.path.join(os.environ.get('HOME'), 'txt')
PROGRAM_PROGRESS_PATH_DIR = os.path.join(PROGRAM_PATH_DIR, 'progress')
PROGRAM_NOTES_PATH_DIR = os.path.join(PROGRAM_PATH_DIR, 'notes')
PROGRAM_WORDS_PATH_DIR = os.path.join(PROGRAM_PATH_DIR, 'words')

KEY_ESCAPE_CODE = 27
HIGHLIGHT_COLOR_PAIRCODE = 1
STATUSBAR_COLOR_PAIRCODE = 2
HELP_KEY_CODES = [ord('h'), ord('H')]
TOGGLE_STATUSBAR_KEY_CODE = ord('.')
SHOW_PERCENTAGE_POINTS_KEY_CODES = [ord('P'), ord('p')]
GOTO_KEY_CODES = [ord('g'), ord('G')]
SAVE_PROGRESS_KEY_CODE = [ord('s'), ord('S')]
ADD_NOTES_KEY_CODE = [ord('n'), ord('N')]
OPEN_RAE_WEBSITE_KEY_CODES = [ord('o'), ord('O')]
OPEN_GOODREADS_WEBSITE_KEY_CODES = [ord('r'), ord('R')]
WORD_BUILDING_KEY_CODES = [ord('w'), ord('W')]


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


def print_save_progress_status(stdscr, bookwnd_nav, filename):
    if not bookwnd_nav.show_status_bar:
        return

    base_filename = os.path.basename(filename)

    status_text = f"Status saved for: '{base_filename}''"

    pos_height = bookwnd_nav.window_height - 1
    pos_width = bookwnd_nav.window_width
    stdscr.attron(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))
    stdscr.addstr(pos_height, pos_width//2, status_text)
    stdscr.attroff(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))


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
        status_text = f"{bookwnd_nav.line_number} of {bookwnd_nav.book_number_lines()}      (%{perc:.1f})"

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
        'P       -> Show Percentage Points',
        'N       -> Open Notes file',
        'T       -> Toggle Status Bar Versions',
        'O       -> Opens RAE Web site with search from the clipboard.',
        'R       -> Opens GoodReads Web site with search from the clipboard.'
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


def get_progress_filepath(filename):
    base_filename = os.path.basename(filename)
    return os.path.join(PROGRAM_PROGRESS_PATH_DIR, base_filename)


def save_progress(filename, bookwnd_nav):
    abs_path = os.path.abspath(filename)
    base_filename = os.path.basename(filename)

    with open(os.path.join(PROGRAM_PROGRESS_PATH_DIR, base_filename), 'w') as progress_file:
        progress_file.write(
            f"{abs_path}|{bookwnd_nav.from_line}|{bookwnd_nav.to_line}")


def parse_progress_file(progress_file_path):
    reading_progress = progress.ReadingProgressBook(
        os.path.basename(progress_file_path), -1, -1)
    try:
        with open(progress_file_path, 'r') as progress_file_object:
            text = progress_file_object.read()
    except FileNotFoundError:
        return None
    else:
        text_fields = text.split('|')
        if len(text_fields) != 3:
            return None
        else:
            reading_progress.from_line = int(text_fields[1])
            reading_progress.to_line = int(text_fields[2])
            return reading_progress


def goto_beginning_book(bookwnd_nav):
    bookwnd_nav.from_line = 0
    bookwnd_nav.to_line = bookwnd_nav.window_height
    bookwnd_nav.current_row = 0
    bookwnd_nav.line_number = 1


def adjust_book_position_toline(bookwnd_nav, to_line):
    bookwnd_nav.from_line = to_line + 1
    bookwnd_nav.to_line = bookwnd_nav.from_line + bookwnd_nav.window_height
    bookwnd_nav.line_number = bookwnd_nav.from_line
    bookwnd_nav.current_row = 0


def open_notes_file(PROGRAM_NOTES_PATH_DIR, filename):
    notes_file = os.path.join(PROGRAM_NOTES_PATH_DIR,
                              os.path.basename(filename))
    if sys.platform == 'linux':
        if shutil.which('xterm') and shutil.which('vim'):
            subprocess.call(
                ["/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", '+$', notes_file])
    else:
        if shutil.which('notepad'):
            subprocess.call(['notepad', notes_file])


def open_url_in_browser(url_search_type):
    if shutil.which('xdg-open'):
        clipboard_content = clipboard.paste()
        if clipboard_content:
            url_search = url_search_type.value.format(clipboard_content)
            subprocess.call(
                ["xdg-open", url_search])


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

        # Initialize stuff ...
        progress_file = get_progress_filepath(filename)
        if os.path.exists(progress_file):
            reading_progress = parse_progress_file(progress_file)
            if reading_progress:
                bookwnd_nav.from_line = reading_progress.from_line
                bookwnd_nav.to_line = reading_progress.to_line
                bookwnd_nav.line_number = reading_progress.from_line + 1
                bookwnd_nav.current_row = 0

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
                    goto_line = utils.go_to(bookwnd_nav, int(input_goto))
                    adjust_book_position_toline(bookwnd_nav, goto_line)
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

            elif key in SAVE_PROGRESS_KEY_CODE:
                save_progress(filename, bookwnd_nav)
                print_save_progress_status(stdscr, bookwnd_nav, filename)

            elif key == curses.KEY_HOME:
                goto_beginning_book(bookwnd_nav)

            elif key == curses.KEY_END:
                goto_line = utils.go_to(
                    bookwnd_nav, bookwnd_nav.book_number_lines())
                adjust_book_position_toline(bookwnd_nav, goto_line)

            elif key in ADD_NOTES_KEY_CODE:
                open_notes_file(PROGRAM_NOTES_PATH_DIR, filename)

            elif key in OPEN_RAE_WEBSITE_KEY_CODES:
                open_url_in_browser(book.URLSearch.rae)

            elif key in OPEN_GOODREADS_WEBSITE_KEY_CODES:
                open_url_in_browser(book.URLSearch.good_reads)

            elif key in WORD_BUILDING_KEY_CODES:
                book_page = book_chunk(lines, bookwnd_nav.from_line,
                                       bookwnd_nav.to_line, bookwnd_nav.book_number_lines())
                if book_page:
                    stdscr.clear()
                    row_line = book_page[bookwnd_nav.current_row]
                    stdscr.attron(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))
                    words = row_line.split()
                    stdscr.addstr(0, 0, row_line)
                    stdscr.attroff(curses.color_pair(STATUSBAR_COLOR_PAIRCODE))
                    stdscr.getch()

            if bookwnd_nav.window_mode == book.WindowMode.reading:
                stdscr.clear()
                print_page(stdscr, lines, bookwnd_nav)


# Create program's directories:
Path(PROGRAM_PATH_DIR).mkdir(parents=True, exist_ok=True)
Path(PROGRAM_PROGRESS_PATH_DIR).mkdir(parents=True, exist_ok=True)
Path(PROGRAM_NOTES_PATH_DIR).mkdir(parents=True, exist_ok=True)
Path(PROGRAM_WORDS_PATH_DIR).mkdir(parents=True, exist_ok=True)

curses.wrapper(main)
