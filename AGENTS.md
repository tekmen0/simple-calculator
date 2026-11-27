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

- Co-locate tests with their code. Use table-driven Go tests for arithmetic, `httptest` for HTTP behavior, and Vitest with React Testing Library for form behavior. Assert results and errors through public interfaces.
- Keep the Go handler, TypeScript client, and README API examples in sync when the wire contract changes.
