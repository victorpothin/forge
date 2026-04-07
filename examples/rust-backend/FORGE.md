# FORGE — Example: Rust + Axum Backend

> This is a filled example of the FORGE template applied to a Rust API using Axum.

---

## LAYER 1 — CONTEXT

### Project
A backend service for processing webhook events from a payment provider. Receives events, validates signatures, persists to database, and publishes to an internal message queue.

### Tech Stack
- Rust (stable)
- Axum 0.7 for HTTP
- SQLx with PostgreSQL
- tokio for async runtime
- serde/serde_json for serialization
- cargo test for unit tests

### Current State
- Handler for `POST /webhooks` exists but does not validate the HMAC signature
- Events are inserted into the database but without deduplication
- No tests exist
- Queue publishing is stubbed with a `todo!()`

### Integrations
- PostgreSQL via `DATABASE_URL` env var
- Message queue: RabbitMQ (not yet implemented — stubbed)
- Payment provider sends HMAC-SHA256 signature in `X-Signature` header

### Entry Points
- `cargo run` → starts on port 8080
- `docker compose up` → starts with PostgreSQL and RabbitMQ

---

**[ GATE 1 ]** ✅ — Context validated by user.

---

## LAYER 2 — PROBLEM

### Priority 1
Webhook handler does not validate the HMAC signature, meaning any caller can send fake events and they will be processed.

### Priority 2
No deduplication — the same event can be inserted multiple times if the provider retries.

---

**[ GATE 2 ]** ✅ — P1 is a security issue, P2 is a data integrity issue. Both confirmed.

---

## LAYER 3 — LOCKED PATH

### [DECISION] — Technical decisions already made
```
[DECISION] Use HMAC-SHA256 for signature validation. The key is in env var WEBHOOK_SECRET.
Reason: Required by the payment provider's documentation.
Scope: Webhook handler only.

[DECISION] Deduplication via event_id unique constraint in the database.
Reason: Simplest approach, already have the event_id field.
Scope: Database and insert logic.
```

### [THOUGHT] — Conceptual directions that are off-limits
```
[THOUGHT] Do not implement the RabbitMQ queue in this session.
Reason: Queue infrastructure is not ready. The todo!() stub stays.
Implication: Do not touch the queue module. Do not suggest queue alternatives.
```

---

**[ GATE 3 ]** ✅ — Locked paths confirmed.

---

## LAYER 4 — PLANNING

### Task List

| # | Task | Depends On | Skills |
|---|---|---|---|
| T1 | Implement HMAC-SHA256 signature validation function | — | execution/task-runner |
| T2 | Integrate signature validation into webhook handler | T1 | execution/task-runner, execution/scope-guard |
| T3 | Add unique constraint on event_id and handle duplicate insert error | — | execution/task-runner |
| T4 | Write unit tests for signature validation | T1 | testing/unit |
| T5 | Write unit tests for duplicate handling | T3 | testing/unit |

### Skills Required

| Skill | Status |
|---|---|
| skills/execution/task-runner.md | existing |
| skills/execution/scope-guard.md | existing |
| skills/testing/unit.md | existing |

---

**[ GATE 4 ]** ✅ — Plan approved. T1 and T3 are parallel (no dependency between them).

---

## LAYER 5 — EXECUTION

### T1 — Implement HMAC-SHA256 signature validation function
- Status: `done`
- Output: `validate_signature(payload: &[u8], secret: &str, header: &str) -> bool` in `src/webhook/signature.rs`

### T3 — Add unique constraint on event_id and handle duplicate insert error
- Status: `done`
- Output: Migration `0002_event_id_unique.sql`, insert logic returns `Ok(InsertResult::Duplicate)` on conflict

### T2 — Integrate signature validation into webhook handler
- Status: `done`
- Output: Handler returns `401 Unauthorized` when signature is missing or invalid

### T4 — Unit tests for signature validation
- Status: `done`
- Output: 3 tests in `src/webhook/signature_tests.rs`, all passing

### T5 — Unit tests for duplicate handling
- Status: `done`
- Output: 2 tests in `src/webhook/insert_tests.rs`, all passing

---

## LAYER 6 — TESTING

### T1 / T4 Tests
- [x] valid signature → returns true
- [x] invalid signature → returns false
- [x] missing header → returns false

### T3 / T5 Tests
- [x] first insert → Ok(Inserted)
- [x] duplicate insert → Ok(Duplicate), no error

---

## LAYER 7 — DOCUMENTATION

### What Was Built
1. HMAC-SHA256 webhook signature validation. The handler now rejects any request without a valid `X-Signature` header.
2. Event deduplication via unique constraint on `event_id`. Duplicate events are silently accepted (idempotent) rather than causing an error.

### How to Use
Set the `WEBHOOK_SECRET` environment variable to your webhook signing secret. The handler will validate all incoming requests automatically.

```
WEBHOOK_SECRET=your_secret_here cargo run
```

Duplicate events: safe to retry. The endpoint returns `200 OK` for both new and duplicate events.

### Dependencies
- `hmac` and `sha2` crates added to `Cargo.toml`
- Migration `0002_event_id_unique.sql` must be applied before deployment

### Known Limitations
- Signature timestamp validation not implemented (replay attack window is open)
- Queue publishing still stubbed — events are not forwarded downstream
