import os


class ReadingProgressBook:
    """ abc """

    def __init__(self, filename, from_line, to_line):
        self.filename = filename
        self.from_line = from_line
        self.to_line = to_line
