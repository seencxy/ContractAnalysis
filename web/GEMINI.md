# GEMINI Project Context: ContractAnalysis Web

This document provides a comprehensive overview of the `ContractAnalysis/web` project, intended to guide AI-driven development and maintenance.

## Project Overview

This is the frontend for a Binance Futures Analysis system, designed to monitor, analyze, and display cryptocurrency contract signals. It's a modern, data-intensive single-page application (SPA).

The application provides functionalities like a real-time dashboard, signal monitoring, strategy analysis, historical data review, and leaderboards for trading pairs.

### Key Technologies

- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **UI Library**: Ant Design 5
- **Charting**: Apache ECharts
- **State Management**:
    - **Server State**: TanStack Query (React Query) for data fetching, caching, and synchronization.
    - **Client State**: Zustand for lightweight global state management.
- **Routing**: React Router v6
- **HTTP Client**: Axios

## Building and Running

The project uses `npm` as the package manager.

- **Install Dependencies**:
  ```bash
  npm install
  ```

- **Run Development Server**:
  Starts the app in development mode with HMR.
  ```bash
  npm run dev
  ```
  The frontend will be available at `http://localhost:5173`.

- **Build for Production**:
  Transpiles TypeScript and bundles the application for production.
  ```bash
  npm run build
  ```
  Output is placed in the `dist/` directory.

- **Lint Files**:
  Analyzes the code for style and syntax errors.
  ```bash
  npm run lint
  ```

- **Preview Production Build**:
  Serves the `dist/` directory to preview the production application.
  ```bash
  npm run preview
  ```

## Development Conventions

### Architecture

The codebase is structured logically by feature and responsibility:

- `src/api/`: Contains Axios instance configuration and typed API endpoint definitions.
- `src/components/`: Shared UI components, categorized into `common`, `layout`, and specialized `charts`.
- `src/hooks/`: Custom hooks, primarily for data fetching with TanStack Query (`src/hooks/queries/`).
- `src/pages/`: Top-level components for each route/page of the application.
- `src/types/`: Centralized TypeScript type and interface definitions.
- `src/utils/`: Utility functions for formatting, colors, etc.
- `src/router.tsx`: Defines all application routes using React Router.

### State Management

- **TanStack Query (`@tanstack/react-query`)**: Used for all interactions with the backend API. Custom hooks in `src/hooks/queries/` encapsulate data fetching logic for specific resources (e.g., `useSignals`, `useStatistics`).
- **Zustand**: Employed for managing global, cross-component client state that is not fetched from the server.

### API Interaction

- **Proxy**: Vite is configured to proxy requests from `/api` to the backend server, which is expected to be running at `http://localhost:8080`.
- **Response Format**: The backend has a standardized response wrapper. All data returned from the API is nested within a `data` property.
  ```typescript
  {
    code: 200,
    message: "success",
    data: { /* Actual payload */ },
    timestamp: 1234567890
  }
  ```
  Data access in React Query hooks should account for this structure (e.g., `response.data`).

### Code Style and Quality

- **Linter**: ESLint is configured to enforce code quality and consistency. The configuration can be found in `eslint.config.js`.
- **TypeScript**: The project is fully typed. Path aliases are configured, with `@/` pointing to the `src/` directory.
- **Component Design**: The use of Ant Design suggests a preference for using its components and design system for the UI. Custom components should be built to integrate well with this system.
