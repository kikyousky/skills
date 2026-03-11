---
name: e2e-testing-setup
description: E2E testing with Playwright in isolated Docker Compose environment, read this when setting up or running e2e tests for web apps
---

# E2E Testing with Playwright and Docker

This guide covers setting up and running E2E tests in an isolated Docker environment.

## Core Principles

1. **Complete Isolation**: No ports exposed to host, only internal Docker network
2. **Production builds for E2E tests** - Never mount source code for E2E tests. Use built images for reliability.
3. **Healthchecks over startup order** - `depends_on` alone is insufficient. Always use healthchecks.
4. **Cache at build time, not runtime** - Use BuildKit cache mounts to speed up builds.
5. **Service names for networking** - Inside Docker network, services communicate via service names, not `localhost`.
6. **Don't run `docker-compose down` between test runs** - Keep services running to save time

---

## Docker Compose Structure

### Minimal Working Setup

```yaml
version: "3.8"

services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    healthcheck:
      test: pg_isready -U postgres
      interval: 2s
      timeout: 5s
      retries: 10
    # Fast ephemeral storage for tests
    tmpfs:
      - /var/lib/postgresql/data
    # Optional: disable fsync for speed (test only, never in prod)
    command:
      - postgres
      - -c
      - fsync=off
      - -c
      - synchronous_commit=off

  api:
    build: ./api
    environment:
      # Use service name "db", not "localhost"
      DATABASE_URL: postgres://postgres:postgres@db:5432/app
    depends_on:
      db:
        condition: service_healthy
    healthcheck:
      test: curl -f http://localhost:3000/health || exit 1
      interval: 5s
      timeout: 5s
      retries: 10

  frontend:
    build: ./frontend
    environment:
      # Internal Docker network URL
      API_URL: http://api:3000
    depends_on:
      api:
        condition: service_healthy
    healthcheck:
      test: curl -f http://localhost:3000 || exit 1
      interval: 5s
      timeout: 5s
      retries: 10

  e2e:
    image: mcr.microsoft.com/playwright:v1.58.2-jammy
    depends_on:
      frontend:
        condition: service_healthy
    environment:
      BASE_URL: http://frontend:3000
      PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD: 1
    volumes:
      - ./e2e:/app
      - ./playwright-report:/app/playwright-report
      - ./test-results:/app/test-results
    working_dir: /app
    command: npx playwright test

  # Debug container for troubleshooting network issues
  debug:
    image: curlimages/curl
    profiles: ["debug"]
    command: sleep infinity

volumes:
  postgres_data:
```

---

## Next.js Dockerfile (Production Build with Caching)

```dockerfile
# syntax=docker/dockerfile:1
FROM node:20-alpine AS deps

WORKDIR /app

COPY package*.json ./

# Cache npm packages across builds
RUN --mount=type=cache,target=/root/.npm \
    npm ci

FROM node:20-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Cache .next/cache for faster rebuilds
# Contains: webpack cache, SWC cache, fetch-cache, optimized images
RUN --mount=type=cache,target=/app/.next/cache \
    npm run build

FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Copy only necessary files for standalone output
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

CMD ["node", "server.js"]
```

### Required next.config.js Setting

```javascript
module.exports = {
  output: 'standalone',
}
```

---

## Healthcheck Patterns

### Aggressive but Safe Healthcheck Timing

```yaml
healthcheck:
  test: curl -f http://localhost:3000/health || exit 1
  interval: 1s       # check frequently
  timeout: 2s        # fail fast
  retries: 30        # allow enough attempts
  start_period: 5s   # grace period before counting failures
```

### Common Healthcheck Commands

| Service | Healthcheck |
|---------|-------------|
| PostgreSQL | `pg_isready -U postgres` |
| HTTP service | `curl -f http://localhost:PORT/health \|\| exit 1` |
| Redis | `redis-cli ping` |
| MySQL | `mysqladmin ping -h localhost` |

---

## Debugging

When services can't communicate, use the debug container:

```bash
# Start debug container
docker-compose --profile debug up -d

# Get shell inside the Docker network
docker-compose exec debug sh

# Test connectivity from inside the network
curl http://api:3000/health
curl http://frontend:3000
nslookup api
```

---

## Performance Optimization

### 1. Enable BuildKit

```bash
export DOCKER_BUILDKIT=1
```

Or in Docker Compose:

```yaml
# docker-compose.yml (top level)
version: "3.8"
```

Then run with:
```bash
DOCKER_BUILDKIT=1 docker-compose build
```

### 2. Keep Containers Running Between Test Runs

```bash
# Start once
docker-compose up -d --wait

# Run tests multiple times (fast - no startup)
docker-compose exec e2e npx playwright test
docker-compose exec e2e npx playwright test

# Tear down when done
docker-compose down
```

### 3. Parallel Service Startup

Structure dependencies for maximum parallelism:

```
db ────────┐
           ├──► api ──► frontend ──► e2e
redis ─────┘
```

Don't create unnecessary serial dependencies.

