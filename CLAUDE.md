# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Test
go test ./...

# Test single package
go test ./jdb/...

# Test single test
go test ./jdb/... -run TestName

# Add dependency
go get github.com/cgalvisleon/et@v1.0.19

# Versioning (bump patch/minor/major, tags and pushes)
./version.sh --r          # patch
./version.sh --n          # minor
./version.sh --m          # major
```

## Architecture

**`jql`** is a Go ORM/query-builder library targeting PostgreSQL (and potentially SQLite). The public API is in `v1/` (package `jql`); the core engine is in `jdb/`.

### Packages

- **`jdb/`** — Core engine. All database primitives live here.
- **`v1/`** — Thin public wrapper over `jdb`. Re-exports types and adds convenience constructors (`ConnectTo`, `Load`, `LoadTo`, `Define`, `From`).
- **`drivers/postgres/`** — PostgreSQL driver using `lib/pq`. Registers itself via `init()` calling `jdb.Register("postgres", ...)`.
- **`drivers/josefina/`** — Secondary driver (registers similarly).
- **`tenant/`** — Multi-tenant helper. Maps tenant IDs to `*jdb.DB` instances.

### Core Concepts

**Driver pattern** (`jdb/driver.go`): Drivers implement the `Driver` interface (`Connect`, `Load`, `Mutate`, `Query`, `Command`) and self-register in their `init()`. A driver must be imported (blank import) to be available.

**DB / Model registry**: Global maps in `jdb/jdb.go` hold loaded `*DB` and `*Model` instances, keyed by normalized name. `LoadDb` creates and connects; `GetDb` retrieves or rehydrates from the catalog.

**Catalog** (`jdb/catalog.go`): A `core.catalog` table in the database persists serialized `DB` and `Model` definitions. This enables model rehydration across restarts without re-calling `Define*` methods. The catalog itself is a `Model` (IsCore=true, skips catalog persistence).

**Model definition** (`jdb/define.go`): Models are built by calling `Define*` methods:
- `DefineColumn` — regular SQL column
- `DefineAttribute` — attribute stored inside a JSON `source` column
- `DefineDetail` — 1-to-many child model (auto-joined on query)
- `DefineRollup` — many-to-1 lookup (auto-joined on query)
- `DefineCalc` — computed field via a `DataContext` callback
- `DefineModel()` — shortcut adding `created_at`, `updated_at`, `status`, `id` (PK), and `source` fields

**Query builder** (`jdb/ql.go`): `NewQuery(model, alias)` returns a `*Ql`. Fluent API:
```go
items, err := jdb.NewQuery(model, "A").
    Select("name", "email").
    Where(jdb.Eq("status", "active")).
    OrderBy("name").
    Limit(1, 20)
```
After `All()`/`Limit()`, details/rollups/calcs are resolved concurrently with goroutines.

**Command builder** (`jdb/cmd.go`): `model.Insert(data)`, `model.Update(data)`, `model.Delete()`, `model.Upsert(data)` return a `*Cmd`. Supports `Where`, trigger hooks (`BeforeInsert`, `AfterUpdate`, etc.), and transaction (`ExecTx`).

**Trigger hooks**: Both `Model` and `Cmd` support `Before/After Insert/Update/Delete` hooks of type `TriggerFunction func(tx *Tx, old, new et.Json) error`. Returning `ErrNotUpdated` from a BeforeUpdate hook is a no-op signal (skips the update without error propagation).

**Two storage modes**:
- **COLUMN**: standard SQL columns — `DefineColumn(name, TypeData, default)`
- **ATTRIB**: stored in a JSON blob field (`source`) — `DefineAttribute(name, TypeData, default)`. Activating at least one attribute auto-creates the `source` column and index.

### Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `DB_DRIVER` | `postgres` | Driver name |
| `DB_NAME` | `josephine` | Database name |
| `DB_HOST` | `localhost` | Host |
| `DB_PORT` | `5432` | Port |
| `DB_USERNAME` | `test` | Username |
| `DB_PASSWORD` | `test` | Password |
| `DB_APP` | `jql` | App name for connection |
| `DEBUG` | `false` | Enable debug SQL logging |
| `MAX_ROWS` | `100` | Default row limit per query |
