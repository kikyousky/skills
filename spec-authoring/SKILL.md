---
name: spec-authoring
description: Author a self-contained executable spec before implementation. Use whenever a request involves a new app, feature, bug fix, refactor, multi-file change, behavior change, contract/API change, config/dependency change, or non-trivial validation, even if the user does not explicitly ask for a spec.
---

# Spec Authoring

Write a single executable spec that becomes the source of truth for drafting, review, approval, implementation, and validation.

Bundled files:

- `references/example-spec.md`: filled example; read when the target spec shape or level of detail is unclear.
- `assets/template.md`: exact spec template; use as the required output structure.

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

The spec must be detailed enough that implementation can begin from the approved file without depending on hidden chat context.

## Required Workflow

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

The file may end in either of these states:

- `status: draft`
- `status: reviewed`

`reviewed` means the spec has been refined and is ready for approval review, but it is still not executable.

It must not end in `status: approved` unless a separate approval step explicitly promotes it.

## Section-Specific Instructions

### Purpose / Big Picture

Write from the user or operator perspective. Explain what new behavior becomes possible and how to observe it working.

State what is in scope and what is intentionally out of scope clearly enough to prevent silent scope expansion during implementation.

### Context and Orientation

Name the relevant repository files and directories using repository-relative paths. Explain enough structure that a novice can navigate confidently.

### Implementation Impact

This section is mandatory and must always answer all five categories.

#### Database / Schema

Be explicit about tables, columns, indexes, constraints, migrations, seeds, and data backfills.

#### GUI / UX

This may stay high-level. Describe the screens, flows, or visible behavior that change.

#### API / Contracts

Be explicit about endpoints, request and response shapes, jobs, events, service contracts, or exported interfaces.

#### Config / Environment

Be explicit about environment variables, feature flags, secrets, deployment config, build config, or external service configuration.

#### Dependencies

Be explicit about added, removed, or upgraded packages, SDKs, services, or infrastructure dependencies, and explain why they are needed.

If any category is unaffected, write `none`.

For implementation impact, GUI may be high-level, but schema, API/contracts, config/env, and dependencies must be concrete and specific.

### Execution > Progress

Seed this with the template checklist and adjust it if the task needs more precise draft-time tracking.

### Execution > Decision Log

Record important planning decisions made during authoring. Do not wait until implementation if you already chose an approach.

### Execution > Plan of Work

Describe the intended implementation strategy in prose. Name exact files or locations to edit or create whenever they are already knowable from repo inspection.

### Execution > Concrete Steps

Write the expected runtime procedure as exact commands plus working directory and expected result. If command names are unknown, inspect the repository and use the actual commands.

### Validation > Test Plan

This is mandatory. Include:

- automated tests
- agent-executed QA scenarios
- expected outputs or observed behavior
- before/after behavior when relevant

### Validation > Test Results

For a newly drafted spec, this section may begin empty or note that validation has not run yet. It must remain present because the same file will later hold concise proof that validation succeeded or failed.

## Assumptions and Unknowns

If you must proceed with an assumption, make it explicit in the relevant section instead of hiding it.

If the assumption affects approval, say so clearly.

## Hard Constraints

- Do not start implementation.
- Do not approve the spec.
- Do not write the spec outside `docs/specs/` in the target project.
- Do not omit any required frontmatter field.
- Do not omit any `Implementation Impact` category.
- Do not leave non-GUI impact categories vague.
- Do not omit `Validation > Test Plan`.
- Do not remove `Validation > Test Results`; it is required even if results are not available yet.

## Handoff

After drafting or updating the spec, report the spec path and direct the operator to run `/spec-review <path>` for approval. Implementation must wait until `/spec-start <path>` is run against a spec whose frontmatter has `status: approved` and a populated `approved_at`.

## Completion Standard

The skill is successful when the resulting spec:

- exists at `docs/specs/{YYYY-MM-DD}-{slug}.md`
- matches this skill's bundled `assets/template.md`
- is self-contained and novice-readable
- includes explicit implementation impact
- includes both `Validation > Test Plan` and `Validation > Test Results`
- ends at `draft` or `reviewed`, not `approved`
