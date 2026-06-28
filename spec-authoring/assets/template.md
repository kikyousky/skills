---
title: <short action-oriented title, e.g. Add project archiving>
slug: <kebab-case-slug, e.g. add-project-archiving>
status: draft
request_type: <new-app|feature|bugfix|refactor|other>
created_at: <YYYY-MM-DDTHH:MM:SSZ>
approved_at: null
completed_at: null
---

# <Short action-oriented title>

This spec is a living document. During implementation, `Execution > Progress`, `Execution > Decision Log`, `Execution > Concrete Steps`, `Validation > Test Results`, and `Outcomes & Retrospective` are the routine-update areas. All other sections are change-controlled and should be updated only when implementation facts actually change.

Use this file as both the required structure and an example of expected specificity. Replace the example content with target-project facts while keeping the section order, level of detail, and integrated “what to write / why it matters” guidance.

## Purpose / Big Picture

Explain the work from the user or operator perspective. State what becomes possible, how to observe success, what is in scope, and what is intentionally out of scope. This prevents implementation from drifting into adjacent work that was never approved.

## Implementation Impact

Document every expected code, data, contract, config, and dependency impact before approval. If a category is unaffected, write `none`. GUI impact may be high-level; all other categories must be concrete enough for review and implementation.

### Database / Schema

State the exact persistence impact because schema changes are durable and hard to undo. Name tables, columns, indexes, constraints, migrations, seeds, and data backfills. When creating or altering tables, include the actual SQL that will run, not only ORM model changes or prose descriptions.

Example:

Update `prisma/schema.prisma` so the `Project` model adds these fields:

- `archivedAt DateTime? @map("archived_at")`
- `archivedById String? @map("archived_by_id")`
- `archivedBy User? @relation("ProjectArchivedBy", fields: [archivedById], references: [id], onDelete: SetNull)`

Add the inverse relation to `User`:

- `archivedProjects Project[] @relation("ProjectArchivedBy")`

Create a Prisma migration under `prisma/migrations/{timestamp}_add_project_archiving/` that runs equivalent SQL:

```sql
ALTER TABLE projects ADD COLUMN archived_at TIMESTAMPTZ NULL;
ALTER TABLE projects ADD COLUMN archived_by_id TEXT NULL;
ALTER TABLE projects ADD CONSTRAINT projects_archived_by_id_fkey FOREIGN KEY (archived_by_id) REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX projects_active_workspace_idx ON projects(workspace_id, updated_at DESC) WHERE archived_at IS NULL;
CREATE INDEX projects_archived_workspace_idx ON projects(workspace_id, archived_at DESC) WHERE archived_at IS NOT NULL;
```

No data backfill is required because existing projects remain active with `archived_at = NULL`. No existing rows should be deleted or rewritten.

### GUI / UX

Describe the visible behavior users will see. Name screens, flows, loading states, disabled states, errors, and labels when known. This can stay higher-level than schema or API details but must still make the user-visible outcome reviewable.

Example:

Update `src/client/pages/ProjectsPage.tsx` so the default list shows only active projects. Add an `Archived` filter option with values `Active` and `Archived` that maps to `GET /api/projects?archived=false` and `GET /api/projects?archived=true`.

Update `src/client/pages/ProjectDetailPage.tsx` so archived projects show an `Archived` badge, disable the `New task` button, disable inline project-name editing, and show helper text: `Archived projects are read-only until an admin restores them.` Admins should see `Archive project` for active projects and `Restore project` for archived projects in the project settings menu.

### API / Contracts

Specify exact endpoint, payload, event, job, service-contract, or exported-interface changes. Contract details affect callers, tests, backward compatibility, and implementation boundaries.

Example:

Update `src/server/routes/projects.ts` and `src/server/services/projectService.ts` with these contract changes:

- `GET /api/projects` accepts optional query `archived=true|false|all`; default is `false`.
- `GET /api/projects?archived=false` returns only rows where `archivedAt` is `null`.
- `GET /api/projects?archived=true` returns only rows where `archivedAt` is not `null`.
- `GET /api/projects?archived=all` returns both active and archived projects and requires the `admin` role.
- `GET /api/projects/:projectId` continues to return archived projects for users who already have project access.
- `POST /api/projects/:projectId/archive` requires `admin`, sets `archivedAt` to the server timestamp, sets `archivedById` to the current user id, and returns `200` with the updated project JSON.
- `POST /api/projects/:projectId/unarchive` requires `admin`, sets `archivedAt` and `archivedById` to `null`, and returns `200` with the updated project JSON.
- `POST /api/projects/:projectId/tasks` returns `409` with `{ "error": "PROJECT_ARCHIVED" }` when `archivedAt` is not `null`.

Extend the project JSON shape returned by list and detail endpoints:

