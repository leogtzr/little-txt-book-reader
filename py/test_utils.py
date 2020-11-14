import unittest
from unittest import TestCase
import utils


class UtilFunctionTestCase(TestCase):
    """Tests for 'utils.py'."""

    def test_x(self):
        """ doc """
        # got = ...
        # self.assertEqual(got, expected)
        pass

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


if __name__ == '__main__':
    unittest.main()
