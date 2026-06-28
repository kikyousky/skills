---
name: golang-backend
description: Use when designing, implementing, reviewing, or refactoring Go backend code, especially HTTP/gRPC/queue boundaries, service/domain/repo package structure, request validation, JSON decoding, config access, naming, or complexity constraints. Provides definitive architecture rules plus scripts for enforcing those constraints.
---

# Go Backend Rules

Strict, mechanically enforceable rules for building Go backends with explicit boundaries.

## Core Principles

1. Treat all external input as untrusted at the boundary.
2. Allow only validated internal input types into `internal/service`.
3. Keep cross-layer access explicit and restricted.
4. Treat naming as part of the contract.
5. Prefer a small rigid rule set over many soft conventions.

## Package Layout

```text
internal/
  transport/
    http/
    grpc/
    queue/
  service/
  domain/      # optional for simple CRUD
  repo/
  provider/
  config/
  platform/
```

Package roles:

- `internal/domain`: optional business entities, value objects, domain rules, and domain errors that do not depend on transport, persistence, config, or external SDKs. Skip this package for simple CRUD when the types are mostly data shapes and the rules are light enough to live in `internal/service`.
- `internal/repo`: persistence adapters and data access implementations. Repositories translate between storage records and domain/service types, and must not contain transport concerns.
- `internal/platform`: low-level process and infrastructure plumbing such as clocks, IDs, logging adapters, filesystem helpers, environment adapters, and other code that supports the app but is not business logic.

## Writing Rules

### Boundary Input

- All external input must use a struct named `*Request`.
- Request structs must live under `internal/transport/...`.
- Request structs are boundary-only types and must not be defined in `internal/service`, `internal/domain`, or `internal/repo`.

### Request Validation

- Every request struct must implement `ValidateAndSanitize() (service.XInput, error)`.
- The returned value must be a type from `internal/service` whose name ends with `Input`.
- Raw request structs must never cross the transport-to-service boundary.

Required flow:

```text
Decode -> ValidateAndSanitize -> service.Input -> service method
```

### Service Input

- All inputs entering `internal/service` must be named `*Input`.
- Service methods may accept only `context.Context`, `service.*Input`, and service-owned interfaces.
- Service code must not accept transport request types.

### JSON Decoding

- HTTP JSON decoding must use `json.Decoder`.
- The decoder must call `DisallowUnknownFields()` before `Decode()`.
- Do not use `json.Unmarshal` for HTTP request decoding.

Required pattern:

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
if err := dec.Decode(&req); err != nil {
    return err
}
```

### Package Isolation

- `internal/domain` must not import `internal/transport` or `internal/service`.
- `internal/service` must not import `internal/transport` or transport frameworks.
- `internal/repo` must not import `internal/transport` or `internal/service`.
- External SDKs should enter business code through `internal/provider` or `internal/repo`.

### Config

- Raw `os.Getenv` is allowed only in `internal/config`.
- Config structs must validate before use.

### Naming

Required suffixes:

- transport input: `*Request`
- transport output: `*Response`
- service input: `*Input`
- repository implementation: `*Repository`
- provider interfaces: `*Provider`

Forbidden suffixes:

- `Manager`
- `Helper`
- `Util`
- `Data`

### Complexity

- Maximum file size: 400 lines, excluding generated files, migrations, and test fixtures.
- Maximum non-constructor function size: 60 lines.
- Maximum non-constructor function parameter count: 4, excluding `context.Context`.

### Logging

- Use structured logging through a logger abstraction.
- Use stable event names instead of prose log messages.
- Include the error object as a structured field on error logs.

## Scripts

This skill ships runnable enforcement scripts under `golang-backend/scripts`.

Execution mode: source only. Assume the user is developing a Go backend and should have a working Go toolchain.

### What It Enforces Definitively

- request types must be named `*Request`
- request types must live under `internal/transport/...`
- transport request types must implement `ValidateAndSanitize() (service.XInput, error)`
- `internal/service`, `internal/domain`, and `internal/repo` import boundaries
- service function parameter restrictions
- raw request structs passed directly into service calls from transport code
- HTTP JSON decoder usage with `DisallowUnknownFields()`
- raw `os.Getenv` outside `internal/config`
- file length, function length, and function parameter count
- banned suffixes such as `Manager`, `Helper`, `Util`, and `Data`

### What It Checks Heuristically

- direct third-party SDK usage in `internal/service`
- structured logging patterns
- provider-only cross-cutting access
- typed bad-request errors

Use these as review rules even when the scripts cannot prove them perfectly.

## Preflight

Check the local Go toolchain before using the scripts:

```bash
go version
```

If `go version` fails, fix the local Go installation first. This skill expects a working source-mode Go environment.

## Run The Enforcement Script

```bash
cd golang-backend/scripts
go version
go run ./cmd/golang-backend-lint ./...
```

## Example Diagnostics

```text
request type `CreateUserRequest` is declared outside `internal/transport`.
Move boundary request structs under `internal/transport/...`.
```

```text
request `CreateUserRequest` does not implement `ValidateAndSanitize() (service.CreateUserInput, error)`.
Boundary request structs must validate and normalize input before entering `internal/service`.
```

```text
service call receives request `CreateUserRequest` directly.
Call `ValidateAndSanitize()` first and pass the resulting `service.*Input` value instead.
```

```text
raw `os.Getenv` usage is forbidden outside `internal/config`.
Load env vars in `internal/config`, validate them, and pass typed config inward.
```
