# LMS-bit

Library management system built with Go and a TUI stack:
- Bubble Tea
- Bubbles
- Lip Gloss

## Current Status

Phase 5 hardening and quality complete:
- Added integration tests for issue/renew/return flow and overdue persistence after reopen
- Added race detector and benchmark workflow via Make targets
- Added `BenchmarkLoanServiceListOverdue10k` baseline performance test
- Improved startup seed error handling with warning logs
- Improved ID generation fallback for safer uniqueness on entropy failure

Phase 4 feature views complete:
- Books view with add, archive, and copy creation workflows
- Members view with registration and status toggle
- Loans view with issue, renew, return, and status filter
- Reports view showing overdue loans
- Dashboard and settings views now populated from persisted data

Phase 3 TUI core framework complete:
- Route-based shell (Dashboard, Books, Members, Loans, Reports, Settings)
- Global and contextual keymaps with help integration
- Shared table and search interaction patterns
- Status bar notifications and confirmation modal scaffold
- Resize-aware layout with minimum terminal fallback

Phase 2 persistence layer complete:
- JSON storage adapter with schema versioning
- Atomic temp-file write and replace for safer persistence
- Repository implementations for books, copies, members, and loans
- Persistence tests for read-write-read consistency and corrupt-file handling

Phase 1 domain and use case foundation complete:
- Domain entities: books, copies, members, loans
- Core lending rules: issue eligibility, renew constraints, return handling, overdue checks
- Application service layer with ports/repositories and DTOs
- Unit tests for critical loan rules and loan service behavior

Phase 0 remains in place:
- Go module and app entrypoint
- Basic Bubble Tea shell with key bindings
- Config and logging scaffolding
- Makefile and CI baseline
- Initial project docs

## Run

```bash
make tidy
make run
```

## Quality Checks

```bash
make test
make test-race
make bench
make lint
```

## Docs

- `docs/SRS.md`
- `docs/MVP_PLAN.md`
- `docs/PHASE5_RELEASE_NOTES.md`
