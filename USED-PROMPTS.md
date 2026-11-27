### Used prompts
#### 1 - Design : 
Read TASK.md and inspect the current repository structure.

Propose an architecture (using code-base design skill), for this take-home assignment using Go for the backend and React with TypeScript for the frontend.

Create AGENTS.md on project root to record our shared conventions for backend & frontend.&#x20;

Consider the following:

- Idiomatic Go and TypeScript conventions.
- Separation of business logic from HTTP transport.
- Clean REST API design and consistent error handling.
- Testability of both frontend and backend.
- Minimal external dependencies.
- Maintainability without overengineering.

Prioritize simplicity, correctness, and the explicit requirements in TASK.md.

Present the proposed directory structure, API contract, testing strategy, and architectural trade-offs. Explain your reasoning.

Do not implement anything yet. I want to review the architecture before proceeding.

Propose init commands for each repo & libraries to be used.



#### 2 - Implementation (Backend) : 
1. Given our design decisions at README.md & using available idiomatic golang skills implement the backend of the project. Which should be extendable to the tasks described in TASK.md.

2. Implement the request logging middleware as well.


#### 3 - Implementation (Frontend) : 
1. Given our design decisions at README.md & using available typescript skills, implement the simple fronted of the project. Which should be extendable to the tasks described in TASK.md. Add tests, vite form & https connections with browser fetch.


2. Update the ui, we select the operation from the dropdown between the first & second operands for more intuitive design which resemble to drawing operation on paper from left to right. And show the latest operation + result in the resulting box. Not only the result. For example, if i calculate 65 + 35, result container should show : 
65 + 35 = 100 instead of 100.



#### 4 - Docs verification : 

