# simple-calculator

A calculator with a React and TypeScript frontend and a Go REST API. 

## Local setup

Install Go 1.26.5 or later and Node.js 24 LTS (24.15 or later) with npm.

```sh
npm --prefix frontend ci
```

The backend uses the Go standard library.

Start the Go API on port `8080` in one terminal. Air rebuilds and restarts the API when you save a Go application source file.

```sh
cd backend
go run github.com/air-verse/air@v1.67.4
```

For a one-off run, use `go run .` from `backend/`.

From the repository root in a second terminal, start Vite.

```sh
cd frontend
npm run dev
```

Open `http://localhost:5173`. Vite forwards requests under `/api` to the Go server so the browser can use relative URLs during development.

### Production build

Run `npm run build` from `frontend/` to type-check the project and build `frontend/dist/`. Serve those files over HTTPS and route `/api` on the same origin to the Go backend. The browser's relative `fetch` URLs use HTTPS when the page is served over HTTPS.

## API usage

### Health check

`GET /api/health` returns `200 OK` with `{"status":"ok"}` when the API is running. `HEAD /api/health` returns the same status without a response body.

```sh
curl -i http://localhost:8080/api/health
```

### Calculations

Send `POST /api/calculations/{operation}` with `Content-Type: application/json`. The supported operations are `add`, `subtract`, `multiply`, and `divide`. Use finite JSON numbers for `left` (the first operand) and `right` (the second operand).

```sh
curl -i http://localhost:8080/api/calculations/divide \
  -H 'Content-Type: application/json' \
  -d '{"left":7,"right":2}'
```

The response is `200 OK` with `{"result":3.5}`.

```sh
curl -i http://localhost:8080/api/calculations/divide \
  -H 'Content-Type: application/json' \
  -d '{"left":7,"right":0}'
```

The response has status `422` with a stable error code and a human-readable message.

```json
{"error":{"code":"division_by_zero","message":"Cannot divide by zero."}}
```

| Status | Error code | Cause |
| --- | --- | --- |
| `400` | `invalid_request` | Invalid JSON or fields, including missing operands, `null`, numeric strings, numbers that overflow `float64`, unknown or duplicate fields, or trailing JSON values |
| `404` | `not_found` | Unknown API path or operation |
| `405` | `method_not_allowed` | A supported operation was called with another method. The response includes `Allow: POST`. |
| `413` | `request_too_large` | The request body exceeds 4 KiB |
| `415` | `unsupported_media_type` | The request lacks an `application/json` content type |
| `422` | `division_by_zero`, `result_out_of_range` | Division by either sign of zero or arithmetic overflow |

## Design decisions

- **One repository, two projects.** This was asked by the Task deliverables, not our choice. This allow for simpler management & ci/cd operations. But bad for CI/CD, increasing the complexity of management of versions, and blast radius.
- **Arithmetic has one entry point.** A pure `Calculate` function in `backend/internal/calculator` owns dispatch and arithmetic errors. This allows shared interface for tests. 
`backend/internal/httpapi` owns JSON decoding and translates domain errors to HTTP responses. The four operation URLs share one handler. Each request is stateless.
- **React owns the interaction.** A single form keeps input and submission state local. It validates input for immediate feedback and displays the server result. A small TypeScript API client uses simple browser `fetch` and validates responses at runtime. Plain CSS supports narrow screens.
- **Numbers use binary floating point.** Go uses `float64` and the browser uses JavaScript `number`. This keeps the JSON contract simple and supports fractional inputs. Values such as `0.1` are approximate, so `0.1 + 0.2` can return `0.30000000000000004`.
- **KISS** This project didn't use pnpm, Gin or other extra libraries to keep the set-up simple and don't overkill the problem at hand.

## Testing and coverage

Table-driven Go tests cover arithmetic and edge cases. `httptest` checks request validation, status codes, and response bodies through the real handler. Vitest and React Testing Library cover calculation submissions, input validation, pending state, and service errors. The API client tests cover request formatting and response validation.

Run the Go tests and print the coverage summary from `backend/`.

```sh
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Run the frontend tests and generate the coverage report from `frontend/`.

```sh
npm test
npm run coverage
```

Open `frontend/coverage/index.html` for the HTML report.

## AI prompts

See [USED-PROMPTS.md](USED-PROMPTS.md) for the prompts used on the assignment.
