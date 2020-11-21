import curses
import os
from curses import textpad
from book import WindowMode
from book import NavigationMode
from constants import STATUSBAR_COLOR_PAIRCODE
from constants import ENTER_KEY_CODES
from constants import PROGRAM_WORDS_PATH_DIR
from constants import KEY_ESCAPE_CODE
from math import ceil


def check_key_down(bookwnd_nav):
    if bookwnd_nav.nav_mode() == NavigationMode.by_page:
        if bookwnd_nav.current_row == (bookwnd_nav.window_height - 1):
            # Reset:
            bookwnd_nav.current_row = 0
            bookwnd_nav.line_number += 1
            bookwnd_nav.from_line += bookwnd_nav.window_height
            bookwnd_nav.to_line = bookwnd_nav.from_line + bookwnd_nav.window_height
        else:
            bookwnd_nav.current_row += 1
            bookwnd_nav.line_number += 1
    else:
        if bookwnd_nav.line_number < bookwnd_nav.book_number_lines():
            bookwnd_nav.current_row = bookwnd_nav.window_height - 1
            bookwnd_nav.line_number += 1
            bookwnd_nav.from_line += 1
            bookwnd_nav.to_line += 1


def check_key_up(bookwnd_nav):
    if bookwnd_nav.nav_mode() == NavigationMode.by_page:
        if bookwnd_nav.current_row == 0:
            if bookwnd_nav.line_number > bookwnd_nav.window_height:
                bookwnd_nav.current_row = bookwnd_nav.window_height - 1
                bookwnd_nav.line_number -= 1
                bookwnd_nav.from_line -= bookwnd_nav.window_height
                bookwnd_nav.to_line -= bookwnd_nav.window_height
        else:
            bookwnd_nav.current_row -= 1
            bookwnd_nav.line_number -= 1
    else:
        if bookwnd_nav.line_number > 0:
            bookwnd_nav.current_row = bookwnd_nav.window_height - 1
            bookwnd_nav.line_number -= 1
            bookwnd_nav.from_line -= 1
            bookwnd_nav.to_line -= 1


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
        'R       -> Opens GoodReads Web site with search from the clipboard.',
        'W       -> Open Word Building Mode with current sentence',
        'V       -> View Words added to the Word Building Database',
        'T       -> View Stats'
    ]

    for idx, help_entry in enumerate(help_entries):
        stdscr.addstr(border_offset + idx + 1, border_offset+1, help_entry)


def word_building_words(bookwnd_nav):
    base_filename = os.path.basename(bookwnd_nav.filename)
    try:
        with open(os.path.join(PROGRAM_WORDS_PATH_DIR, base_filename), 'r') as f:
            lines = [line.rstrip('\n') for line in f.readlines()]
    except FileNotFoundError:
        return []
    else:
        return lines


def word_count(lines):
    count = 0
    for line in lines:
        count += len(line.split())
    return count


def view_stats(lines, bookwnd_nav, stdscr):
    screen_height, screen_width = stdscr.getmaxyx()
    border_offset = 3
    box = [[border_offset, border_offset], [
        screen_height-border_offset, screen_width-border_offset]]
    textpad.rectangle(
        stdscr, box[0][0], box[0][1], box[1][0], box[1][1])
    perc = percent(bookwnd_nav.line_number,
                   bookwnd_nav.book_number_lines())

    stat_entries = [
        f"Progress {perc:.1f}%, line: {bookwnd_nav.line_number}",
        f"Pages:   {bookwnd_nav.book_number_lines()}",
        f"Words:   ~{word_count(lines)}"
    ]

    for idx, entry in enumerate(stat_entries):
        stdscr.addstr(border_offset + idx + 1, border_offset+1, entry)

    bookwnd_nav.window_mode = WindowMode.reading
    stdscr.getch()


