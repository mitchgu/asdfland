from pathlib import Path



with Path('wordlists/google-10000-english-no-swears.txt').open() as f:
    wordlist = f.readlines()

wordlist = [w.lower().strip() for w in wordlist if len(w.strip()) <= 8 and len(w.strip()) >= 3][:5000]

with Path('wordlists/google.txt').open('w') as f:
    f.write("\n".join(wordlist))