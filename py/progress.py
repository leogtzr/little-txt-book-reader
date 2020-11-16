import os


class ReadingProgressBook:
    """ abc """

    def __init__(self, filename, from_line, to_line):
        # TODO:
        pass
        self.filename = filename
        self.from_line = from_line
        self.to_line = to_line

    # TODO: ...
    @staticmethod
    def from_profress_file(filename):
        base_filename = os.path.basename(filename)
        #progress_book = ReadingProgressBook()
