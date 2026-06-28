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

5. Keep the resulting frontmatter status at `draft`.
6. Stop after draft creation. Do not approve the spec. Do not start implementation.

## Required Frontmatter

The draft must include exactly the frontmatter keys from `spec-authoring/assets/template.md`:

- `title`
- `slug`
- `status`
- `request_type`
- `created_at`
- `approved_at`
- `completed_at`

For a new draft, use `status: draft`, `approved_at: null`, and `completed_at: null`.

## Required Structure

The draft must follow the current template section order exactly:

- `Purpose / Big Picture`
- `Implementation Impact`
  - `Database / Schema`
  - `GUI / UX`
  - `API / Contracts`
  - `Config / Environment`
  - `Dependencies`
- `Execution`
  - `Progress`
  - `Surprises & Discoveries`
  - `Decision Log`
  - `Plan of Work`
  - `Concrete Steps`
- `Validation`
  - `Test Plan`
  - `Test Results`
- `Outcomes & Retrospective`

Replace the template's example content with target-project facts. Preserve the integrated guidance style only when it helps the spec remain readable; do not leave generic placeholders as final content.

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
- Do not omit `Execution > Surprises & Discoveries`; write `none yet` if no discoveries exist in a draft.
- Do not omit `Outcomes & Retrospective`; write that implementation has not started yet for a draft.
- Do not mark the spec `approved`.
- Do not invoke `/spec-implement`, `/start-work`, or any implementation flow.

## Output Contract

The command succeeds only when:

- exactly one dated spec file is created or updated in `docs/specs/`
- the file follows the `spec-authoring` bundled `assets/template.md`
- the file is self-contained and reviewable
- the status is `draft`
- `approved_at` and `completed_at` are `null`
