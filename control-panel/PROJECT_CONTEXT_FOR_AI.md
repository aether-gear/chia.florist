# Project Context: Chia Florist Control Panel, WAF, & Staff Admin APIs

This document provides the current context, architecture, and feature descriptions for the Chia Florist control-panel workspace. AI agents should use this document to understand the project structure, features, and database mappings before generating code or modifying files.

## 📌 Project Overview
The control-panel is a hybrid workspace encompassing:
1. **Frontend Dashboard (Control Panel)**: Built using React, Vite, TypeScript, and Tailwind CSS.
2. **Web Application Firewall (WAF)**: Built using Golang, operating as a proxy filtering layer.
3. **Staff & Admin Management Dashboard**: Integrated directly with the `service-core` Golang backend to manage business operations.

---

## 🛠️ Feature Modules & Functionalities

### 1. Staff Profile & Session Management
- **Description**: Displays the currently logged-in staff member's profile and session parameters.
- **Features**:
  - View staff ID, user ID, account creation time, and last login.
  - Update personal information (Full Name, Phone Number, Avatar URL).
- **Backend API**: `GET /profile` (reads profile info), `PUT /profile` (saves profile updates).

### 2. Shops & Operations Management
- **Description**: Manages branch physical locations, shipping options, and product inventory levels.
- **Features**:
  - **General Info**: Configure shop name, description, and toggle shop active status.
  - **Addresses**: Add branch physical locations with full address parsing. Interacts with RajaOngkir to fetch dynamic Indonesian administrative regions (Province → City → District → Village).
  - **Couriers**: List and toggle active shipping courier codes (e.g. JNE, POS, TIKI) for the shop.
  - **Products Inventory**: Assign and add stock levels dynamically for specific products at the current shop.
- **Backend API**: `GET /shops` (lists shops), `POST /shops` (saves shop metadata), `POST /shops/{shopID}/addresses`, `POST /shops/{shopID}/products/{productID}/inventories` (adds inventory stock).

### 3. Staff & Merchant Workspace Management (Admin)
- **Description**: Provides platform-wide administrative controls to register new branch workspaces (staff entities) and assign user accounts.
- **Features**:
  - **Daftar Staff**: Retrieve a list of registered staff workspaces with details (creation dates, description, logo).
  - **Create Staff**: Creates a new staff/merchant workspace. Generates a unique user record and staff record within a unified transaction.
  - **Add Account**: Assigns a new user account (email, password, role) to an existing staff workspace.
- **Backend API**: `GET /staff` (lists workspaces), `POST /staff` (creates staff workspace), `POST /staff/{staffID}/accounts` (registers account to workspace).

### 4. Customer Database (Admin)
- **Description**: Displays registered customer profiles for auditing and platform management.
- **Features**:
  - Lists name, username, phone, and last login timestamp.
- **Backend API**: `GET /customers`.

### 5. Order Management
- **Description**: Operations team panel to monitor customer orders and fulfillment.
- **Features**:
  - Filter orders by status (Pending, Confirmed, Shipped, Delivered, Cancelled) and search by order number.
  - Order Detail view showing itemized breakdowns (price, quantity, subtotal, shipping fee, courier type) and customer profiles.
- **Backend API**: `GET /orders`.

### 6. Payment Settings (Admin)
- **Description**: Configuration of settlement options to receive customer payments.
- **Features**:
  - **Payment Methods**: List supported payment gateways (e.g. Bank Transfer, GoPay, QRIS).
  - **Payment Accounts**: Create and list settlement destination accounts (bank account numbers or e-wallet phone numbers). Handles optional values (e.g., bank accounts have `null` phone numbers, e-wallets have `null` account numbers).
- **Backend API**: `GET /payments/methods` (lists methods), `GET /payments/accounts` (lists accounts), `POST /payments/accounts` (saves new account).

---

## 🚀 Recent Critical Bugfixes & Migrations
- **Role Mismatch Fixed**: Corrected the backend Role constant `RoleStaffAdmin` value from `"admin"` to `"staff_admin"` to align with the SQL roles table seed code.
- **Database NULL Scan Mismatches**: Patched `PhoneNumber` types from `string` to pointer `*string` inside the Payment & Order domain models and DTO structures. This prevents Go pgx driver crashes when scanning `NULL` database values for non-phone settlement methods.
- **Staff Creation Implementation**: Overhauled the backend `CreateStaff` HTTP handler and usecase from returning mock JSON to performing transactional SQL writes on `users` and `staff` tables.

---

## 📂 Workspace Structure
The workspace contains both Go and Node.js ecosystems:

**Frontend (UI & Tooling)**
- `src/` & `public/`: Contains the React dashboard source code and static assets.
- `vite.config.ts`, `tsconfig.json`: Build and TypeScript configuration.
- `package.json`: Node dependencies (react, lucide-react, react-router-dom, etc.).

**Backend (Golang WAF & service-core)**
- `service-core/`: Golang service core package representing the business APIs (Shops, Customers, Staff, Payments, Orders).
- `service-core/internal/bootstrap/router.go`: chi router mappings for all staff and admin API chains.
- `service-core/migrations/`: SQL migration files defining database schema layout.

---

## 📋 Ongoing Focus Areas
- [ ] Connect the frontend dashboard to the WAF configuration files (allow editing `waf-rules.json` via UI).
- [ ] Review and analyze `waf-logs.json` for false positives in attack detections.
- [ ] Implement manual shipment dispatch integrations in `/shipments`.
