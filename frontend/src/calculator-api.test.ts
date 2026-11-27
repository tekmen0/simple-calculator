import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { calculate } from "./calculator-api";

const fetchMock = vi.fn<typeof fetch>();

beforeEach(() => {
  vi.stubGlobal("fetch", fetchMock);
});

afterEach(() => {
  fetchMock.mockReset();
  vi.unstubAllGlobals();
});

describe("calculate", () => {
  it.each(["add", "subtract", "multiply", "divide"] as const)(
    "posts numeric operands to the %s endpoint",
    async (operation) => {
      fetchMock.mockResolvedValueOnce(Response.json({ result: 0 }));

      await expect(calculate(operation, { left: -1.25, right: 4 })).resolves.toEqual({
        ok: true,
        result: 0,
      });
      expect(fetchMock).toHaveBeenCalledExactlyOnceWith(`/api/calculations/${operation}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: '{"left":-1.25,"right":4}',
      });
    },
  );

  it("returns the code and message of a structured service error", async () => {
    const error = { code: "division_by_zero", message: "Cannot divide by zero." };
    fetchMock.mockResolvedValueOnce(Response.json({ error }, { status: 422 }));

    await expect(calculate("divide", { left: 7, right: 0 })).resolves.toEqual({
      ok: false,
      error,
    });
  });

  it.each([
    ["missing result", "{}", 200],
    ["string result", '{"result":"3.5"}', 200],
    ["nonfinite result", '{"result":1e999}', 200],
    ["null body", "null", 200],
    ["non-JSON proxy failure", "<html>Bad gateway</html>", 502],
    ["result with a failed HTTP status", '{"result":3.5}', 500],
    ["missing error code", '{"error":{"message":"Try again."}}', 500],
    ["missing error message", '{"error":{"code":"internal_error"}}', 500],
    ["empty error message", '{"error":{"code":"internal_error","message":" "}}', 500],
    ["null error", '{"error":null}', 500],
  ])("rejects an unexpected response with %s", async (_, body, status) => {
    fetchMock.mockResolvedValueOnce(new Response(body, { status }));

    await expect(calculate("add", { left: 1, right: 2 })).resolves.toMatchObject({
      ok: false,
      error: { code: "invalid_response" },
    });
  });

  it("reports a rejected fetch as a connection failure", async () => {
    fetchMock.mockRejectedValueOnce(new TypeError("Failed to fetch"));

    await expect(calculate("add", { left: 1, right: 2 })).resolves.toMatchObject({
      ok: false,
      error: { code: "network_error" },
    });
  });
});
