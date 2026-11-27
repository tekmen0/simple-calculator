import { describe, expect, it } from "vitest";
import { parseOperand } from "./parse-operand";

describe("parseOperand", () => {
  it.each([
    ["0", 0],
    ["-0", -0],
    ["+42", 42],
    ["0007", 7],
    ["-12.75", -12.75],
    [".5", 0.5],
    ["2.", 2],
    ["1.e2", 100],
    ["-.5E+2", -50],
    ["6.25e-2", 0.0625],
    ["\t +4.5 \r\n", 4.5],
    ["1.7976931348623157e308", Number.MAX_VALUE],
    ["5e-324", Number.MIN_VALUE],
    ["1e-324", 0],
  ])("parses %j as %s", (input, expected) => {
    expect(parseOperand(input)).toBe(expected);
  });

  it.each([
    "",
    " \t\n",
    "-",
    ".",
    "1..2",
    "1e",
    "1e+",
    "e2",
    "1e1.5",
    "1e2e3",
    "12abc",
    "1 2",
    "1\n2",
    "0x10",
    "0b10",
    "0o10",
    "1_000",
    "1,5",
    "−1",
    "１２",
    "NaN",
    "Infinity",
    "1e309",
    "-1e309",
  ])("rejects %j", (input) => {
    expect(parseOperand(input)).toBeNull();
  });
});
