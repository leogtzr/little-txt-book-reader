from enum import Enum


class WindowMode(Enum):
    reading = 1
    help = 2
    goto = 3
    word_building = 4


class URLSearch(Enum):
    rae = "https://dle.rae.es/{}"
    good_reads = 'https://www.goodreads.com/search?q={}'


class BookWindowNavigation:
    '''This class will contain everything related with the object navigation
        Current page number, navigation mode (help, reading), etc
    '''

    def __init__(self, book_number_lines, window_height, window_width, filename):
        self._book_number_of_lines = book_number_lines
        self.window_height = window_height
        self.window_width = window_width
        self.from_line = 0
        self.to_line = window_height
        self.current_row = 0
        self.line_number = 1
        self.window_mode = WindowMode.reading
        self.show_status_bar = True
        self.show_percentage_points = False
        self.filename = filename

    def book_number_lines(self):
        return self._book_number_of_lines
