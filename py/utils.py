import curses
from book import WindowMode
from constants import STATUSBAR_COLOR_PAIRCODE
from constants import ENTER_KEY_CODES


def words_with_brackets(words, select_idx, stdscr):
    text_sentence_brackets = ''

    for word_idx, word in enumerate(words):
        text_sentence_brackets += f"[{word}] " if word_idx == select_idx else f"{word} "

    return text_sentence_brackets


def word_building_row_sentence_user_input(bookwnd_nav, stdscr, words):
    words_count = len(words)
    word_select_idx = 0

    if words_count > 0:

        stdscr.addstr(0, 0, words_with_brackets(
            words, word_select_idx, stdscr))

        while True:
            key = stdscr.getch()
            if key in ENTER_KEY_CODES:
                stdscr.addstr(
                    0, 0, f"Add word: {words[word_select_idx]}                                                                                    ")
                stdscr.getch()
                # TODO: add word to file.
                # break
            elif key == curses.KEY_LEFT and word_select_idx > 0:
                word_select_idx -= 1
            elif key == curses.KEY_RIGHT and word_select_idx < words_count:
                word_select_idx += 1

            stdscr.addstr(0, 0, words_with_brackets(
                words, word_select_idx, stdscr))


def word_building_mode(bookwnd_nav, stdscr, lines):
    stdscr.clear()

    book_page = book_chunk(lines, bookwnd_nav.from_line,
                           bookwnd_nav.to_line, bookwnd_nav.book_number_lines())
    if book_page:

        bookwnd_nav.window_mode = WindowMode.word_building
        row_line = book_page[bookwnd_nav.current_row]
        words = row_line.split()

        word_building_row_sentence_user_input(
            word_building_row_sentence_user_input, stdscr, words)

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