```json
{
  "id": "proj_123",
  "name": "Project Alpha",
  "workspaceId": "ws_123",
  "archivedAt": "2026-06-14T12:00:00.000Z",
  "archivedById": "user_123"
}
```

For active projects, `archivedAt` and `archivedById` must be `null`.

### Config / Environment

State whether the work changes environment variables, feature flags, secrets, deployment configuration, build settings, or third-party service configuration. This protects local setup, deployment, and rollback safety.

Example:

No new environment variables, feature flags, secrets, deployment settings, or third-party service configuration are required.

The implementation uses the existing PostgreSQL connection, existing Prisma migration workflow, and existing role middleware. Deployment only needs the normal migration step already used by this repository: `npm run prisma:migrate:deploy`.

### Dependencies

State whether packages, SDKs, services, or infrastructure dependencies are added, removed, or upgraded. Name exact dependencies and explain why each is needed because dependency changes affect install, build, security, and maintenance risk.

Example:

No packages, SDKs, services, or infrastructure dependencies are added, removed, or upgraded.

The existing dependencies are sufficient: Prisma handles schema and migration changes, Express handles the new routes, and the current React component library provides the badge, menu, and disabled-button styles.

## Execution

Track the implementation path, decisions, and commands in the same source-of-truth document that approved the work.

### Progress

Use a list with checkboxes to summarize granular steps. Every stopping point must be documented here, even if it requires splitting a partially completed task into two (“done” vs. “remaining”). This section must always reflect the actual current state of the work.
Use timestamps to measure rates of progress. Add task-specific checklist items when useful.

Example:
- [x] (2026-06-14 10:15Z) Initial repo inspection completed.
- [x] (2026-06-14 10:40Z) Spec drafted with scope, impact, and validation.
- [ ] Spec approved.
- [ ] Validation completed.
- [ ] Outcomes and retrospective completed.
- [ ] Example partially completed step (completed: API routes; remaining: browser QA).

### Surprises & Discoveries

Document unexpected behaviors, bugs, optimizations, or insights discovered during implementation. Provide concise evidence.

- Observation: …
  Evidence: …

### Decision Log

Record every decision made while working on the plan in the format:

- Decision: …
  Rationale: …
  Date/Author: …

### Plan of Work

Describe, in prose, the sequence of edits and additions. For each edit, name the file and location (function, module) and what to insert or change. Keep it concrete and minimal.

### Concrete Steps

State the exact commands to run and where to run them (working directory). When a command generates output, show a short expected transcript so the reader can compare. This section must be updated as work proceeds.

## Validation

Define how success will be proven before implementation starts, then record what actually happened.

### Test Plan

Specify required automated and agent-executed checks before work begins so completion is objective. State what to test, how to test it, and whether each check is a unit test, API integration test, E2E test, migration/schema check, typecheck, lint, or manual/agent QA.

Example:

- Unit test: `npm test -- tests/unit/projectArchivePolicy.test.ts`
  Verify `canEditProject`, `canCreateTask`, and `canRestoreProject` return the expected decisions for active projects, archived projects, admins, and non-admins.
- API integration test: `npm test -- tests/api/projects.test.ts`
  Verify `GET /api/projects` defaults to active projects, `archived=true` returns archived projects, archive/unarchive endpoints require admin, and creating a task in an archived project returns `409 PROJECT_ARCHIVED`.
- Migration/schema check: `npm run prisma:migrate:test`
  Verify the migration creates `archived_at`, `archived_by_id`, the foreign key, and the active/archived workspace indexes using the expected SQL.
- E2E test: `npm run test:e2e -- projects-archive.spec.ts`
  Verify an admin can archive a project from settings, the project moves to the archived filter, archived detail pages are read-only, and restore makes task creation available again.
- Typecheck/lint: `npm run typecheck` and `npm run lint`
  Verify the new response fields, route handlers, and UI code satisfy repository quality gates.
- Agent QA: run the app locally and exercise the archive/restore flow in a browser.
  Capture the tested browser path, user role, and any screenshots or trace links needed to prove the visible behavior.

### Test Results

For a draft spec, write `Not run yet; implementation has not started.` After validation, replace that with concise proof so completion does not depend on external memory or chat logs.

Example:

  Command: npm test -- tests/api/projects.test.ts
  Result: PASS - 18 project API tests passed, including active default filtering, archived filter, admin-only archive/unarchive, direct archived detail read, and `409 PROJECT_ARCHIVED` task creation failure.

  Command: npm run typecheck
  Result: PASS - no TypeScript errors for project response fields or UI prop changes.

## Outcomes & Retrospective

For a draft spec, write `Not started yet; implementation outcomes are not available.` During implementation or at completion, replace that with outcomes, gaps, and lessons learned. Compare the actual result against the original `Purpose / Big Picture` so the spec closes the loop.
