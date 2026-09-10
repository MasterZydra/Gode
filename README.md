# Gode

A lightweight text editor built with Go and powered by the Fyne GUI toolkit.

The long-term goal is to evolve Gode into a lightweight IDE.

## Code highlighting features
The editor has built-in code highlighting for
- Go
- Dockerfile

## Shortcuts and editing features

### Commands
- Ctrl+Shift+P opens the command palette.
- Format Document formats the current Go or JSON file.

### Selection and cursors
- Shift + arrow keys extend the current selection.
- Shift + click extends the active selection.
- Alt + click adds another cursor.

### Clipboard and editing
- Ctrl+A selects the entire file.
- Ctrl+C copies the selected text
- Ctrl+X cuts the selected text
- Ctrl+V pastes at all active cursors.
- Ctrl+V after copying a line inserts it directly below the current line.
- Ctrl+V works for the final line even without a trailing newline.

### Word navigation and selection
- Ctrl+Left/Right jumps between word boundaries.
- Ctrl+Shift+Left/Right selects to the previous/next word boundary.
- Double-clicking a word selects the entire word.
