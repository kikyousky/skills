---
description: Implement an approved spec as the source of truth.
---

# /spec-implement

Use this command to begin implementation from an approved spec.

## Purpose

This is the explicit handoff from spec approval into implementation. It turns an approved spec into the single source of truth for execution.

## Input

Accept one spec path under `docs/specs/`.

## Start Preconditions

Refuse to implement unless all of the following are true:

- the file exists
- the file path is under `docs/specs/`
- frontmatter `status` is exactly `approved`
- frontmatter `approved_at` is populated
- frontmatter `completed_at` is `null`

If any precondition fails, stop and direct the operator back to `/spec-review`.

When all preconditions pass, implementation may begin. Leave `status` as `approved` and `completed_at` as `null` until validation and retrospective work are done.

## Handoff Rules

Once the preconditions pass, treat the approved spec as the single source of truth for implementation.

The implementation agent must:

- follow the approved scope
- avoid inventing new scope
- avoid bypassing the spec with hidden chat assumptions
- treat unresolved questions as blockers unless the spec already records a safe assumption
- keep the spec current as a living document while implementation proceeds
- record unexpected behaviors, bugs, optimizations, or insights in `Execution > Surprises & Discoveries` with concise evidence
- record course-changing decisions in `Execution > Decision Log` with short evidence snippets when optimizer behavior, performance tradeoffs, unexpected bugs, or inverse/unapply semantics affect the approach

## Routine-Update Areas During Implementation

During implementation, routine updates go to:

- `Execution > Progress`
- `Execution > Surprises & Discoveries`
- `Execution > Decision Log`
- `Execution > Concrete Steps`
- `Validation > Test Results`
- `Outcomes & Retrospective`

All other sections are change-controlled and should be edited only when implementation facts, validated decisions, commands, or acceptance expectations actually change. Update `Outcomes & Retrospective` at major milestones or completion to compare the result against the original purpose and capture gaps or lessons learned.

## Ralph-Loop Style Execution Behavior

After implementation starts, execution should proceed milestone by milestone without asking for trivial continuation approval between normal steps.

The implementation agent should:

1. read the approved spec
2. execute the next planned work
3. update `Execution > Progress` at every stopping point, splitting partially completed work into done vs. remaining work when needed
4. run validation from `Validation > Test Plan`
5. record unexpected findings in `Execution > Surprises & Discoveries` and decisions in `Execution > Decision Log`
6. record concise proof in `Validation > Test Results`
7. update `Outcomes & Retrospective` at major milestones or completion
8. when the approved spec is satisfied, validation is complete, and retrospective is written, set frontmatter `status` to `completed` and populate `completed_at`
9. continue until the approved spec is satisfied or a real blocker appears

## Output Contract

When the command is allowed to proceed, it should clearly report:

- that the spec is approved
- that the spec is now the source of truth
- that implementation is beginning
- that routine updates must be written back into the spec during execution
- that outcomes and retrospective notes must be written before completion

When the command is blocked, it should clearly report:

- that implementation did not start
- which approval precondition failed
- that the operator must use `/spec-review` first
