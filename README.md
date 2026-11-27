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

### Run with Docker

From the repository root, build and start both services with Docker Compose.

```sh
docker compose up --build
```

Open `http://localhost:8088`. Nginx serves the React build and forwards `/api` to the Go container. To stop the stack, press Ctrl+C or run `docker compose down`.

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

- **One repository, two projects.** Go and React have separate build and test commands. Keeping them in one repository lets a change to the JSON contract update the handler, client, and examples in the same commit.
- **Four operations share one contract.** Addition, subtraction, multiplication, and division all fit the same two-input form and JSON request. The assignment treats exponentiation, square root, and percentage as optional. These would add UI and validation cases, including a single-operand form and request for square root.
- **Stateless arithmetic.** The pure `Calculate` function in `backend/internal/calculator` owns dispatch and arithmetic errors. Each request supplies both operands. The browser displays the latest equation, and a feature like history synced across devices would introduce persistent storage and user identity.
- **Standard-library routing, validation, and logging.** `backend/internal/httpapi` uses `net/http` and `encoding/json` for the health and calculation routes, request validation, and JSON errors. One handler serves all four operations. Request logging uses a small `slog` middleware because it applies across routes. A routing framework or schema library would add a dependency for a small, fixed API.
- **Local form state and native controls.** One React form owns the operands, selected operation, and submission state. A global store or form library would add wiring for state used in one place. Native controls provide keyboard interaction, and plain CSS handles the responsive layout.
- **The frontend API client builds requests.** `calculate()` sends the shared JSON POST with browser `fetch` and validates the response at runtime. Client generation would add a schema and build step for one frontend consumer. A shared `fetch` wrapper would make sense if several frontend calls needed the same authentication headers or timeout policy.
- **Numbers use binary floating point.** Go `float64` and JavaScript `number` support fractions with ordinary JSON numbers. Results can round, as in `0.1 + 0.2 = 0.30000000000000004`. Exact decimal arithmetic would require an explicit precision and rounding policy in both projects.
- **Same-origin serving.** Vite proxies `/api` to Go over local HTTP. In production, the host manages TLS and routes `/api` on the frontend's HTTPS origin. Both setups use the same relative URLs. Separate origins would add CORS configuration, and local HTTPS would add certificate setup.
- **Direct local setup.** Go and npm launch the two processes using the project's existing toolchains. Docker Compose builds and starts both services with one command.
- **Built-in code checks.** Strict TypeScript and behavior tests cover the single-form frontend. The backend uses `go vet` for static analysis and `gofmt` for formatting. If the frontend grows to use effects or custom hooks, ESLint with React's Hooks rules can catch conditional hook calls and missing dependencies. [Rules of Hooks](https://react.dev/reference/eslint-plugin-react-hooks/lints/rules-of-hooks) · [Dependency checks](https://react.dev/reference/eslint-plugin-react-hooks/lints/exhaustive-deps)

## Testing and coverage

Go tests live under `backend/tests/`, grouped by the package they exercise. React and TypeScript tests live under `frontend/tests/`.

Table-driven Go tests cover arithmetic and edge cases. `httptest` checks request validation, status codes, and response bodies through the real handler. Vitest and React Testing Library cover calculation submissions, input validation, pending state, and service errors. The API client tests cover request formatting and response validation.

Each automated suite runs independently. A browser smoke check exercises the running Vite and Go processes together. A dedicated browser suite would add browser setup and server lifecycle management for the single calculation workflow.

### Coverage summary

| Project | Metric | Coverage |
| --- | --- | ---: |
| Go backend | Statements | 91.1% |
| React frontend | Lines | 98.33% |
| React frontend | Branches | 94.91% |

From `backend/`, run the Go tests, print coverage, run static analysis, and list any files that need formatting. The `-coverpkg` flag measures coverage in the implementation packages exercised by the separate test packages.

```sh
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go vet ./...
gofmt -l .
```

The `total` line from `go tool cover` shows the combined coverage across the Go test suites.

Run the frontend tests, generate coverage, and check types from `frontend/`.

```sh
npm test
npm run coverage
npm exec -- tsc --noEmit
```

Open `frontend/coverage/index.html` for the HTML report.

## AI prompts

See [USED-PROMPTS.md](USED-PROMPTS.md) for the prompts used on the assignment.
