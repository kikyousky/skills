---
title: <short action-oriented title>
slug: <kebab-case-slug>
status: draft
request_type: <new-app|feature|bugfix|refactor|other>
created_at: <YYYY-MM-DDTHH:MM:SSZ>
approved_at: null
source_of_truth: true
implementation_mode: spec-first
---

# <Short action-oriented title>

This spec is a living document. During implementation, `Execution > Progress`, `Execution > Decision Log`, `Execution > Concrete Steps`, and `Validation > Test Results` are the routine-update areas. All other sections are change-controlled and should be updated only when implementation facts actually change.

## Purpose / Big Picture

Explain why this work matters from the user or operator perspective. Describe what will be possible after the change that is not possible now. State how to observe the result working in a concrete way.

State what is in scope and what is intentionally out of scope. The boundary must be concrete enough to prevent adjacent work from being added silently during implementation.

## Context and Orientation

Describe the relevant repository area as if the implementer knows nothing. Name important files, directories, modules, routes, commands, schemas, or services using repository-relative paths. Define any non-obvious terms in plain language.

## Implementation Impact

### Database / Schema

State whether this work changes tables, columns, indexes, constraints, migrations, seeds, or data backfills. If none, write `none`.

### GUI / UX

State whether this work changes screens, layouts, forms, components, user flows, loading states, error states, or other visible behavior. High-level description is acceptable. If none, write `none`.

### API / Contracts

State whether this work changes endpoints, request or response payloads, events, background jobs, internal service contracts, or exported interfaces. Be specific. If none, write `none`.

### Config / Environment

State whether this work changes environment variables, feature flags, secrets, deployment configuration, build settings, or third-party service configuration. Be specific. If none, write `none`.

### Dependencies

State whether this work adds, removes, or upgrades libraries, SDKs, packages, services, or infrastructure dependencies. Name the exact dependency and why it is needed. If none, write `none`.

## Execution

### Progress

- [ ] Initial repo inspection completed.
- [ ] Spec drafted with scope, impact, and validation.
- [ ] Spec reviewed.
- [ ] Spec approved.
- [ ] Implementation started from approved spec.
- [ ] Validation completed.

### Decision Log

- Decision: <decision made>
  Rationale: <why this path was chosen>
  Date/Author: <date / author>

### Plan of Work

Describe the intended implementation strategy in prose. Name the exact files or locations to edit or create and what will change in each place. Keep the sequence concrete and minimal.

### Concrete Steps

State the actual commands to run, the working directory, and the expected output or observable effect.

Example format:

  Working directory: <repo-root>
  Command: <exact command>
  Expected: <specific output or observable effect>

## Validation

### Test Plan

Describe what must be tested before the work is complete.

Include:

- automated test commands
- manual QA scenarios
- expected outputs, UI states, API responses, or error messages
- before/after behavior when relevant

### Test Results

Capture what was actually run and the concise proof that validation passed or failed.

Examples of acceptable evidence:

- short terminal transcript
- API request/response snippet
- screenshot path or UI note
- concise diff excerpt
- log line proving success or failure
