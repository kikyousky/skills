---
title: Add email and password sign-in
slug: add-email-password-sign-in
status: approved
request_type: feature
created_at: 2026-06-09T15:00:00Z
approved_at: 2026-06-09T18:20:00Z
source_of_truth: true
implementation_mode: spec-first
---

# Add email and password sign-in

This spec is a living document. During implementation, `Execution > Progress`, `Execution > Decision Log`, `Execution > Concrete Steps`, and `Validation > Test Results` are the routine-update areas. All other sections are change-controlled and should be updated only when implementation facts actually change.

## Purpose / Big Picture

Users can currently browse the application without an account, but they cannot create an account or return later with a password-based sign-in flow. After this change, a new user can register with email and password, sign in from a login screen, reach the authenticated area of the product, sign out, and sign back in using the same credentials.

The change is successful when a user can complete the full register → authenticated session → sign out → sign in flow and when invalid credentials clearly fail without exposing whether an email address exists.

This change covers first-party email and password registration, login, logout, session creation, protected-route enforcement, and the tests needed to prove those paths work. It does not add social login, password reset, email verification, account settings, profile editing, or multi-factor authentication.

## Context and Orientation

The target repository is assumed to be a web application with a browser-facing UI, a backend API layer, and a relational database. The relevant files are likely to include a user persistence layer, authentication endpoints, UI routes or pages for sign-in and registration, route protection middleware, and automated test suites.

Key repository-relative paths for this example are:

- `db/schema.sql` for user schema changes
- `src/lib/session.ts` for session creation and teardown
- `src/lib/db.ts` for user reads and writes
- `src/app/login/page.tsx` or `src/pages/login.tsx` for sign-in UI
- `src/app/register/page.tsx` or `src/pages/register.tsx` for registration UI
- `src/app/api/auth/` or `src/api/auth/` for auth endpoints
- `src/middleware.ts` for protected-route enforcement
- `tests/integration/` for API validation
- `tests/e2e/` for end-to-end browser validation

In this example, “session” means the authenticated server-recognized state for the current user. “Password hash” means the one-way stored representation of a password, never the raw password itself.

## Implementation Impact

### Database / Schema

Add persistent password storage support for registered users. Update `db/schema.sql` so the user table contains a password hash column. If the repository uses migrations, add a migration file that introduces the new column and ensures newly registered users always receive a stored hash. If a not-null migration is risky for existing rows, use a safe migration sequence and document the transitional state in the migration comments.

### GUI / UX

Add or update registration and login screens. Users should see email and password form fields, validation errors, loading states, and a successful redirect into the authenticated area after sign-in or registration.

### API / Contracts

Add or update authentication endpoints for registration, login, and logout. Registration must accept email and password input, create a user, create a session, and return success. Login must accept email and password input, validate credentials, create a session, and return a generic invalid-credentials failure when authentication fails. Logout must clear the session.

### Config / Environment

No new environment variables are required if the repository already has database and session configuration. If password hashing cost is configurable in the target project, document the exact environment variable or config field used. Otherwise write `none` in the real spec.

### Dependencies

Add one password hashing dependency if the repository does not already have one, for example `bcrypt` or `bcryptjs`. Do not add a second auth framework unless the repository has no viable session mechanism. Reuse the existing session system if it already exists.

## Execution

### Progress

- [x] Initial repo inspection completed.
- [x] Spec drafted with scope, impact, and validation.
- [x] Spec reviewed.
- [x] Spec approved.
- [x] Implementation started from approved spec.
- [x] Validation completed.

### Decision Log

- Decision: Reuse the repository's existing session mechanism instead of introducing JWTs.
  Rationale: Reusing the existing session path reduces risk and keeps route protection aligned with current middleware patterns.
  Date/Author: 2026-06-09 / Atlas example

- Decision: Require explicit implementation-impact coverage for schema, API/contracts, config, and dependencies.
  Rationale: Authentication work is security-sensitive and benefits from very concrete scope definition before coding begins.
  Date/Author: 2026-06-09 / Atlas example

### Plan of Work

First update the user persistence layer so password hashes can be stored safely. Then implement backend registration, login, and logout flows using the repository’s existing session helper. After backend behavior exists, add or update the login and registration screens and ensure successful authentication redirects into the protected area. Finally, enforce route protection for authenticated pages and add both integration and end-to-end validation.

The implementation should touch the schema definition, the user persistence layer, the auth endpoint layer, the session helper if needed, the login and registration UI routes, route protection middleware, and automated tests.

### Concrete Steps

  Working directory: repository root
  Command: ls db src/lib src/app src/pages src/middleware.ts tests/integration tests/e2e
  Expected: confirm the schema location, auth endpoints, session helper, UI route structure, middleware path, and validation directories

  Working directory: repository root
  Command: npm test -- auth
  Expected: run the auth-focused integration coverage before applying the new sign-in flow

  Working directory: repository root
  Command: npm test -- auth --runInBand
  Expected: verify schema, registration, login, logout, duplicate-email, and invalid-credential backend paths pass after the auth changes

  Working directory: repository root
  Command: npx playwright test tests/e2e/auth.spec.ts
  Expected: verify the browser flow for register, authenticated redirect, logout, login, and protected-route enforcement

## Validation

### Test Plan

Run all relevant automated and agent-executed checks:

- schema or migration validation command, if the repository has one
- backend unit or integration tests for registration, login, logout, duplicate-email handling, and invalid-credential handling
- end-to-end browser test for register → authenticated → logout → login
- protected-route check verifying unauthenticated access redirects or fails appropriately

Expected outcomes:

- a new user can register and immediately reach the authenticated area
- the same user can sign out and sign back in
- duplicate-email registration fails clearly
- wrong-password sign-in fails with a generic error
- protected routes reject unauthenticated access
- password hashes, not raw passwords, are stored

### Test Results

Validation completed against the implemented auth flow.

  Command: npm test -- auth
  Result: PASS - 12 auth integration tests passed, covering registration, login, logout, duplicate email, invalid password, and protected-route access.

  Command: npx playwright test tests/e2e/auth.spec.ts
  Result: PASS - 4 browser scenarios passed, covering register, authenticated redirect, logout, login, and logged-out redirect to `/login`.

  Evidence: `POST /api/auth/register` returned `201`, `POST /api/auth/login` returned `200`, and `POST /api/auth/logout` returned `204` for the seeded test user flow.
  Evidence: `GET /dashboard` while logged out redirected to `/login`.
  Evidence: test database row for `casey@example.com` stored a bcrypt hash beginning with `$2b$` and did not store the raw password.
