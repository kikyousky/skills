# skills

## Spec-First Workflow

This repository now includes a spec-first workflow for creating a self-contained executable spec before implementation starts.

### When to use it

Use this workflow for:

- new applications
- new features
- bug fixes
- refactors
- any change that affects multiple files, behavior, contracts, configuration, validation, or dependencies

### Skill model

The workflow is driven by the skill, one reusable template, one optional example, and thin wrapper commands.

- `spec-authoring/assets/template.md` defines the exact spec shape
- `spec-authoring/SKILL.md` defines how to author the spec
- `spec-authoring/references/example-spec.md` shows a fully populated approved spec

### Repository contents

- `spec-authoring/assets/template.md` — canonical spec template
- `spec-authoring/SKILL.md` — spec-authoring skill
- `spec-authoring/references/example-spec.md` — filled example spec
- `commands/spec-draft.md` — draft command
- `commands/spec-review.md` — review and approval gate
- `commands/spec-start.md` — implementation handoff command

### Command sequence

#### 1. `/spec-draft`

Start from a raw request. Inspect the target repository, derive a slug, and create a spec draft at:

`docs/specs/{YYYY-MM-DD}-{slug}.md`

The result must stay in `draft` or `reviewed` status.

#### 2. `/spec-review`

Review a drafted spec under `docs/specs/`. This is the only supported approval gate. It either:

- rejects the spec with concrete remediation notes, or
- updates frontmatter to `status: approved` and sets `approved_at`

#### 3. `/spec-start`

Start implementation only from an approved spec. This command treats the approved spec as the source of truth and hands execution into implementation mode.

### Target-project storage rule

Write specs into the target project at:

`docs/specs/{YYYY-MM-DD}-{slug}.md`

Do not store working specs elsewhere. For this workflow, the supported spec location is `docs/specs/`.

### Approval rule

Implementation may begin only when:

- the file is under `docs/specs/`
- frontmatter `status` is `approved`
- frontmatter `approved_at` is populated

`reviewed` is an optional pre-approval state for a spec that has been refined and is ready for approval review, but it is still not executable.

### Source-of-truth contract

Once approved, the spec becomes the single source of truth for implementation. The implementation agent must follow the approved scope and must not silently expand it.

### Living-document update requirement

During implementation, routine updates go to:

- `Execution > Progress`
- `Execution > Decision Log`
- `Execution > Concrete Steps`
- `Validation > Test Results`

All other sections are change-controlled and should be edited only when implementation facts actually change.

### Example

See `spec-authoring/references/example-spec.md` for a fully populated approved example using the final template.
