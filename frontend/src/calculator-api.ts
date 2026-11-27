export type Operation = "add" | "subtract" | "multiply" | "divide";

type Operands = {
  left: number;
  right: number;
};

type CalculationResponse =
  | { ok: true; result: number }
  | { ok: false; error: { code: string; message: string } };

const invalidResponse = {
  ok: false,
  error: {
    code: "invalid_response",
    message: "The calculator service returned an unexpected response. Try again.",
  },
} satisfies CalculationResponse;

export async function calculate(
  operation: Operation,
  operands: Operands,
): Promise<CalculationResponse> {
  let response: Response;
  try {
    response = await fetch(`/api/calculations/${operation}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(operands),
    });
  } catch {
    return {
      ok: false,
      error: {
        code: "network_error",
        message: "Could not reach the calculator service. Try again.",
      },
    };
  }

  let body: unknown;
  try {
    body = await response.json();
  } catch {
    return invalidResponse;
  }

  if (typeof body !== "object" || body === null) {
    return invalidResponse;
  }

  if (response.ok) {
    if ("result" in body && typeof body.result === "number" && Number.isFinite(body.result)) {
      return { ok: true, result: body.result };
    }
    return invalidResponse;
  }

  if ("error" in body && isApiError(body.error)) {
    return { ok: false, error: body.error };
  }
  return invalidResponse;
}

function isApiError(error: unknown): error is { code: string; message: string } {
  return (
    typeof error === "object" &&
    error !== null &&
    "code" in error &&
    typeof error.code === "string" &&
    error.code.trim() !== "" &&
    "message" in error &&
    typeof error.message === "string" &&
    error.message.trim() !== ""
  );
}
