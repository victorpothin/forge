# Skill: Testing Acceptance

## Purpose
Verify that the completed task actually solves the problem it was created for. Acceptance tests validate the problem statement, not just the code.

## Trigger
Use after unit tests pass, before marking a task as fully done.

## Instructions for AI

1. Go back to the Problem layer and find the problem this task addresses
2. Derive acceptance criteria from the problem description:
   - "Given [context], when [action], then [expected outcome]"
3. Verify each criterion against the implemented output
4. If any criterion is not met, the task is NOT done — report and fix

## Acceptance criteria format

```
Given {context}
When {action}
Then {outcome}
```

## Example

Problem: "Users cannot reset their password"

```
Given a user with a valid account
When they request a password reset
Then they receive an email with a reset link

Given a user who clicks a valid reset link
When they submit a new password
Then their password is updated and the link is invalidated

Given a user who clicks an expired reset link
When they try to submit a new password
Then they see an error and are prompted to request a new link
```

## Output Format

```
## Acceptance Tests — T{n}
Problem addressed: P{n}

- [x] Given ... When ... Then ... — PASS
- [x] Given ... When ... Then ... — PASS
- [ ] Given ... When ... Then ... — FAIL

Problem P{n}: [SOLVED / PARTIALLY SOLVED / NOT SOLVED]
```

## Rules
- If a problem has multiple acceptance criteria, all must pass for the problem to be considered solved
- "Partially solved" is valid — document what remains and create a follow-up task
- Do not invent acceptance criteria the user did not ask for
