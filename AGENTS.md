# Calculator project conventions

## Go backend

- Use `gofmt` and the standard library for HTTP, JSON, arithmetic, and tests. Keep startup and server configuration in `main`.
- Keep arithmetic and domain errors in the pure calculator package. Let the HTTP adapter validate request shape, own wire types, and translate domain errors to HTTP responses.
- Return JSON errors with a stable `error.code` and a human-readable `error.message`.

## React frontend

- Use strict TypeScript, React, plain CSS, and browser `fetch`. Keep form state local.
- Validate input for immediate feedback and delegate arithmetic to the Go API. Keep request construction and response validation in the API client. Narrow decoded JSON from `unknown`.
- Use labeled, keyboard-operable controls, accessible status messages, and layouts that work at narrow viewport widths.

## Shared conventions

- Keep Go tests under `backend/tests/` in subdirectories named for the packages they test. Keep frontend tests under `frontend/tests/`. Use table-driven Go tests for arithmetic, `httptest` for HTTP behavior, and Vitest with React Testing Library for form behavior. Assert results and errors through public interfaces.
- Keep the Go handler, TypeScript client, and README API examples in sync when the wire contract changes.
- After each change, run the tests and checks defined in [CI](.github/workflows/ci.yml).

## Submitting changes

- Push a feature branch and open a pull request against the intended base branch. Never push directly to `main`.
- Keep the diff scoped to the request. In the pull request, explain what changed, why, and how you verified it. Wait for CI to pass and address review feedback before merging.

## Reviewing changes

- Check the diff against the requested behavior and these conventions. Prioritize bugs, regressions, and missing tests for changed behavior. Flag unnecessary abstractions or unrelated edits when they make the change harder to maintain.
- Report actionable findings in severity order. Cite the file and line, explain the impact, and mark optional suggestions as such.
