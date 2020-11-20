def print_words_brackets(words, select_idx):

    for word_idx, word in enumerate(words):
        if word_idx == select_idx:
            print(f"[{word}] ", end='')
        else:
            print(f"{word} ", end='')


text = 'campos de sus antepasados. Jamás se metió en política, así es que ninguna'
words = text.split()

print_words_brackets(words, 2)
print()

print_words_brackets(words, 5)
print()
