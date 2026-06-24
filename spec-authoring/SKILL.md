---
name: spec-authoring
description: Author a self-contained executable spec before implementation. Use for new apps, substantial features, risky bug fixes, broad refactors, multi-file behavior changes, API/contract changes, dependency/config changes with runtime impact, or work needing non-trivial validation. Skip for small localized edits, documentation-only tweaks, command/skill text updates, and obvious single-file maintenance.
---

# Spec Authoring

Write a single executable spec that becomes the source of truth for drafting, review, approval, implementation, validation, and retrospective follow-up.
The spec is also a living document that is updated during implementation.

Bundled files:

- `assets/template.md`: exact spec template and inline guidance; use as the required output structure.

## What This Skill Produces

Create or update one spec file in the target project at:

`docs/specs/{YYYY-MM-DD}-{slug}.md`

Where:

- `{YYYY-MM-DD}` is the current date
- `{slug}` is a lowercase kebab-case summary of the request

The file must use the exact frontmatter and section structure from this skill's bundled `assets/template.md`.

## Primary Goal

Produce a self-contained spec that lets a stateless agent understand:

- why the work exists
- what repository area it affects
- what concrete implementation impact it has
- how the work is intended to be executed
- how the work must be validated
- what happened after implementation and what should be learned from it

The spec must be detailed enough that implementation can begin from the approved file without depending on hidden chat context.

## Required Workflow

### 0. Decide whether a spec is warranted

Skip this skill for trivial or low-risk work where a spec would add process overhead without improving correctness, such as:

- small localized edits with obvious acceptance criteria
- documentation-only tweaks
- command or skill text updates
- formatting, typo fixes, or mechanical maintenance

Use the skill when the work has meaningful design uncertainty, user-visible behavior changes, multi-file implementation risk, persistent data impact, external integrations, security implications, or non-trivial validation needs.

### 1. Inspect the target repository first

Before asking the user clarifying questions, inspect the target repository and resolve everything that is discoverable from files, code, config, tests, routes, schemas, docs, or commands.

Always prefer repository evidence over guesswork.

### 2. Ask questions only for non-discoverable blockers

If something is ambiguous but the simplest valid interpretation is safe, state the interpretation and proceed.

Ask clarifying questions only when the ambiguity would materially change:

- scope
- acceptance criteria
- data model
- external integrations
- security or compliance behavior
- user-visible behavior

### 3. Draft the spec using the template

Use this skill's bundled `assets/template.md` as the exact structure. Do not invent extra top-level sections. Do not reorder required sections.

### 4. Stop before implementation

This skill drafts the spec. It must never start coding, never auto-approve the spec, and never invoke implementation.

During authoring, the file must end with:

- `status: draft`

The full lifecycle is `draft -> approved -> completed`.

`/spec-review` is the approval gate: it reviews `draft`, leaves failures as `draft`, and promotes passing specs directly to `approved` with `approved_at` populated. `/spec-start` begins implementation from an `approved` spec without changing status. Implementation changes the spec to `completed` with `completed_at` populated only after validation and `Outcomes & Retrospective` are done.

## Assumptions and Unknowns

If you must proceed with an assumption, make it explicit in the relevant section instead of hiding it.

If the assumption affects approval, say so clearly.

## Hard Constraints

- Do not write the spec outside `docs/specs/` in the target project.
- Do not omit any required frontmatter field.
- Do not omit any `Implementation Impact` category.
- Do not leave non-GUI impact categories vague.
- Do not omit `Validation > Test Plan`.
- Do not remove `Validation > Test Results`; it is required even if results are not available yet.
- Do not omit `Outcomes & Retrospective`; it is required even if implementation has not started yet.

## Handoff

After drafting or updating the spec, report the spec path and direct the operator to run `/spec-review <path>` for approval. Implementation must wait until `/spec-start <path>` is run against a spec whose frontmatter has `status: approved` and a populated `approved_at`.
