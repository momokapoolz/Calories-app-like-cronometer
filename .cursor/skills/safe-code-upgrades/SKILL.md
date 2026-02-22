---
name: safe-code-upgrades
description: Upgrades, fixes bugs in, or refactors existing code without breaking other functionality. Performs impact analysis, applies incremental safe changes with rollback awareness, prevents regressions, profiles performance changes, and follows a reproduce→isolate→fix→verify bug cycle. Use when modifying existing systems, fixing production bugs, refactoring, improving performance, or any change that must preserve current behavior.
---

# Safe Code Upgrades

## Core Principle

**Never change what you haven't understood.** Always read before touching. Always verify after changing.

---

## Phase 1 — Impact Analysis (Before Any Change)

1. **Read the target code fully.** Understand its current behavior, not just the surface.
2. **Map dependencies:**
   - What calls this code? (callers / consumers)
   - What does this code call? (dependencies)
   - What data does it read/write? (state, DB, files, cache)
3. **Identify contracts:** function signatures, return types, error shapes, side effects, and any behavior other code relies on.
4. **Flag blast radius:** list every file/function that could be affected by the change.
5. **Check for hidden coupling:** global state, shared caches, event listeners, middleware chains.

> If the blast radius is large, break the change into smaller independent pieces.

---

## Phase 2 — Safe Refactoring Steps

Apply changes **incrementally** — one logical unit at a time.

1. **Establish a baseline** before changing anything:
   - Note current behavior (outputs, logs, metrics, test results).
   - If tests exist, run them and record the green state.
2. **Make the smallest possible change** that moves toward the goal.
3. **Preserve the public contract** unless explicitly changing it:
   - Same function signatures (or add overloads/wrappers).
   - Same error types and shapes.
   - Same side effects unless the bug is *in* the side effect.
4. **Use the Strangler Fig pattern for large refactors:**
   - Keep the old code path alive.
   - Route a subset of calls to the new path.
   - Verify parity, then fully migrate, then delete old code.
5. **Commit/checkpoint after each green state.** Never accumulate multiple unverified changes.

---

## Phase 3 — Bug Isolation Cycle

```
Reproduce → Isolate → Fix → Verify → Prevent recurrence
```

| Step | Action |
|------|--------|
| **Reproduce** | Write the smallest possible reproducer (test, script, curl command). Do not guess — confirm the bug is real and consistent. |
| **Isolate** | Binary-search the call chain. Narrow to the exact function/line causing incorrect behavior. |
| **Root cause** | Ask "why" at least twice. Fix the root, not the symptom. |
| **Fix** | Apply the minimal change. Avoid touching unrelated code in the same diff. |
| **Verify** | Re-run the reproducer — confirm it passes. Run all existing tests. |
| **Prevent** | Add a regression test that would have caught this bug. |

---

## Phase 4 — Regression Prevention

Before submitting any change:

- [ ] All pre-existing tests still pass.
- [ ] A new test covers the specific behavior changed or the bug fixed.
- [ ] Edge cases are tested: empty input, nil/null, zero values, max values, concurrent access (if applicable).
- [ ] Error paths are tested, not just the happy path.
- [ ] If no test framework exists, manually document the verification steps taken.

**Do not rely on "it looks right."** If you can't write a test, at minimum write a comment explaining the expected invariant.

---

## Phase 5 — Performance Upgrades

Only optimize after confirming a real bottleneck:

1. **Profile first** — never guess. Use language-appropriate profiling tools.
2. **Record the baseline metric** (latency, memory, CPU, query count).
3. **Change one thing at a time** and measure after each change.
4. **Preserve behavior** — a faster function with wrong output is worse than a slow correct one.
5. **Document the trade-off** — if performance was gained at the cost of readability or memory, say so.

---

## Checklist Before Finishing Any Change

- [ ] Did I read all affected code before changing it?
- [ ] Does the public contract remain intact (or is the change intentional and communicated)?
- [ ] Do all existing tests pass?
- [ ] Is there a test or verifiable proof for the new behavior?
- [ ] Is the diff minimal — no unrelated changes bundled in?
- [ ] Would someone reading this diff understand *why* the change was made?

---

## Anti-Patterns to Avoid

| Anti-Pattern | What to Do Instead |
|---|---|
| "While I'm in here…" cleanup bundled with a bug fix | Separate commits/PRs |
| Fixing the symptom (e.g., catching a panic) | Find and fix the root cause |
| Deleting old code before the new code is proven | Keep old path until new path is verified |
| Skipping tests because "it's a small change" | Small changes cause the most silent regressions |
| Changing function behavior to fix a caller | Fix the caller; keep the function contract stable |
