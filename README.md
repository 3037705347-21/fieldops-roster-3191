# FieldOps Roster

FieldOps Roster is a local command-line tool for small field-service teams that need to schedule work orders without an external dispatch platform. Dispatch coordinators can validate a workspace file, assign open work to technicians, mark jobs complete or cancelled, and produce a regional backlog report.

## Project Structure

- `cmd/rosterctl`: CLI entrypoint and command parsing.
- `internal/model`: domain objects, status transitions, and validation helpers.
- `internal/store`: JSON workspace loading, saving, and integrity checks.
- `internal/planner`: assignment planning by priority, due time, region, skill, and remaining capacity.
- `internal/service`: workflow operations that update orders and assignments.
- `internal/report`: regional operational summaries for open, scheduled, completed, cancelled, and overdue work.
- `testdata`: sample workspace data for local runs.

## Workflows

1. Validate a workspace before dispatching.
2. Assign open work orders to active technicians with matching region, skill, and available capacity.
3. Complete or cancel scheduled orders while preserving valid status transitions.
4. Generate a JSON report grouped by region for backlog and capacity review.

## Commands

Build the CLI:

```bash
go build ./...
```

Run all tests:

```bash
go test ./...
```

Validate sample data:

```bash
go run ./cmd/rosterctl validate -file testdata/workspace.json
```

Plan assignments and update the workspace file:

```bash
go run ./cmd/rosterctl plan -file testdata/workspace.json
```

Complete an order:

```bash
go run ./cmd/rosterctl complete -file testdata/workspace.json -order wo-1001 -at 2026-08-16T15:00:00Z
```

Cancel an order:

```bash
go run ./cmd/rosterctl cancel -file testdata/workspace.json -order wo-1002
```

Print a regional report:

```bash
go run ./cmd/rosterctl report -file testdata/workspace.json -now 2026-08-16T12:00:00Z
```

## Workspace Format

The workspace is a JSON file containing technicians, work orders, existing assignments, a generation timestamp, and a business date. Timestamps use RFC3339 UTC format. The tool uses only local files and has no required environment variables.
