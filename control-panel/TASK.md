### Frontend Dashboard Development Task

### Objective

Build the initial frontend structure for a dashboard application using **React** and **shadcn/ui**.

The goal is to create a clean starter application with core navigation, mock data, and page scaffolding that can be connected to real backend services later.

### Tech Stack

#### Core
- React
- TypeScript
- Vite

#### UI
- shadcn/ui
- TailwindCSS
- Recharts (for graphs)

#### State / Utilities
Use common lightweight patterns most React developers use.

Recommended:
- React Context (if needed)
- React Router
- Utility-first folder organization

### Project Structure

Use a commonly accepted React structure.

```txt
src/
  components/
    ui/
    shared/
    dashboard/

  pages/
    auth/
    dashboard/
    security/
    shop/
    product/
    order/
    transaction/
    shipment/

  layouts/
    DashboardLayout.tsx

  data/
    mockAccounts.ts
    mockSales.ts
    mockSecurity.ts
    mockAIInsight.ts
    mockLogs.ts

  routes/
    index.tsx

  lib/

  hooks/

  types/

  App.tsx
````

### Important Notes

All data should be hardcoded for now.

Since backend services do not exist yet, create mock static data inside:

```txt
/src/data
```

Mock:

* account data
* selling analytics
* AI insights
* security metrics
* WAF logs
* summary cards

Do not integrate APIs yet.

### Pages To Build

### 1. Login Page

Build a clean login page.

Requirements:

* email input
* password input
* login button
* app branding / logo placeholder
* simple validation UI
* redirect to dashboard after login (mock)

### 2. Dashboard Home Page

Main landing page after login.

Must contain:

#### Header

Show:

* page title
* user profile/avatar
* notification placeholder

#### Sidebar Navigation

Navigation items:

* Dashboard
* Security
* Shop
* Products
* Orders
* Transactions
* Shipments

#### Main Dashboard Content

Include:

##### Selling Graph

Example:

* daily sales chart
* weekly revenue chart

Use Recharts.

##### AI Insight Card

Hardcoded AI-generated business insights.

Example:

* top selling product
* unusual traffic spike
* conversion suggestion
* revenue trend prediction

##### WAF Summary

Cards for:

* blocked requests
* suspicious traffic
* threat level
* active rules

##### WAF Logs Table

Mock log entries:

* timestamp
* IP
* endpoint
* attack type
* action blocked/allowed

### 3. Security Page

Dedicated page for security monitoring.

Include:

#### Security Overview Cards

* attack attempts
* blocked IPs
* anomaly score
* active incidents

#### Threat Graph

Trend visualization

#### Recent Security Logs

Detailed mock logs table

#### Rule Status

Simple table of WAF rules

Example:

* SQL Injection
* XSS
* Path Traversal
* Rate Limit

Status:

* active
* warning
* disabled

### Empty Placeholder Pages

Create empty starter pages for:

#### Shop Page

Text:
"Shop management coming soon"

#### Product Page

Text:
"Product management coming soon"

#### Order Page

Text:
"Order management coming soon"

#### Transaction Page

Text:
"Transaction management coming soon"

#### Shipment Page

Text:
"Shipment management coming soon"

Each page should still use dashboard layout and sidebar.

### Layout Requirements

### Dashboard Layout

Reusable layout containing:

#### Sidebar

Persistent desktop sidebar

#### Main Content Area

Responsive

#### Mobile Behavior

Collapsible sidebar

### UI Guidelines

Keep design:

* clean
* modern
* minimal

Use shadcn components wherever possible:

* Card
* Button
* Table
* Badge
* Input
* Tabs
* Sheet
* Dropdown Menu

### Mock Data Expectations

Create realistic sample values.

Examples:

#### Sales

* revenue history
* order count
* product performance

#### AI Insight

Natural language summaries

#### Security

Threat statistics

#### WAF Logs

At least 20 mock entries

### Deliverables

Complete:

* routing
* login flow mock
* dashboard layout
* all listed pages
* graphs
* hardcoded datasets
* responsive UI

Focus on frontend structure and maintainability.

No backend integration required.

```
```
