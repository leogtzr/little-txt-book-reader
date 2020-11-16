import unittest
import utils
import book


class UtilFunctionTestCase(unittest.TestCase):
    """Tests for 'utils.py'."""

    def test_percent(self):
        param_list = [(150, 30, 20.0)]
        for total_lines, current_index, want in param_list:
            with self.subTest():
                got = utils.percent(current_index, total_lines)
                self.assertEqual(got, want, f"got={got}, expected={want}")

    def test_lines_to_change_percentage_point(self):
        param_list = [(100, 1000, 10), (22, 1000, 8)]
        for current_line, total_lines, want in param_list:
            with self.subTest():
                got = utils.lines_to_change_percentage_point(
                    current_line, total_lines)
                self.assertEqual(got, want, f"got={got}, expected={want}")

    def test_go_to(self):
        param_list = [(2000, 53, 120, 106)]

        for book_number_of_lines, window_height, goto_line, want in param_list:
            bookwnd_nav = book.BookWindowNavigation(
                book_number_of_lines, window_height, -1)
            got = utils.go_to(bookwnd_nav, goto_line)
            self.assertEqual(got, want, f"got={got}, expected={want}")


if __name__ == '__main__':
    unittest.main()
