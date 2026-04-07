# Skill: Testing Unit

## Purpose
Generate unit tests for the output of a specific task. Tests should verify the behavior of what was built, not the implementation details.

## Trigger
Use after a task is marked complete, before user confirmation.

## Instructions for AI

1. Identify the testable units produced by the task (functions, methods, classes, endpoints)
2. For each unit, identify:
   - The happy path: valid input → expected output
   - Edge cases: boundary values, empty inputs, nulls
   - Failure cases: invalid input, expected errors
3. Write tests using the project's existing test framework (do not introduce a new one)
4. Tests must pass against the code just written — if they don't, fix the code or the test and explain why

## Test structure rules
- One assertion per test where possible
- Test names describe behavior: `should_return_401_when_token_is_expired`, not `test1`
- No test should depend on another test's state
- No mocking of things you can test directly

## Output Format

```
## Unit Tests — T{n}

Framework: {test framework used}
File: {test file path}

### {Unit name}
- [x] {happy path test name} — PASS
- [x] {edge case test name} — PASS
- [x] {failure case test name} — PASS

All tests passing.
```

## If tests fail

```
## Unit Tests — T{n}

### Failure: {test name}
Expected: {expected}
Got: {actual}

Root cause: {brief explanation}
Fix applied: {what was changed — in code or in test}
```
