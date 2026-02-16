# MVP Build Plan
## Library Management System (Go + Bubble Tea)

Version: 1.0  
Date: 2026-02-16

---

## 1. Goal

Build a production-quality MVP for a local-first library management TUI that supports:
- Book and copy management
- Member management
- Loan issue/renew/return workflows
- Overdue tracking (without fines)
- Search/filter/reporting
- Reliable local persistence

This plan maps implementation into clear phases with test gates.

---

## 2. MVP Scope and Priorities

### Must-Have (MVP)
1. Books, copies, members CRUD (with archival/deactivation where appropriate)
2. Loan issue, renew, return with business rule enforcement
3. Overdue status and overdue report
4. Keyboard-first TUI with contextual help
5. Persistent local storage and safe writes
6. Unit tests for core business logic

### Should-Have (If time allows)
1. Dashboard aggregates
2. CSV/JSON export
3. CSV import with validation report

### Out of Scope
1. Monetary fines
2. Multi-branch and network synchronization
3. External notification services

---

## 3. Architecture Strategy

### Style
Use a **layered, dependency-inward architecture**:
- `ui` depends on `application`
- `application` depends on `domain` and repository interfaces
- `infrastructure` implements storage/logging/config adapters

No business logic in TUI components; Bubble Tea models orchestrate use cases only.

### Core Design Principles
- Keep domain models and rules deterministic and testable.
- Use interfaces at boundaries (repositories, clock, ID generator).
- Keep Bubble Tea update loop thin by delegating to application services.
- Centralize style tokens and keymaps.

---

## 4. Recommended Project Structure (Best Practice)

```text
lib-management/
  cmd/
    lms/
      main.go

  internal/
    app/
      usecase/
        book_service.go
        member_service.go
        loan_service.go
        report_service.go
      dto/
        book.go
        member.go
        loan.go
      ports/
        repository.go
        tx.go
        logger.go
        clock.go

    domain/
      book/
        entity.go
        rules.go
      copy/
        entity.go
      member/
        entity.go
        rules.go
      loan/
        entity.go
        rules.go
      shared/
        errors.go
        types.go

    infra/
      storage/
        json/
          store.go
          book_repo.go
          member_repo.go
          loan_repo.go
        sqlite/
          db.go
          migrations/
      config/
        config.go
      log/
        logger.go
      time/
        clock.go
      id/
        generator.go

    ui/
      tui/
        model.go
        route.go
        keymap.go
        styles.go
        statusbar.go
      views/
        dashboard/
          model.go
          view.go
        books/
          list.go
          form.go
          detail.go
        members/
          list.go
          form.go
          detail.go
        loans/
          issue.go
          return.go
          renew.go
          list.go
        reports/
          overdue.go
        shared/
          table.go
          confirm.go
          search.go
          help.go

  pkg/
    version/
      version.go

  test/
    integration/
      app_flow_test.go

  docs/
    SRS.md
    MVP_PLAN.md

  .golangci.yml
  Makefile
  go.mod
  go.sum
  README.md
```

### Why this structure
- `cmd` keeps entrypoint thin.
- `internal` protects app internals from external imports.
- Domain remains framework-agnostic and highly testable.
- UI is isolated, allowing future alternative interfaces (CLI/API).

---

## 5. Implementation Phases

## Phase 0 - Foundation (Day 1)

### Deliverables
- Repo initialization and module setup
- Basic app bootstrap and main Bubble Tea loop
- Config loading, logger, graceful shutdown
- Base Makefile targets

### Tasks
1. Initialize Go module and dependency set.
2. Add Bubble Tea/Bubbles/Lip Gloss dependencies.
3. Create `cmd/lms/main.go` with app wiring.
4. Add `make run`, `make test`, `make lint`.
5. Add CI baseline (test + lint).

### Exit Criteria
- App starts with placeholder screen.
- CI pipeline green on empty baseline tests.

---

## Phase 1 - Domain and Use Cases (Days 2-4)

### Deliverables
- Domain entities and business rules
- Application services and interfaces
- Unit tests for rule validation

### Tasks
1. Define entities: Book, Copy, Member, Loan.
2. Implement rules:
   - copy availability checks
   - member eligibility checks
   - due date calculation
   - renewal constraints
   - overdue determination
