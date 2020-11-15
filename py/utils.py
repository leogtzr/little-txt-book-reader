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
