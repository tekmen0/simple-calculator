import { useState, type SubmitEvent } from "react";
import { calculate, type Operation } from "./calculator-api";
import { parseOperand } from "./parse-operand";

const operations = [
  { value: "add", label: "Add", symbol: "+" },
  { value: "subtract", label: "Subtract", symbol: "−" },
  { value: "multiply", label: "Multiply", symbol: "×" },
  { value: "divide", label: "Divide", symbol: "÷" },
] as const satisfies readonly { value: Operation; label: string; symbol: string }[];

type Submission =
  | { status: "idle" }
  | { status: "pending" }
  | { status: "success"; equation: string }
  | { status: "error"; message: string };

export function Calculator() {
  const [operation, setOperation] = useState<(typeof operations)[number]>(operations[0]);
  const [inputs, setInputs] = useState({ left: "", right: "" });
  const [showValidation, setShowValidation] = useState(false);
  const [submission, setSubmission] = useState<Submission>({ status: "idle" });
  const pending = submission.status === "pending";
  const operands = { left: parseOperand(inputs.left), right: parseOperand(inputs.right) };

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending) {
      return;
    }

    const { left, right } = operands;
    if (left === null || right === null) {
      setShowValidation(true);
      const invalidField = event.currentTarget.elements.namedItem(left === null ? "left" : "right");
      if (invalidField instanceof HTMLInputElement) {
        invalidField.focus();
      }
      return;
    }

    setShowValidation(false);
    setSubmission({ status: "pending" });
    const response = await calculate(operation.value, { left, right });
    setSubmission(
      response.ok
        ? { status: "success", equation: `${left} ${operation.symbol} ${right} = ${response.result}` }
        : { status: "error", message: response.error.message },
    );
  }

  function renderOperand(name: "left" | "right", label: string) {
    const invalid = showValidation && operands[name] === null;
    const errorId = `${name}-error`;
    return (
      <div className="expression-field">
        <label htmlFor={name}>{label}</label>
        <input
          id={name}
          name={name}
          type="text"
          inputMode="decimal"
          autoComplete="off"
          required
          disabled={pending}
          value={inputs[name]}
          aria-invalid={invalid || undefined}
          aria-describedby={invalid ? errorId : undefined}
          onChange={(event) => {
            setInputs({ ...inputs, [name]: event.currentTarget.value });
          }}
        />
        {invalid && <p id={errorId} className="field-error">Enter a finite number.</p>}
      </div>
    );
  }

  let resultText = "—";
  if (submission.status === "success") {
    resultText = submission.equation;
  } else if (pending) {
    resultText = "Calculating…";
  }

  return (
    <main className="calculator">
      <h1>Calculator</h1>

      <form
        noValidate
        onSubmit={handleSubmit}
        onChange={() => {
          if (submission.status === "error") {
            setSubmission({ status: "idle" });
          }
        }}
      >
        <div className="result">
          <span id="result-label">Result</span>
          <output
            className="result-value"
            aria-labelledby="result-label"
            htmlFor="left operation right"
          >
            {submission.status !== "error" && resultText}
          </output>
          {submission.status === "error" && <p role="alert">{submission.message}</p>}
        </div>

        <div className="expression">
          {renderOperand("left", "First number")}
          <div className="expression-field">
            <label htmlFor="operation">Operation</label>
            <select
              id="operation"
              name="operation"
              value={operation.value}
              disabled={pending}
              onChange={(event) => {
                const selected = operations.find(({ value }) => value === event.currentTarget.value);
                if (selected) {
                  setOperation(selected);
                }
              }}
            >
              {operations.map(({ value, label, symbol }) => (
                <option key={value} value={value} aria-label={label}>
                  {symbol}
                </option>
              ))}
            </select>
          </div>
          {renderOperand("right", "Second number")}
        </div>

        <button type="submit" disabled={pending}>Calculate</button>
      </form>
    </main>
  );
}
