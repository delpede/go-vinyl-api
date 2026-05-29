# Go Vinyl API — Spec Alignment TODO

Goal: make the implementation, `openspec/specs/record-api/spec.md`, and
`api/openapi.yaml` all agree. Spec-driven; spec.md is the source of truth.

## Done
- [x] **Search (`q`) over `title` + `artist`** — `GetRecords` in `cmd/api/main.go`
      now honors `?q=` with case-insensitive substring match (`ILIKE %q%`),
      parameterized. Empty/absent `q` returns all records. Removed stray
      `fmt.Println` debug line.
      NOTE: compiles + vets, but NOT yet run against a live DB.

## Next up (in order)

### 1. Metadata columns — biggest item, unblocks the rest
Six optional fields exist in `openapi.yaml` + generated structs
(`internal/api/generated/api.gen.go`) but are persisted nowhere:
`label, format, releaseYear, country, condition, notes`.
- [ ] New migration adding columns: `label, format, release_year, country,
      condition, notes` (mirror `migrations/20260515081402_create_records.up.sql`,
      add matching `.down.sql`).
- [ ] `PostRecords` INSERT (`cmd/api/main.go:~70`) — add the 6 columns +
      `RETURNING` them.
- [ ] `PutRecordsId` UPDATE (`cmd/api/main.go:~135`) — same.
- [ ] All SELECTs (list + get-by-id) and the `Scan(...)` calls — add the 6
      columns. They're pointers (`*string`/`*int`) in the generated `Record`.
- [ ] Finish search: extend the `WHERE` in `GetRecords` to also match
      `label` and `notes` (marked with a `TODO` in the code).

### 2. Required-field validation — DECISION NEEDED
Spec/openapi mark `title` + `artist` required, but generated structs have no
`binding:"required"` tags, so `POST {}` currently inserts blank strings and
returns 201 instead of 400.
- [ ] Decide: enforce required fields (add validation / binding tags) or
      accept current behavior? Spec has no explicit "missing required field"
      scenario, so it's a judgment call.

### 3. Tidy `openapi.yaml` (docs-only, three-way alignment)
- [ ] Document `400 Bad Request` (malformed body) on POST/PUT.
- [ ] Document `500` on DB failure.
- [ ] Define an `Error` schema (`{ "error": string }`) and reference it from
      the 404/400/500 responses (currently 404s have no `content`).

### 4. Tests (optional but recommended)
- [ ] No tests exist. Add coverage for spec scenarios: search hit/miss/empty,
      404s, 201/200/204.

## How to verify locally
- `docker-compose up` for Postgres (see `docker-compose.yaml`), apply migrations,
  then run `go run ./cmd/api` and exercise endpoints (`:8080`).
- Or ask Claude to run the `/verify` flow end-to-end.

## Reference
- Spec: `openspec/specs/record-api/spec.md`
- OpenAPI: `api/openapi.yaml`
- Handlers: `cmd/api/main.go`
- Generated (do not edit): `internal/api/generated/api.gen.go`
