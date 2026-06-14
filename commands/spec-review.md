---
description: Review a drafted spec and approve or reject it before implementation.
---

# /spec-review

Use this command to review a drafted spec under `docs/specs/` and either reject it with concrete remediation notes or approve it by updating frontmatter.

## Purpose

Provide the only supported approval gate between spec drafting and implementation start.

## Input

Accept one existing spec path under `docs/specs/`.

Reject immediately if:

- the path is outside `docs/specs/`
- the file does not exist
- the file does not follow the required frontmatter and section structure

## Review Outcomes

This command has only two valid outcomes:

1. **Reject**
   - leave the spec unapproved
   - provide concrete remediation notes
   - identify missing, vague, or conflicting sections

2. **Approve**
   - change frontmatter `status` to `approved`
   - set `approved_at` to the approval timestamp
   - preserve the rest of the spec except for any required review notes explicitly requested by the operator

## Deterministic Checklist

Review the spec in this order:

### 1. Frontmatter validity

Verify these keys exist:

- `title`
- `slug`
- `status`
- `request_type`
- `created_at`
- `approved_at`
- `source_of_truth`
- `implementation_mode`

Verify `status` is one of:

- `draft`
- `reviewed`
- `approved`

### 2. Section structure

Verify the spec body uses this exact structure:

- `Purpose / Big Picture`
- `Context and Orientation`
- `Implementation Impact`
- `Execution`
  - `Progress`
  - `Decision Log`
  - `Plan of Work`
  - `Concrete Steps`
- `Validation`
  - `Test Plan`
  - `Test Results`

### 3. Purpose quality

Verify the purpose explains:

- why the work matters
- what new behavior exists after the change
- how to observe the result

### 4. Context quality

Verify the context section names relevant repository paths and gives enough orientation for a novice.

### 5. Implementation impact coverage

Verify all five categories are present:

- `Database / Schema`
- `GUI / UX`
- `API / Contracts`
- `Config / Environment`
- `Dependencies`

Rules:

- if a category is unaffected, it must say `none`
- `GUI / UX` may stay high-level
- `Database / Schema`, `API / Contracts`, `Config / Environment`, and `Dependencies` must be concrete and specific

Reject the spec if non-GUI impact sections are vague, generic, or hand-wavy.

### 6. Scope boundaries

Verify the spec states what is in scope and what is intentionally out of scope.

Rules:

- the boundary may live in `Purpose / Big Picture`, `Context and Orientation`, or `Execution > Plan of Work`
- included work must be concrete enough that an implementer can tell what belongs in this change
- excluded work must be explicit enough to prevent silent scope expansion during implementation

Reject the spec if the intended change could expand silently because the boundaries are missing, vague, or contradictory.

### 7. Execution quality

Verify:

- `Progress` exists and reflects a plausible draft/review state
- `Decision Log` contains planning decisions when meaningful choices were already made
- `Plan of Work` describes the intended implementation strategy in prose
- `Concrete Steps` includes exact commands or explicitly states what must be run once the real repository commands are confirmed

### 8. Validation quality

Verify:

- `Validation > Test Plan` exists
- the test plan includes automated tests and agent-executed QA where relevant
- expected behavior or outputs are stated
- `Validation > Test Results` exists even if validation has not run yet

If validation has already run, verify `Test Results` includes concise proof such as short transcripts, logs, snippets, screenshots, or URLs.

## Approval Rules

Approve only when all checklist items pass.

Do not approve partial specs.
Do not approve specs with missing sections.
Do not approve specs with vague non-GUI impact categories.
Do not approve specs with missing or ambiguous scope boundaries.
Do not approve specs without a real `Validation > Test Plan`.
Do not approve specs that leave room for silent scope expansion during implementation.

## Output Style

For rejection, report:

- overall verdict: `REJECT`
- each failing checklist item
- the exact remediation needed

For approval, report:

- overall verdict: `APPROVE`
- confirmation that `status` was changed to `approved`
- confirmation that `approved_at` was populated
