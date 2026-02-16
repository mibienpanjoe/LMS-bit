# Software Requirements Specification (SRS)
## Library Management System (TUI)

Version: 1.0  
Date: 2026-02-16  
Status: Draft

---

## 1. Introduction

### 1.1 Purpose
This SRS defines the requirements for a **Library Management System** built in **Go** as a **Terminal User Interface (TUI)** using:
- Bubble Tea (application state and update loop)
- Bubbles (reusable TUI components)
- Lip Gloss (consistent styling/theme)

The document is intended for developers, testers, and maintainers.

### 1.2 Scope
The system provides a keyboard-first workflow for librarians to:
- Manage books and copies
- Manage member accounts
- Issue, renew, and return books
- Track due dates and overdue items
- Search and filter books/members/loans
- Persist data across sessions

MVP targets a single-library, local-first application.

### 1.3 Definitions
- **TUI**: Terminal User Interface
- **LMS**: Library Management System
- **CRUD**: Create, Read, Update, Delete
- **MVP**: Minimum Viable Product
- **ISBN**: International Standard Book Number

### 1.4 Document Conventions
- Requirement IDs:
  - `FR-*` Functional
  - `NFR-*` Non-functional
  - `UI-*` User interface
  - `DATA-*` Data constraints
- Priority: **Must**, **Should**, **Could**

---

## 2. Overall Description

### 2.1 Product Perspective
Standalone cross-platform terminal app with clean layering:
- Domain models and business rules
- Application services/use cases
- Storage adapters (JSON or SQLite)
- Bubble Tea presentation layer

### 2.2 Product Functions
- Book catalog and copy inventory management
- Member registration and lifecycle management
- Loan issuance, renewal, and return
- Overdue tracking (no monetary penalties)
- Search/filter/sort across records
- Basic dashboard and reports

### 2.3 User Class
1. **Librarian (Primary)**: full operational workflows.

### 2.4 Operating Environment
- Linux/macOS/Windows terminals
- ANSI color capable terminal
- Recommended minimum terminal size: 100x30

### 2.5 Constraints
- Must be implemented in Go.
- Must use Bubble Tea + Bubbles + Lip Gloss.
- All major actions must be keyboard driven.

### 2.6 Assumptions
- Single local installation for MVP.
- Data stored locally with regular backups by operator.
- System clock/timezone is valid for due-date logic.

---

## 3. Specific Requirements

### 3.1 Functional Requirements

#### 3.1.1 Book Catalog
- **FR-010 (Must):** Create book records with title, authors, ISBN, category, publisher, year.
- **FR-011 (Must):** Edit book metadata.
- **FR-012 (Must):** Archive/deactivate books from circulation.
- **FR-013 (Must):** Validate ISBN format where provided.
- **FR-014 (Must):** Support multiple physical copies per book title.

#### 3.1.2 Copy Inventory
- **FR-020 (Must):** Create/update copy records with unique copy ID/barcode.
- **FR-021 (Must):** Track copy status: Available, Loaned, Reserved (optional), Damaged, Lost.
- **FR-022 (Should):** Track copy condition notes.

#### 3.1.3 Member Management
- **FR-030 (Must):** Register members with unique member ID.
- **FR-031 (Must):** Update member details (name, contact, status).
- **FR-032 (Must):** Deactivate/reactivate members.
- **FR-033 (Should):** View member borrowing history.

#### 3.1.4 Loans and Returns
- **FR-040 (Must):** Issue available copy to active member.
- **FR-041 (Must):** Prevent issue when copy unavailable.
- **FR-042 (Must):** Prevent issue when member inactive/blocked.
- **FR-043 (Must):** Set due date from configurable loan policy.
- **FR-044 (Must):** Record returns and restore availability.
- **FR-045 (Should):** Support renewals with max renewal count.
- **FR-046 (Must):** Mark loans as overdue when past due date.

#### 3.1.5 Search and Navigation
- **FR-050 (Must):** Search books by title, author, ISBN, category.
- **FR-051 (Must):** Search members by name/member ID/contact.
- **FR-052 (Must):** Filter loans by Active, Overdue, Returned.
- **FR-053 (Should):** Sort list views by key fields.

#### 3.1.6 Dashboard and Reports
- **FR-060 (Should):** Dashboard with total books, copies available, active loans, overdue loans, active members.
- **FR-061 (Should):** Overdue report grouped by member and due date.
- **FR-062 (Could):** Most borrowed books report.

