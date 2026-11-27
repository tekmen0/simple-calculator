import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Calculator } from "./Calculator";
import "./styles.css";

const root = document.getElementById("root");
if (!root) {
  throw new Error("The calculator root element is missing.");
}

createRoot(root).render(
  <StrictMode>
    <Calculator />
  </StrictMode>,
);
