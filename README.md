# LMS-bit

Library management system built with Go and a TUI stack:
- Bubble Tea
- Bubbles
- Lip Gloss

## Current Status

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
make lint
```

## Docs

- `docs/SRS.md`
- `docs/MVP_PLAN.md`
