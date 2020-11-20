import curses

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
ENTER_KEY_CODES = [curses.KEY_ENTER, 10, 13]