3. Create service interfaces and DTOs.
4. Add unit tests for all critical rule branches.

### Exit Criteria
- `go test ./...` passes for domain/app packages.
- Business rules are independent of TUI/storage.

---

## Phase 2 - Persistence Layer (Days 5-6)

### Deliverables
- JSON storage adapter (MVP default)
- Repository implementations
- Safe write strategy and load validation

### Tasks
1. Implement repository interfaces in `infra/storage/json`.
2. Use atomic file replace strategy for safe writes.
3. Add schema/version field for future migrations.
4. Add persistence tests (read-write-read consistency).

### Exit Criteria
- Data survives restart simulations.
- Corrupt file handling returns actionable errors.

---

## Phase 3 - TUI Core Framework (Days 7-8)

### Deliverables
- Navigation shell, routes, keymaps
- Shared components (table/search/help/confirm)
- Unified Lip Gloss theme tokens

### Tasks
1. Implement root model with route switching.
2. Define global and context key bindings.
3. Build reusable table wrapper and command palette style search.
4. Add status bar notifications for operations.
5. Ensure terminal resize handling.

### Exit Criteria
- User can move between views with keyboard only.
- Styling is consistent and readable.

---

## Phase 4 - Feature Views (Days 9-12)

### Deliverables
- Books and copies screens
- Members screens
- Loans issue/renew/return screens
- Overdue report screen

### Tasks
1. Books:
   - list/search/sort
   - create/edit/archive
   - copy inventory actions
2. Members:
   - list/search
   - create/edit/deactivate
   - history view
3. Loans:
   - issue flow with validation
   - return flow
   - renewal flow
   - active/overdue filters
4. Reports:
   - overdue list grouped/sorted by due date

### Exit Criteria
- End-to-end MVP workflow is usable manually in TUI.

---

## Phase 5 - Hardening and Quality (Days 13-14)

### Deliverables
- Integration tests for key flows
- Error handling polish
- Performance checks on larger datasets
- Documentation and release notes

### Tasks
1. Integration tests for:
   - add member/book/copy
   - issue/renew/return
   - overdue transition
2. Add robust input validation and error messages.
3. Benchmark key operations against NFR target.
4. Finalize README with keybindings and setup.

### Exit Criteria
- MVP acceptance criteria fully met.
- Test suite and lint are green.

---

## 6. Testing Strategy

### Unit Tests (Highest Priority)
- Domain rules and use case services
- Edge cases: max loans reached, invalid renewals, inactive member

### Integration Tests
- Storage + use case interactions
- End-to-end command sequences without TUI rendering dependency

### Manual TUI Verification Checklist
1. Full keyboard navigation works.
2. Form validation is clear and recoverable.
3. Terminal resize keeps layout usable.
4. Overdue indicators are visible and accurate.

### Tooling
- `go test ./...`
- `go test -race ./...`
- `golangci-lint run`

---

## 7. MVP Definition of Done

MVP is done when all are true:
1. SRS Must requirements are implemented.
2. Tests cover core business flows and pass reliably.
3. No critical data-loss path in normal operation.
4. TUI is keyboard-driven with help and confirmations.
5. Documentation includes setup, keybindings, and architecture notes.

---

## 8. Risks and Mitigation

1. **Risk:** TUI complexity grows quickly.  
   **Mitigation:** Reusable shared components and strict view boundaries.

2. **Risk:** Business logic leaks into UI update handlers.  
   **Mitigation:** Route all mutations through use case services only.

3. **Risk:** Data corruption in JSON storage.  
   **Mitigation:** Atomic writes, backups, and startup validation.

4. **Risk:** Performance degradation with large tables.  
   **Mitigation:** Pagination/virtualization patterns and indexed in-memory lookups.

---

## 9. Suggested Milestones

- **M1 (End Day 4):** Domain + services + unit tests complete
- **M2 (End Day 6):** Persistence stable and tested
- **M3 (End Day 8):** TUI shell and shared components complete
- **M4 (End Day 12):** Feature complete MVP workflows
- **M5 (End Day 14):** Hardened release candidate
