// Restrict Number() to decimal input so blank strings and radix prefixes are rejected.
const decimalNumber = /^[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?$/;

export function parseOperand(input: string): number | null {
  const text = input.trim();
  if (!decimalNumber.test(text)) {
    return null;
  }

  const value = Number(text);
  return Number.isFinite(value) ? value : null;
}
