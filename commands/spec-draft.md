---
description: Draft a repository-aware executable spec before implementation begins.
---

# /spec-draft

Use this command to start spec creation for a target project before any implementation begins.

## Purpose

Turn a raw request into a repository-aware executable spec draft using the `spec-authoring` workflow.

## Required Behavior

1. Read the target repository before drafting.
2. Use the `spec-authoring` skill and its bundled `assets/template.md`.
3. Derive a lowercase kebab-case slug from the request.
4. Write the output spec to:

   `docs/specs/{YYYY-MM-DD}-{slug}.md`

5. Keep the resulting frontmatter status at `draft` or `reviewed`.
6. Stop after draft creation. Do not approve the spec. Do not start implementation.

## Input Expectations

The operator should provide:

- the target project or repository path if it is not obvious
- the raw request
- any known constraints, standards, or non-goals

If details are missing, inspect the repository first and convert minor gaps into explicit assumptions instead of blocking.

## Non-Negotiable Rules

- Do not write the spec anywhere outside `docs/specs/` in the target project.
- Do not skip repository inspection.
- Do not omit implementation-impact categories. Use `none` where appropriate.
- Do not leave non-GUI impact categories vague.
- Do not mark the spec `approved`.
- Do not invoke `/spec-start`, `/start-work`, or any implementation flow.

## Output Contract

The command succeeds only when:

- exactly one dated spec file is created or updated in `docs/specs/`
- the file follows the `spec-authoring` bundled `assets/template.md`
- the file is self-contained and reviewable
- the status is `draft` or `reviewed`