### 4. Use Alpine Images

```yaml
# Slower
image: node:20

# Faster
image: node:20-alpine
```

---

## BuildKit Cache Mounts

Cache mounts persist data between builds without including it in the final image.

### Basic Syntax

```dockerfile
RUN --mount=type=cache,target=<path>,id=<id>,sharing=<mode> \
    <command>
```

| Option | Purpose | Default |
|--------|---------|---------|
| `target` | Path inside container to mount cache | Required |
| `id` | Unique identifier for this cache | Target path |
| `sharing` | `shared`, `private`, or `locked` | `shared` |

### Common Cache Locations

| Package Manager | Cache Path |
|-----------------|------------|
| npm | `/root/.npm` |
| yarn | `/usr/local/share/.cache/yarn` |
| pnpm | `/root/.local/share/pnpm/store` |
| pip | `/root/.cache/pip` |
| go | `/go/pkg/mod` |
| Next.js build | `/app/.next/cache` |

### Cache Isolation by ID

When working with multiple projects on the same machine, use unique IDs to prevent cache collisions:

```dockerfile
# Project A
RUN --mount=type=cache,target=/root/.npm,id=project-a-npm npm ci
RUN --mount=type=cache,target=/app/.next/cache,id=project-a-next npm run build

# Project B
RUN --mount=type=cache,target=/root/.npm,id=project-b-npm npm ci
RUN --mount=type=cache,target=/app/.next/cache,id=project-b-next npm run build
```

Without unique IDs, projects share the same cache by default (based on target path). This can cause:
- Cache thrashing when projects have different dependencies
- Unexpected cache invalidation
- Subtle build inconsistencies

### Important Requirements

1. **Enable BuildKit**:
   ```bash
   export DOCKER_BUILDKIT=1
   docker-compose build
   ```

2. **First line must be syntax directive**:
   ```dockerfile
   # syntax=docker/dockerfile:1
   FROM node:20-alpine
   ```

3. **Cache the package manager's cache, not node_modules**:
   ```dockerfile
   # Bad: overwrites node_modules with empty cache on first run
   RUN --mount=type=cache,target=/app/node_modules npm ci

   # Good: cache npm's download cache
   RUN --mount=type=cache,target=/root/.npm npm ci
   ```

### Managing Cache

```bash
# List build cache
docker builder prune --dry-run

# Clear build cache
docker builder prune

# Clear everything including cache mounts
docker builder prune -a
```

### Debugging Cache

Verifying that caching is working can be difficult. Use these techniques:

#### 1. Time Comparison

```bash
# Cold cache
docker builder prune -f
time DOCKER_BUILDKIT=1 docker build -t myapp .

# Warm cache (should be much faster)
time DOCKER_BUILDKIT=1 docker build -t myapp .
```

#### 2. Inspect Cache Storage

```bash
# List all BuildKit cache entries with sizes
docker builder du --verbose

# Output shows cache IDs and sizes:
# ID                    RECLAIMABLE   SIZE      LAST ACCESSED
# myproject-npm         true          245MB     5 minutes ago
# myproject-gomod       true          180MB     5 minutes ago
```

#### 3. Go Module Cache Paths

Go uses two separate caches:

| Path | Contents |
|------|----------|
| `/go/pkg/mod` | Downloaded module source code |
| `/root/.cache/go-build` | Compiled packages |

Both should be cached:
```dockerfile
RUN --mount=type=cache,target=/go/pkg/mod,id=myproject-gomod \
    --mount=type=cache,target=/root/.cache/go-build,id=myproject-gobuild \
    go build -o /app/server .
```

---

## Workflow Commands

### Full Test Run

```bash
# Build images (uses cache)
docker-compose build

# Start all services and wait for health
docker-compose up -d --wait

# Run E2E tests
docker-compose run --rm e2e

# View results
open playwright-report/index.html
```

### Quick Iteration (Containers Already Running)

```bash
# Rebuild only changed service
docker-compose build frontend

# Restart that service
docker-compose up -d frontend

# Run tests
docker-compose exec e2e npx playwright test
```

### Force Full Rebuild

```bash
docker-compose build --no-cache
```

### Clean Everything

```bash
docker-compose down -v --rmi local
```

### Benchmark Startup Time

```bash
time docker-compose up -d --wait
```

---

## Troubleshooting Checklist

When E2E tests fail to start:

1. Check all services are healthy: `docker-compose ps`
2. Check service logs: `docker-compose logs <service>`
3. Verify network connectivity from debug container
4. Ensure environment variables use service names, not localhost
5. Verify healthcheck endpoints exist and return 200

When tests fail intermittently:

1. Check if it's a timing issue - add proper `waitFor` assertions
2. Verify test data is reset between runs
3. Check for port conflicts on host machine
4. Review resource constraints (memory, CPU)

When builds are slow:

1. Verify BuildKit is enabled
2. Check Dockerfile layer ordering (dependencies before source)
3. Use cache mounts for node_modules and .next/cache
4. Consider using smaller base images (alpine)
