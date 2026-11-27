// @vitest-environment jsdom

import { cleanup, render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "../src/Calculator";

const fetchMock = vi.fn<typeof fetch>();

beforeEach(() => {
  vi.stubGlobal("fetch", fetchMock);
});

afterEach(() => {
  cleanup();
  fetchMock.mockReset();
  vi.unstubAllGlobals();
});

describe("Calculator", () => {
  it("places the operation between the operands in tab order and displays the equation", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(Response.json({ result: 0.4 }));
    render(<Calculator />);

    await user.tab();
    expect(document.activeElement).toBe(screen.getByRole("textbox", { name: "First number" }));
    await user.keyboard("-2");
    await user.tab();
    const operation = screen.getByRole("combobox", { name: "Operation" });
    expect(document.activeElement).toBe(operation);
    expect(within(operation).getAllByRole("option").map((option) => option.textContent)).toEqual([
      "+", "−", "×", "÷",
    ]);
    expect(within(operation).getByRole("option", { name: "Divide" }).textContent).toBe("÷");
    await user.selectOptions(operation, "divide");
    await user.tab();
    expect(document.activeElement).toBe(screen.getByRole("textbox", { name: "Second number" }));
    await user.keyboard("-.5e1{Enter}");

    const result = screen.getByRole("status", { name: "Result" });
    await within(result).findByText("-2 ÷ -5 = 0.4");
    expect(fetchMock).toHaveBeenCalledWith(
      "/api/calculations/divide",
      expect.objectContaining({ body: '{"left":-2,"right":-5}' }),
    );
  });

  it("marks missing operands and clears each error when the input is corrected", async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "Calculate" }));

    const first = screen.getByRole("textbox", { name: "First number" });
    const second = screen.getByRole("textbox", { name: "Second number" });
    expect(document.activeElement).toBe(first);
    expect(first.getAttribute("aria-invalid")).toBe("true");
    expect(second.getAttribute("aria-invalid")).toBe("true");
    expect(screen.getAllByText("Enter a finite number.")).toHaveLength(2);
    expect(fetchMock).not.toHaveBeenCalled();

    await user.type(first, "0");
    expect(first.hasAttribute("aria-invalid")).toBe(false);
    expect(second.getAttribute("aria-invalid")).toBe("true");
  });

  it.each(["abc", "0x10", "1e999"])("rejects %s before calling the service", async (input) => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.type(screen.getByRole("textbox", { name: "First number" }), "1");
    const second = screen.getByRole("textbox", { name: "Second number" });
    await user.type(second, input);
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(document.activeElement).toBe(second);
    expect(screen.getByText("Enter a finite number.")).toBeDefined();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("prevents overlapping requests and shows the submitted equation while the form is edited", async () => {
    const user = userEvent.setup();
    const request = Promise.withResolvers<Response>();
    fetchMock
      .mockReturnValueOnce(request.promise)
      .mockResolvedValueOnce(Response.json({ result: 245 }));
    render(<Calculator />);

    const first = screen.getByRole("textbox", { name: "First number" });
    await user.type(first, "65");
    await user.type(screen.getByRole("textbox", { name: "Second number" }), "35");
    const operation = screen.getByRole<HTMLSelectElement>("combobox", { name: "Operation" });
    const submit = screen.getByRole<HTMLButtonElement>("button", { name: "Calculate" });
    await user.click(submit);

    const result = screen.getByRole("status", { name: "Result" });
    expect(result.textContent).toBe("Calculating…");
    expect(submit.disabled).toBe(true);
    expect(operation.disabled).toBe(true);
    await user.click(submit);
    expect(fetchMock).toHaveBeenCalledTimes(1);

    request.resolve(Response.json({ result: 100 }));
    await within(result).findByText("65 + 35 = 100");
    expect(submit.disabled).toBe(false);
    await user.clear(first);
    await user.type(first, "7");
    await user.selectOptions(operation, "multiply");
    expect(result.textContent).toBe("65 + 35 = 100");
    await user.click(submit);
    await within(result).findByText("7 × 35 = 245");
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("shows a service error and lets the user correct and resubmit the inputs", async () => {
    const user = userEvent.setup();
    fetchMock
      .mockResolvedValueOnce(
        Response.json(
          { error: { code: "division_by_zero", message: "Cannot divide by zero." } },
          { status: 422 },
        ),
      )
      .mockResolvedValueOnce(Response.json({ result: 3.5 }));
    render(<Calculator />);

    await user.selectOptions(screen.getByRole("combobox", { name: "Operation" }), "divide");
    await user.type(screen.getByRole("textbox", { name: "First number" }), "7");
    const second = screen.getByRole("textbox", { name: "Second number" });
    await user.type(second, "0");
    const submit = screen.getByRole("button", { name: "Calculate" });
    await user.click(submit);

    expect((await screen.findByRole("alert")).textContent).toBe("Cannot divide by zero.");
    await user.clear(second);
    expect(screen.queryByRole("alert")).toBeNull();
    await user.type(second, "2");
    await user.click(submit);
    await within(screen.getByRole("status", { name: "Result" })).findByText("7 ÷ 2 = 3.5");
  });
});
