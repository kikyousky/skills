---
description: Summarize the current session into a Markdown file.
---

# /summary

Summarize the current session and write the result to a Markdown file.

## Input

Accept an optional Markdown file path.

## Behavior

- If a file is specified, merge the current session summary into that file.
- If no file is specified, create a new Markdown summary file.
- For long sessions, identify the main ideas and section titles first, then fill in each section one by one.
- Use `TodoWrite` to track progress when necessary.
