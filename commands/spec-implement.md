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

The implementation agent must not stop merely because one milestone is complete, tests are partly passing, context is getting long, or a convenient summary point has been reached. The default behavior is to finish the entire approved spec in the current turn. The only valid reasons to stop before completion are a real blocker, a required user decision, missing credentials/access, an unsafe/destructive action requiring approval, or validation that cannot be completed in the current environment.

The implementation agent should:

1. read the approved spec
2. create or maintain a local task checklist covering every item in the spec's `Plan of Work`, `Acceptance Criteria`, and `Validation > Test Plan`
3. execute planned work until the checklist is complete or a valid blocker appears
4. after each milestone, immediately continue to the next incomplete checklist item instead of asking whether to proceed
5. update `Execution > Progress` at every stopping point, splitting partially completed work into done vs. remaining work when needed
6. run validation from `Validation > Test Plan`
7. treat any failing validation as implementation work: diagnose the failure, fix in-scope bugs, rerun the relevant validation, and repeat until it passes or a valid blocker/out-of-scope defect is proven
8. record unexpected findings in `Execution > Surprises & Discoveries` and decisions in `Execution > Decision Log`
9. record concise proof in `Validation > Test Results`
10. update `Outcomes & Retrospective` at major milestones or completion
11. when the approved spec is satisfied, validation is complete, and retrospective is written, set frontmatter `status` to `completed` and populate `completed_at`
12. continue until the approved spec is satisfied or a valid blocker appears

Failed tests, type checks, linters, build commands, smoke checks, or manual validation steps are not completion points. The implementation agent must inspect the failure, determine whether it is caused by the current spec work, and fix every in-scope bug discovered during validation before stopping.

The agent may stop with failing validation only when the failure is proven to be outside the approved scope, caused by unavailable credentials/services/environment, requires a user decision, or requires an unsafe/destructive action that needs approval. In that case, document the evidence, remaining failures, and exact next action in `Validation > Test Results` and `Execution > Progress`.

## Completion Gate

The command is not complete, and the agent must keep working in the same turn, until all of the following are true:

- every in-scope item from `Plan of Work` is implemented or explicitly marked not done with a blocker
- every `Acceptance Criteria` item is satisfied or explicitly marked blocked with evidence
- every feasible command or check from `Validation > Test Plan` has been run
- every in-scope bug discovered during validation has been fixed and the relevant validation has been rerun
- `Validation > Test Results` records what passed, what failed, what was fixed after failure, and what could not run
- `Execution > Progress` has no ambiguous remaining work
- `Outcomes & Retrospective` states the final outcome, known gaps, and follow-up work
- frontmatter `status` is `completed` and `completed_at` is populated, unless a valid blocker prevents completion

If a valid blocker prevents completion, the final response must identify the blocker, the exact remaining checklist items, and the next action needed from the operator. Do not present a partial milestone as finished.

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