#### 3.1.7 Persistence and Data Ops
- **FR-070 (Must):** Persist all records across application restarts.
- **FR-071 (Should):** Export books/members/loans to CSV or JSON.
- **FR-072 (Should):** Import CSV with row-level validation errors.
- **FR-073 (Could):** Full backup/restore command.

#### 3.1.8 UX Behavior
- **FR-080 (Must):** All core workflows accessible via keyboard shortcuts.
- **FR-081 (Must):** Contextual help view with key bindings.
- **FR-082 (Must):** Confirmation prompt for destructive operations.
- **FR-083 (Must):** Inline validation and actionable error feedback.

### 3.2 Interface Requirements

#### 3.2.1 TUI Layout
- **UI-001 (Must):** Header (app/title/context), main content area, footer/status/help area.
- **UI-002 (Must):** Consistent styling and spacing via centralized Lip Gloss style tokens.
- **UI-003 (Must):** Responsive behavior for reduced terminal size with graceful fallback.
- **UI-004 (Should):** Views for Dashboard, Books, Members, Loans, Reports, Settings.

#### 3.2.2 Software Interfaces
- **SI-001 (Must):** Storage adapter abstraction for JSON or SQLite backend.
- **SI-002 (Should):** Domain/service interfaces independent of TUI package.

### 3.3 Data Requirements

#### 3.3.1 Core Entities
1. **Book**: book_id, title, authors[], isbn, category, publisher, year, status
2. **Copy**: copy_id, book_id, barcode, status, condition_note
3. **Member**: member_id, name, email, phone, joined_at, status
4. **Loan**: loan_id, copy_id, member_id, issued_at, due_at, returned_at, renewal_count, status
5. **Config**: loan_days_default, max_loans_per_member, max_renewals, timezone

#### 3.3.2 Data Integrity Rules
- **DATA-001 (Must):** `member_id` and `copy_id` must be unique.
- **DATA-002 (Must):** Only Available copies can be issued.
- **DATA-003 (Must):** Returned loans cannot be returned again.
- **DATA-004 (Must):** Inactive/blocked members cannot borrow.
- **DATA-005 (Must):** Due date must be after issue date.
- **DATA-006 (Should):** ISBN uniqueness configurable.

### 3.4 Non-Functional Requirements

#### 3.4.1 Performance
- **NFR-001 (Must):** Common operations complete under 200ms at ~10k records locally.
- **NFR-002 (Should):** Startup under 2 seconds in standard environment.

#### 3.4.2 Reliability
- **NFR-010 (Must):** Safe writes to avoid data corruption (atomic replace or transactions).
- **NFR-011 (Must):** Graceful error handling for malformed data and failed operations.

#### 3.4.3 Security
- **NFR-020 (Should):** Optional local auth mode (if enabled) with hashed passwords.
- **NFR-021 (Must):** Input validation for all persisted entities.

#### 3.4.4 Maintainability
- **NFR-030 (Must):** Clear package boundaries (domain/app/storage/ui).
- **NFR-031 (Must):** Unit tests for business-critical rules.
- **NFR-032 (Should):** Interface-driven adapters for replaceable storage.

#### 3.4.5 Usability
- **NFR-040 (Must):** Discoverable keybindings and predictable navigation.
- **NFR-041 (Must):** Consistent status and error messaging.

#### 3.4.6 Portability
- **NFR-050 (Must):** Runs on Linux/macOS/Windows terminals.

---

## 4. Business Rules

- **BR-001:** A member may hold up to `max_loans_per_member` active loans.
- **BR-002:** Default due date = `issued_at + loan_days_default`.
- **BR-003:** Renewal allowed only while loan is active and renewal count < `max_renewals`.
- **BR-004:** Reference-only books are non-circulating and cannot be issued.
- **BR-005:** Overdue status is computed from current date and `due_at` when `returned_at` is null.

---

## 5. Acceptance Criteria (MVP)

The MVP is accepted when:
1. Librarian can create/update/search books, copies, and members.
2. Librarian can issue, renew, and return loans with due-date enforcement.
3. Overdue loans are clearly tracked and reportable.
4. Data persists correctly across restarts.
5. All core workflows are keyboard-driven with help visible in-app.
6. Critical business rules are covered by automated tests.
7. TUI styling is consistent and readable across common terminal sizes.

---

## 6. Out of Scope for MVP

- Monetary fines and payment tracking
- Multi-branch support
- Networked multi-user synchronization
- Email/SMS reminders
- Barcode hardware integrations
