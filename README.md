# go-reloaded

A Go command-line tool that edits text for you: it reads a file, applies inline modifiers such as `(up)` or `(hex)`, fixes punctuation, quotes and articles, and writes the cleaned result to another file.

---

## Description

`go-reloaded` is a small text-completion and auto-correction tool. Inside the input text you place modifiers like `(cap)` or `(bin)`; the program applies each modifier to the word before it, removes the modifier, and then normalises the punctuation, the single quotes and the `a` / `an` articles. It is written in Go using only the standard library.

---

## Features

- Case modifiers: `(up)`, `(low)`, `(cap)`, with an optional word count such as `(up, 3)`
- Number conversion: `(hex)` and `(bin)` to decimal
- Punctuation fixed: `. , ! ? : ;` attach to the previous word, groups like `...` and `!?` stay together
- Single quotes fixed: `' word '` becomes `'word'`
- Articles fixed: `a` becomes `an` before a vowel or `h`, and `an` becomes `a` otherwise
- Each line is processed independently; blank lines are kept
- Standard library only

---

## Usage

```bash
go run . <input.txt> <output.txt>
```

**Example:**

```bash
go run . input.txt result.txt
```

Both arguments are required. The output file is created, or overwritten if it already exists.

---

## Modifiers

A modifier applies to the word (or words) written **before** it. Write the count with a space after the comma: `(up, 2)`.

| Modifier | Effect |
|----------|--------|
| `(up)` | Previous word to UPPERCASE |
| `(low)` | Previous word to lowercase |
| `(cap)` | Previous word capitalised (first letter upper, rest lower) |
| `(up, N)` | Previous `N` words to UPPERCASE |
| `(low, N)` | Previous `N` words to lowercase |
| `(cap, N)` | Previous `N` words capitalised |
| `(hex)` | Previous word, read as hexadecimal, replaced by its decimal value |
| `(bin)` | Previous word, read as binary, replaced by its decimal value |

If the word before `(hex)` or `(bin)` is not a valid number, it is left unchanged and the modifier is still removed. If `N` is larger than the number of words available, the modifier applies to all of them. A modifier with no word before it is ignored.

### Automatic fixes

| Fix | Rule |
|-----|------|
| Punctuation | `. , ! ? : ;` are attached to the previous word and followed by a space, even if the input has them glued to the next word |
| Quotes | Standalone `'` marks are paired: the opening one attaches to the next word, the closing one to the previous word |
| Articles | `a` before a word starting with `a e i o u h` becomes `an`; `an` before any other letter becomes `a`; the capital letter of the article is preserved |
| Spacing | Repeated spaces collapse to one |

---

## Examples

| Input | Output |
|-------|--------|
| `1E (hex) files were renamed` | `30 files were renamed` |
| `It has been 10 (bin) years` | `It has been 2 years` |
| `Ready, set, go (up) !` | `Ready, set, GO!` |
| `This is so exciting (up, 2)` | `This is SO EXCITING` |
| `it (cap) was the best of times (cap, 6)` | `It Was The Best Of Times` |
| `I should stop coming so close to the fire , he said ...` | `I should stop coming so close to the fire, he said...` |
| `I was sitting over there ,and then BAMM !!` | `I was sitting over there, and then BAMM!!` |
| `There is no greater agony than bearing a untold story inside you.` | `There is no greater agony than bearing an untold story inside you.` |
| `It was a honest mistake` | `It was an honest mistake` |
| `' I am a optimist ' and ' awesome '` | `'I am an optimist' and 'awesome'` |

**Whole file** (`input.txt` → `result.txt`):

```text
harold wilson (cap, 0) : ' I am a optimist ,but a optimist who carries a raincoat . '

G (hex)
'a omar'
jhj,,, on
```

```text
harold wilson: 'I am an optimist, but an optimist who carries a raincoat.'

G
'a omar'
jhj,,, on
```

---

## Processing Flow

```
input.txt
   │  read the file, split it into lines
   ▼
for each line  ->  fixLine
   │
   ├─ strings.Fields ...... split the line into words
   ├─ processWords ........ apply (hex) (bin) (up) (low) (cap) (up, N) (low, N) (cap, N)
   ├─ handleArticles ...... a <-> an
   ├─ fixPunctuation ...... attach . , ! ? : ; to the previous word
   ├─ fixQuotes ........... ' word '  ->  'word'
   └─ strings.Join ........ words back into one line
   │
   ▼
join the lines, write the file
   │
   ▼
output file
```

The order matters: modifiers run first, so `(up)` and friends never reach the later steps, and articles are checked before the punctuation is separated from the words.

---

## Project Structure

```
go-reloaded/
├── main.go        # all the logic: modifiers, articles, punctuation, quotes, file I/O
├── go.mod
├── input.txt      # sample input
└── result.txt     # output produced from input.txt
```

---

## Implementation Details

All the logic lives in `main.go`, split into small single-purpose functions:

| Function | Responsibility |
|----------|----------------|
| `main` | Validates the arguments, reads the input file, processes it line by line, writes the output file |
| `fixLine` | Runs the whole pipeline on one line; returns an empty string for a blank line |
| `processWords` | Walks the words; applies the single-word modifiers and calls `applyTransform` for the `(tag, N)` form |
| `applyTransform` | Applies `up`, `low` or `cap` to the previous `N` words, stopping at the start of the line |
| `capitalize` | First letter to upper case, the rest to lower case |
| `handleArticles` | Corrects `a` / `an` according to the first letter of the next word |
| `fixPunctuation` | Splits punctuation from words, then attaches each mark to the previous word |
| `fixQuotes` | Pairs standalone `'` marks and attaches them to the neighbouring words |

Input is split with `strings.Fields`, so a modifier such as `(up, 2)` arrives as two words, `(up,` and `2)`, which `processWords` reads together.

---

## Error Handling

| Situation | Behaviour |
|-----------|-----------|
| Wrong number of arguments | Prints `Usage: go run . <input.txt> <output.txt>` and exits with code 1 |
| Input file missing or unreadable | Prints `error of reading file: ...` and exits with code 1 |
| Output file cannot be written | Prints `error of writing file: ...` and exits with code 1 |
| Empty input file | Writes an empty output file |
| Windows line endings (`\r\n`) | Converted to `\n` in the output |

---

## Known Limitations

- The count must be written with a space: `(up, 2)` works, `(up,2)` is left as text.
- Modifiers are case-sensitive: `(UP)` is not recognised.
- A dot, comma or similar mark inside a word is split from it: `3.14` becomes `3. 14` and `example.com` becomes `example. com`.
- `(cap)` works on bytes, so a word that starts with a non-ASCII letter (for example `élan`) is corrupted.

---

## Author

**Omar Rady** — [@orady](https://01.nextera.education/git/orady)
