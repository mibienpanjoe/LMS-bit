# Phase 5 Release Notes

Date: 2026-02-17

## Highlights

- Added end-to-end integration tests for core MVP workflows.
- Added performance benchmark coverage for overdue loan listing.
- Added race-detector validation as a standard quality gate.
- Improved startup resilience with explicit seed error reporting.
- Improved ID fallback behavior to avoid repeated IDs on RNG failure.

## Changes Included

### Testing and Verification

- Added `test/integration/app_flow_test.go`:
  - `TestIssueRenewReturnFlow`
  - `TestOverdueAndPersistenceAfterReopen`
- Added `internal/app/usecase/loan_service_bench_test.go`:
  - `BenchmarkLoanServiceListOverdue10k`
- Updated `Makefile` with `bench` target.

### Reliability and Error Handling

- `cmd/lms/main.go`
  - `seedInitialData` now returns an error.
  - Startup now logs warning on seed failure instead of silently ignoring failures.
- `internal/infra/id/generator.go`
  - Fallback ID now uses a timestamp-based unique string when RNG is unavailable.

## Validation Results

- `go test ./...` passed.
- `go test -race ./...` passed.
- `go build ./...` passed.

Benchmark baseline (local run):

- `BenchmarkLoanServiceListOverdue10k-8`
  - `3399714 ns/op`
  - `4258310 B/op`
  - `18 allocs/op`

## Notes

- Benchmark is currently focused on domain/usecase computation path for 10k in-memory loans.
- Storage write throughput benchmarking can be added in a later phase if needed.