def view_words(bookwnd_nav, stdscr):
    bookwnd_nav.window_mode = WindowMode.view_words

    words = word_building_words(bookwnd_nav)

    stdscr = curses.initscr()
    curses.noecho()
    curses.cbreak()
    curses.start_color()
    stdscr.keypad(1)
    curses.init_pair(1, curses.COLOR_BLACK, curses.COLOR_CYAN)
    highlightText = curses.color_pair(1)
    normalText = curses.A_NORMAL
    stdscr.border(0)
    curses.curs_set(0)
    max_row = 10

    if len(words) >= bookwnd_nav.window_height:
        max_row = bookwnd_nav.window_height // 2
    else:
        max_row = len(words)

    box = curses.newwin(max_row + 2, 100, 1, 1)
    box.box()
    row_num = len(words)

    pages = int(ceil(row_num / max_row))
    position = 1
    page = 1

    for i in range(1, max_row + 1):
        if i == position:
            box.addstr(i, 2, str(i) + " - " + words[i - 1], highlightText)
        else:
            box.addstr(i, 2, str(i) + " - " + words[i - 1], normalText)

    stdscr.refresh()
    box.refresh()

    x = stdscr.getch()
    while x not in [KEY_ESCAPE_CODE]:
        if x == curses.KEY_DOWN:
            if page == 1:
                if position < i:
                    position = position + 1
                else:
                    if pages > 1:
                        page += 1
                        position = 1 + (max_row * (page - 1))
            elif page == pages:
                if position < row_num:
                    position += 1
            else:
                if position < max_row + (max_row * (page - 1)):
                    position += 1
                else:
                    page += 1
                    position = 1 + (max_row * (page - 1))
        if x == curses.KEY_UP:
            if page == 1:
                if position > 1:
                    position -= 1
            else:
                if position > (1 + (max_row * (page - 1))):
                    position -= 1
                else:
                    page -= 1
                    position = max_row + (max_row * (page - 1))
        if x == curses.KEY_LEFT:
            if page > 1:
                page -= 1
                position = 1 + (max_row * (page - 1))

        if x == curses.KEY_RIGHT:
            if page < pages:
                page += 1
                position = (1 + (max_row * (page - 1)))

        box.erase()
        stdscr.border(0)
        box.border(0)

        # TODO: create a function for this.
        for i in range(1 + (max_row * (page - 1)), max_row + 1 + (max_row * (page - 1))):
            if i + (max_row * (page - 1)) == position + (max_row * (page - 1)):
                box.addstr(i - (max_row * (page - 1)), 2, str(i) +
                           " - " + words[i - 1], highlightText)
            else:
                box.addstr(i - (max_row * (page - 1)), 2, str(i) +
                           " - " + words[i - 1], normalText)
            if i == row_num:
                break

        stdscr.refresh()
        box.refresh()
        x = stdscr.getch()

    bookwnd_nav.window_mode = WindowMode.reading


def words_with_brackets(words, select_idx, stdscr):
    text_sentence_brackets = ''

    for word_idx, word in enumerate(words):
        text_sentence_brackets += f"[{word}] " if word_idx == select_idx else f"{word} "

    return text_sentence_brackets


def write_to_words_file(bookwnd_nav, word):
    base_filename = os.path.basename(bookwnd_nav.filename)

    with open(os.path.join(PROGRAM_WORDS_PATH_DIR, base_filename), 'a') as word_file:
        word_file.write(f"{word}\n")


def word_building_row_sentence_user_input(bookwnd_nav, stdscr, words):
    words_count = len(words)
    word_select_idx = 0

    if words_count > 0:

        stdscr.addstr(0, 0, words_with_brackets(
            words, word_select_idx, stdscr))

        while True:
            key = stdscr.getch()
            if key in ENTER_KEY_CODES:
                selected_word = words[word_select_idx]
                write_to_words_file(bookwnd_nav, selected_word)
                break
            elif key == curses.KEY_LEFT and word_select_idx > 0:
                word_select_idx -= 1
            elif key == curses.KEY_RIGHT and (word_select_idx < (words_count - 1)):
                word_select_idx += 1

            stdscr.addstr(0, 0, words_with_brackets(
                words, word_select_idx, stdscr))


def word_building_mode(bookwnd_nav, stdscr, lines, filename):
    stdscr.clear()

    book_page = book_chunk(lines, bookwnd_nav.from_line,
                           bookwnd_nav.to_line, bookwnd_nav.book_number_lines())
    if book_page:

        bookwnd_nav.window_mode = WindowMode.word_building
        row_line = book_page[bookwnd_nav.current_row]
        words = row_line.split()

        word_building_row_sentence_user_input(
            bookwnd_nav, stdscr, words)

        bookwnd_nav.window_mode = WindowMode.reading


def book_chunk(lines, from_line, to_line, book_number_of_lines):
    return lines[from_line:to_line]


def percent(current_number_line, total_lines):
    return float(current_number_line * 100.0) / float(total_lines)


def lines_to_change_percentage_point(current_line, total_lines):
    start = current_line
    lines_to_change_percentage = -1
    percentage_with_currentLine = int(percent(current_line, total_lines))
    while True:
        current_line += 1
        next_percentage = int(percent(current_line, total_lines))
        if next_percentage > percentage_with_currentLine:
            lines_to_change_percentage = current_line
            break

    return lines_to_change_percentage - start


def go_to(bookwnd_nav, goto_line):
    number_of_lines = bookwnd_nav.book_number_lines()
    sum = 0
    while (sum < number_of_lines) and (sum < goto_line):
        sum += bookwnd_nav.window_height

    if sum > goto_line:
        sum -= bookwnd_nav.window_height
    return sum
