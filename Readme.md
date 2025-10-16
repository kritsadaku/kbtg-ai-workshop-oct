# Programming Workshop Assignments

This document outlines a series of programming workshops designed to build a complete application, from initial design and prototyping to backend development, testing, and advanced AI integration.

## Workshop Overview

| Workshop | Topic | Description |
|----------|-------|-------------|
| Workshop 1 | CSS & JavaScript Animation | Create custom CSS and JavaScript animations for a given HTML file. |
| Workshop 2 | AI-Powered Prototyping | Build and deploy a UI prototype from a Product Requirements Document (PRD) using an AI tool. |
| Workshop 3 | Frontend Development | Build a React frontend application using AI-assisted setup, component conversion, and a design token system for consistency. |
| Workshop 4 | Backend Development & Testing | Design and implement backend services with proper database integration, containerization, and automated testing workflows. |
| Workshop 5 | AI-Assisted Tooling & Testing with MCP | Integrate the Model Context Protocol (MCP) with development tools to automate coding and E2E testing tasks. |
| Workshop 6 | AI Instruction Development | Create a comprehensive AI instruction system including system prompts, documentation, and chat mode references to guide project development. |

## Workshop 1: CSS & JavaScript Animation

**Objective**: Animate a static HTML page for KBTG using custom CSS and JavaScript.

**Tasks**:
- You will be given a pre-existing `index.html` file.
- Your task is to write CSS and JavaScript to add animations at the specified points within the file.

## Workshop 2: AI-Powered Prototyping

**Objective**: Create and deploy a UI Proof-of-Concept (POC) from a Product Requirements Document (PRD).

**Tasks**:
- Build the user interface on an AI platform (e.g., v0.dev, lovable.dev, bolt.new) based on the provided PRD.
- Deploy the generated UI to a public hosting service.

## Workshop 3: Frontend Development

**Objective**: Build a complete frontend application using React and modern frontend tooling. Emphasize the use of design tokens to ensure visual consistency across pages and prepare the project for scalability.

**Tasks**:
- Setup a React project with TailwindCSS using AI-assisted scaffolding.
- Convert an existing Vue component into a React component.
- Create a new page using predefined design tokens to maintain a unified style.
- (Optional) Add Storybook to document and specify frontend components.

### Setup Guide

-   **Vite > React (Typescript)**: https://vite.dev/guide/
-   **Tailwind + Vite**: https://tailwindcss.com/docs/installation/using-vite
-   **Storybook**: https://storybook.js.org/docs/get-started/frameworks/react-vite
-   **Interaction Test**: https://storybook.js.org/docs/writing-tests/interaction-testing
-   **React Router**: https://reactrouter.com/start/declarative/installation

## Workshop 4: Backend Development & Testing

**Objective**: Build and test a backend system that supports the frontend application. Focus on API design, database integration, containerization, and ensuring code quality through unit testing.

**Tasks**:
- Initialize the backend project with the help of AI (Go + Fiber or any preferred stack).
- Implement APIs according to the specifications in Git (for docs and images) using an SQLite database.
- Design the database schema and output the result in Mermaid format.
- Write additional code and unit tests based on the provided test cases.

### Setup Guide

- **Go + Fiber** : https://docs.gofiber.io/

## Workshop 5: AI-Assisted Tooling & Testing with MCP

**Objective**: To integrate the Model Context Protocol (MCP) with development tools like Git and Playwright to automate and assist in coding and E2E testing tasks.

**Tasks**:
- Use `git-mcp` to connect to the provided GitHub repository and add new helper code based on a given specification.
- Use `playwright-mcp` to create an end-to-end test that verifies the contact information on the KBTG homepage is located in Nonthaburi.

### Setup Guide
- **Install MCP in VS Code**: https://code.visualstudio.com/docs/copilot/customization/mcp-servers
- **Git MCP**: https://gitmcp.io/
- **Playwright MCP**: https://github.com/microsoft/playwright-mcp

## Workshop 6: AI Instruction Development

**Objective**: To create comprehensive AI instructions for a project by developing a complete instruction system that includes coding standards, documentation, and reference materials.

**Tasks**:
- Create a **System Prompt** file at `.github/copilot-instructions.md` to define coding standards, high-level project ideas, and clear boundaries for the AI.
- Establish a `specs` folder to store all project **Documentation** as a collection of markdown files.
- Configure a **Chat Mode** to act as a database assistant by referencing a `scripts` folder for context.
- Develop essential **Scripts** for database tasks (e.g., seeding data from a CSV, counting records) and make them accessible to the Chat Mode prompt.

### Guideline
- **Copilot Instruction Prompt**: https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions
- **Copilot Chat Mode**: https://code.visualstudio.com/docs/copilot/customization/custom-chat-modes
- **Copilot Prompt Template**: https://code.visualstudio.com/docs/copilot/customization/prompt-files
- **awesome-copilot**: https://github.com/github/awesome-copilot
- **Spec kit**: https://github.com/github/spec-kit (have video guideline in repo)
