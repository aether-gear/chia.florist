### Frontend & WAF Integration Context

### Objective

This document provides a summary of the current workspace state, recent structural changes, and the technical context of the unified `control-panel` project. This will serve as a foundational context for the AI to understand the project architecture before implementing new feature requests.

### Unified Architecture Overview

The project has transitioned from a fragmented structure into a unified, full-stack ecosystem:
- **Frontend**: React + Vite + TypeScript + shadcn/ui.
- **Backend (WAF)**: A lightweight, highly concurrent Golang server (`main.go`) managing security rules, rate limiting, and IP blacklisting.
- **Persistence**: Flat-file JSON databases (`waf-logs.json`, `waf-rules.json`, `waf-blocked.json`, etc.) acting as the single source of truth for both systems.

### Visual Project Structure

```txt
d:\Antigravity-Workspace\web-dummy\chia.florist\control-panel\
│
├── main.go                     # Golang Backend (WAF Server & API endpoints)
├── main.exe                    # Compiled backend binary
├── go.mod                      # Go dependencies
│
├── waf-logs.json               # Live persistence for WAF traffic logs
├── waf-rules.json              # Live persistence for Regex Threat Rules
├── waf-blocked.json            # Live persistence for IP Bans/Whitelists
├── waf-filters.json            # Live persistence for custom keyword filters
│
├── simulate_spam.cjs           # Script to simulate DDoS/Rate Limit attacks
├── simulate_attacks.cjs        # Script to simulate SQLi/XSS/LFI attacks
│
├── src/                        # React Frontend Source
│   ├── components/             # Reusable UI components (shadcn/ui)
│   ├── data/                   # Mock data / legacy static data
│   ├── layouts/                # Dashboard Layouts (Sidebar, Header)
│   ├── pages/                  
│   │   ├── auth/               
│   │   ├── dashboard/          
│   │   └── security/           
│   │       └── SecurityPage.tsx # Highly interactive real-time WAF dashboard
│   │
│   ├── routes/                 # React Router definitions
│   ├── App.tsx                 # Root React Component
│   ├── index.css               # Global Tailwind styles
│   └── main.tsx                # React DOM entry point
│
├── package.json                # Node.js dependencies (Vite, React, Tailwind, Recharts)
├── tailwind.config.js          # Tailwind CSS configuration
├── vite.config.ts              # Vite bundler configuration
└── TASK.md                     # Original Frontend scaffolding task document
```

### What Has Changed (Recent Updates)

#### 1. Live Backend Integration
The frontend is no longer relying solely on hardcoded mock data. The `SecurityPage.tsx` now actively polls the Golang backend (`http://localhost:8080/api/...`) every 3 seconds for real-time traffic statistics, IP status, and rule states.

#### 2. Advanced Interactive Threat Landscape Graph
- **Real-time Data Parsing**: The graph aggregates live traffic logs into 10-minute intervals.
- **Full-Day Pre-population**: The "Today" view guarantees a full 24-hour visual spectrum (00:00 to 24:00) regardless of log sparsity.
- **Interactive Brushing (Click & Drag)**: Administrators can click, drag, and release on the graph to seamlessly zoom into a specific time window. 
- **Auto-Switch Custom Mode**: Dragging on the graph automatically switches the global view to "Custom Range", dynamically populating the `<Input type="datetime-local">` fields and filtering both the graph and the recent logs table instantly.

#### 3. Real-time Security Controls
- **IP Blacklist Manager**: A functional interface allowing admins to manually ban, whitelist, or mute IPs. Changes trigger immediate backend enforcement and persistence into `waf-blocked.json`.
- **Active Rules Registry**: Admins can seamlessly toggle specific attack vector protections (SQLi, XSS, etc.) on or off.
- **Dynamic Log Filtering**: The recent logs table now includes intuitive filters for **Row Count** (5, 10, 20, 50) and **Status** (All, Passed, Blocked).

#### 4. Automated Threat Mitigation Proof-of-Concept
Running scripts like `simulate_spam.cjs` successfully triggers the Golang backend's rate-limiting protocols (yielding HTTP 429 errors), which immediately reflects on the React dashboard as blocked IP events.

### Next Steps & Custom Feature Requests

*(Insert your friend's new feature request or prompt details here. The AI now has full context of the Vite + Golang WAF architecture and recent UX updates.)*
