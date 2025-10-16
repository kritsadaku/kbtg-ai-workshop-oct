# Programming Workshop Assignments

This document outlines a series of programming workshops designed to build a complete frontend application, from initial design to final refactoring.

## Workshop Overview

| Workshop | Topic | Description |
|----------|-------|-------------|
| Workshop 1 | CSS & JavaScript Animation | Create custom CSS and JavaScript animations for a given HTML file. |
| Workshop 2 | AI-Powered Prototyping | Build and deploy a UI prototype from a Product Requirements Document (PRD) using an AI tool. |
| Workshop 3 | Frontend Development | -- need text -- |
| Workshop 4 | Backend Development & Testing | Design and implement backend services with proper database integration, containerization, and automated testing workflows. |
| Workshop 5 | AI Instruction Development | Create comprehensive AI instructions including system prompts, documentation, chat mode references, and essential scripts for project development. |

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
- Setup a React project with TailwindCSS and React Router through AI-assisted scaffolding.  
- Define design tokens through a UI design system.  
- Create new pages using the design tokens to maintain a unified style.  

### Setup Guide

-   **Vite > React (Typescript)**: https://vite.dev/guide/
-   **Tailwind + Vite**: https://tailwindcss.com/docs/installation/using-vite
-   **Storybook**: https://storybook.js.org/docs/get-started/frameworks/react-vite
-   **Interaction Test**: https://storybook.js.org/docs/writing-tests/interaction-testing
-   **React Router**: https://reactrouter.com/start/declarative/installation

## Workshop 4: Backend Development & Testing

**Objective**: Build and test a backend system that supports the frontend application. Focus on API design, database integration, containerization, and ensuring code quality through unit, integration, and end-to-end testing.

**Tasks**:
- Initialize the backend project with the help of AI (Go + Fiber, Python + FastAPI, or any preferred stack).  
- Convert the UI specification into Swagger documentation, database schema, and backend code (using docs to guide development).  
- Create `Dockerfile` and `docker-compose` for configuration coverage.  
- Implement unit tests with code coverage reports.  
- Perform integration tests to validate database interactions.  
- Conduct end-to-end testing using Playwright, leveraging MCP for analysis support.  

### Setup Guide

- **Go + Fiber** : https://docs.gofiber.io/
- **Playwright**: https://playwright.dev/docs/intro#whats-installed

## Workshop 5: AI Instruction Development

**Objective**: To create comprehensive AI instructions for a project by developing a complete instruction system that includes coding standards, documentation, and reference materials.

**Tasks**:
- Use an existing repository (or create your own) to build your own instruction system, which must include:
  - **System Prompt** = `copilot-instruction.md` for defining coding standards, high-level ideas, and clear boundaries of what can and cannot be done
  - **Documentation** = Store `docs` in a dedicated folder (a collection of markdown files that compile project specifications)
  - **Chat Mode** = Specify reference data sources that can be used for citations
  - **Scripts** = Essential scripts that work in conjunction with the system (such as database seeding with initial values)

**Repository**:
- Repo for task: https://github.com/mikelopster/kbtg-be-go-lab
- MCP install in vs code: https://code.visualstudio.com/docs/copilot/customization/mcp-servers