# Project Context: Chia Florist Control Panel & WAF
This document provides the current context, architecture, and recent updates for the Chia Florist control-panel workspace. AI agents should use this document to understand the project structure and ongoing tasks before generating code or modifying files.

## 📌 Project Overview
The control-panel is a hybrid workspace encompassing the frontend dashboard for Chia Florist and a newly integrated, custom-built Web Application Firewall (WAF) written in Golang.

Core Objectives:
- Maintain a fast, modern frontend using Vite, TypeScript, and Tailwind CSS.
- Secure the application using a robust Golang-based WAF.
- Continuously test and monitor firewall rules against spam and simulated attacks.

## 🚀 Recent Updates & Migrations
- **Golang WAF Migration:** The WAF system has been successfully migrated to Golang (see `WAF_MIGRATION_REPORT.md` and `CHANGELOG_WAF_GOLANG.md` for historical context). The compiled executable is `main.exe`.
- **WAF Configuration Files:** The firewall is now fully driven by JSON configurations:
  - `waf-rules.json`: Defines the active security rules.
  - `waf-filters.json`: Defines specific filter parameters.
  - `waf-blocked.json`: Records IP/entities currently blocked.
  - `waf-logs.json`: General WAF activity logging.
- **Attack Simulation Tooling:** Added Node.js/CommonJS scripts (`simulate_attacks.js`, `simulate_spam.js`, and their `.cjs` variants) to rigorously test the WAF rules in development.

## 📂 Workspace Structure
The workspace contains both Go and Node.js ecosystems:

**Frontend (UI & Tooling)**
- `src/` & `public/`: Contains the frontend source code and static assets.
- `vite.config.ts`, `tsconfig.*.json`: Build and TypeScript configurations.
- `tailwind.config.js`, `postcss.config.js`: Styling engine configurations.
- `components.json`: Likely used for a UI library (e.g., shadcn/ui).
- `package.json` & `package-lock.json`: Node dependencies.

**Backend (Golang WAF)**
- `main.go`: The main entry point for the Golang backend/WAF.
- `go.mod`: Go module dependencies.
- `main.exe`: The compiled WAF binary.

**Documentation & AI Context**
- `Task-md/`: Folder containing task breakdowns and specific AI instructions.
- `PROJECT_CONTEXT_FOR_AI.md`: This file (or similar) to guide AI agents.
- `notefromme.txt`: Personal notes or immediate directives for the AI.

## 📋 Current TODOs & Focus Areas
- [x] Migrate WAF to Golang.
- [x] Set up JSON-based rule engine (waf-rules.json).
- [x] Create attack simulation scripts.
- [ ] Review and analyze recent waf-logs.json for false positives.
- [ ] Optimize main.go for better request throughput.
- [ ] Connect the frontend dashboard to the WAF configuration files (allow editing waf-rules.json via UI).

## 🛠️ Directives for AI Agents
When interacting with this workspace:
- **Language & Tech Stack:** Write frontend code in TypeScript/React (Vite) and backend WAF code in Golang.
- **Testing WAF:** If modifying `main.go` or WAF rules, always suggest running `node simulate_attacks.js` to verify the changes do not break existing protections.
- **Log Modification:** Do not manually alter `waf-logs.json` or `waf-blocked.json` via scripts unless explicitly instructed to build a flush/clean mechanism.
