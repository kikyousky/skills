# Repository Instructions

This repository stores reusable skills and commands for coding and daily work.

## Project Structure(first level folder only)

```text
.
├── AGENTS.md
├── README.md
├── .gitignore
├── commands/
├── e2e-testing-develop/
├── e2e-testing-setup/
├── golang-backend/
└── spec-authoring/
```

## Skills

- `spec-authoring`: Use before implementation for new apps, substantial features, risky bug fixes, broad refactors, multi-file behavior changes, API/contract changes, dependency/config changes with runtime impact, or work needing non-trivial validation. Skip for small localized edits, documentation-only tweaks, command/skill text updates, and obvious single-file maintenance. For complex changes, write specs to the target project at `docs/specs/{YYYY-MM-DD}-{slug}.md`, approve them through `/spec-review`, and implement them through `/spec-implement`.
- `golang-backend`: Use for Go backend architecture, package boundaries, validation, JSON decoding, config access, naming, and complexity checks.
- `e2e-testing-setup`: Use when setting up or running Playwright E2E tests in an isolated Docker Compose environment.
- `e2e-testing-develop`: Use when developing or debugging Playwright tests, Page Object Models, traces, artifacts, flaky tests, and reports.

## Standard Skill Format

All skills in this repository must conform to the Agent Skills specification at `https://agentskills.io/home`. Treat that site as the authoritative standard for skill structure, metadata, and behavior. Do not introduce repo-local skill conventions that conflict with it.

Each skill directory must be named after the skill and contain a `SKILL.md` file:

```text
<skill-name>/
  SKILL.md
  references/        # optional supporting docs
  scripts/           # optional helper scripts
  assets/            # optional templates or static files
```

`SKILL.md` must use this frontmatter shape:

```markdown
---
name: skill-name
description: Concrete trigger-oriented description of when to use this skill.
---

# Skill Title

Skill instructions, workflows, examples, and references.
```

- `name` is required, lowercase, hyphen-separated, and must match the folder name. Max 64 characters. 
- `description` is required and must state both what the skill does and when to trigger it. Max 1024 characters. 

## How to write a great skill

1. make the skill descriptions a little bit "pushy" to avoid "undertrigger" the skill. For instance, instead of "How to build a simple fast dashboard to display internal Anthropic data.", you might write "How to build a simple fast dashboard to display internal Anthropic data. Make sure to use this skill whenever the user mentions dashboards, data visualization, internal metrics, or wants to display any kind of company data, even if they don't explicitly ask for a 'dashboard.'"
2. Keep the content of skill as lean as possible. Remove things that aren't pulling their weight.
3. Keep SKILL.md under 500 lines; if you're approaching this limit, use reference files, but you must reference files clearly with guidance on when to read them
4. When a skill supports multiple domains/frameworks, organize by variant, models will read only the relevant reference file:

cloud-deploy/
├── SKILL.md (workflow + selection)
└── references/
    ├── aws.md
    ├── gcp.md
    └── azure.md

5. For large reference files (>300 lines), include a table of contents

## Standard Command Format

Commands live in `commands/` so I can copy them elsewhere.

Commands file should use this format:

```markdown
---
description: Short sentence describing what the command does.
---

# Command Name

Command instructions.
```


## Ignore Rules

Lark skills are intentionally ignored by Git through `/lark-*/`.

## Maintenance Rules

- Preserve existing user changes in this repository; do not revert unrelated modified or untracked files.
