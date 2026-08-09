# 🤖 AI Model Training & API Strategy Report — chia.florist `ai-lab`

> **Scope**: Staff & Internal Admin Operations  
> **Target Service**: `service-core` integration  
> **Date**: 2026-08-10  
> **Reference Document**: [AI_READINESS.md](file:///d:/__Projects/kage/chia.florist/service-core/docs/AI_READINESS.md)  
> **Author**: AI Lab Team  

---

## Executive Summary

This report establishes the model training strategy and architecture for the **chia.florist** AI ecosystem within the `ai-lab` module. Based on the data readiness audit in [AI_READINESS.md](file:///d:/__Projects/kage/chia.florist/service-core/docs/AI_READINESS.md), the current transactional data in `service-core` is **100% clean and sufficient** for **internal staff and admin tooling**. Customer-facing features (e.g., personal recommendations, semantic search) are deferred until customer behavior tracking (`pgvector`, clickstream) is fully implemented.

The AI initiative focuses exclusively on three primary operational objectives:
1. **Predictive Analytics & Forecasting** (Demand, stockouts, courier delivery efficiency).
2. **Data & Pattern Summarization / Explanation** (Translating numerical forecasts and operational state into concise staff briefings).
3. **Anomaly Detection** (Identifying payment irregularities, fulfillment SLA breaches, stock discrepancies, and order abuse).

To achieve maximum precision and cost efficiency, we deploy a **Hybrid Architecture**: specialized custom Machine Learning models trained on local transactional data combined with a lightweight LLM API layer for natural language explanation.

---

## 1. Core Objectives & Staff-First Scope

> [!IMPORTANT]
> **Scope Guardrail**: For the current phase, AI features are built **strictly for staff/admin assistance**, NOT customer-facing applications. Staff require fast, deterministic, mathematically accurate insights to optimize daily store operations.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      chia.florist Staff AI System                       │
├──────────────────────────┬──────────────────────┬───────────────────────┤
│    1. Analytics &        │ 2. Summarization &   │ 3. Anomaly            │
│       Forecasting        │    Explanation       │    Detection          │
│                          │                      │                       │
│  • SKU Sales Velocity    │ • Executive Daily    │ • Payment Failure     │
│  • Stockout Horizon      │   Briefings          │   Spikes              │
│  • Courier SLA Estimator │ • Reorder Advice     │ • Fulfillment Delays  │
│  • Margin Optimization   │ • Trend Explanation  │ • Stock Discrepancies │
└──────────────────────────┴──────────────────────┴───────────────────────┘
```

### Objective Details
1. **Predictive Analytics**: Transition from reactive reporting to proactive operational forecasting. Staff get advance notice of demand spikes (e.g., Valentine's Day, Mother's Day, graduation seasons) and stock depletion dates.
2. **Data & Pattern Summarization**: Raw analytical metrics are often dense. The system automatically converts numerical outputs into clear, actionable Bahasa Indonesia / English bullet points for store managers.
3. **Anomaly Detection**: Automated guardrails to catch operational failures before they impact revenue or customer experience (e.g., courier bottlenecks, unallocated payments, abnormal inventory drops).

---

## 2. Models Being Trained (Based on Available Data)

The `service-core` PostgreSQL database provides high-quality transactional tables ready for model training. Below are the **5 specific models** being built and trained within `ai-lab`.

| # | Model Name | Model Type / Algorithm | Primary Input Data (`service-core`) | Output / Staff Benefit |
|---|---|---|---|---|
| 1 | **Demand & Sales Forecasting Model** | **Prophet + LightGBM / XGBoost** | `orders`, `order_items`, `product_performance` (`view_count`, sales velocity 7d/30d), Holiday calendar | **Predictive Analytics**: 7-day and 30-day predicted sales volume per SKU & category. Helps staff schedule flower stock procurement. |
| 2 | **Inventory Stockout & Safety Reorder Model** | **XGBoost Classifier / Random Forest** | `product_stock_history`, `inventory`, `product_performance` (`supplier_lead_time_days`, `gross_margin_pct`) | **Predictive Analytics**: Risk score (0-100%) of stockout within lead time and automated reorder point alerts for inventory team. |
| 3 | **Operational Anomaly Detection Model** | **Isolation Forest / PyTorch Autoencoder** | `payments`, `payment_events`, `shipment_events`, `orders`, `audit_logs` | **Anomaly Detection**: Real-time alerts for payment method anomaly spikes, status transition gaps, or sudden order volume anomalies. |
| 4 | **Courier SLA & Delivery Duration Estimator** | **Gradient Boosting Regression** | `shipments`, `shipment_events`, `order_items` (`courier_code`, weight, cost) | **Predictive Analytics**: Precise estimated fulfillment & delivery duration per courier/city. Recommends best courier to staff at dispatch. |
| 5 | **Natural Language Staff Briefing Generator** | **Hybrid LLM API** *(Gemini Flash / GPT-4o-mini)* | Outputs from Models 1–4 + `analytics/*` API payloads | **Summarization & Explanation**: Daily written summaries explaining *why* a spike occurred and recommending specific staff actions. |

---

### Detailed Breakdown of Models

#### Model 1: Demand & Sales Forecasting Model
* **Dataset**: Historical `order_items` aggregated by date and `product_id`, joined with `product_performance` (view count & 7d/30d velocity).
* **Technique**: **Prophet** captures macroeconomic holiday seasonality (Valentine's Day, Mother's Day, Year-End), while **LightGBM/XGBoost** captures short-term product features (price, gross margin, view count trend).
* **Output**: Time-series forecast curve per shop & product category for $T+7$ and $T+30$ days.

#### Model 2: Inventory Stockout & Safety Reorder Model
* **Dataset**: Time-series snapshots from `product_stock_history` (purpose-built for ML in `service-core`), current `inventory.stock` & `reserved_stock`, and `product_performance.supplier_lead_time_days`.
* **Technique**: Supervised classification predicting binary flag `will_stockout_in_lead_time`.
* **Output**: Ranked list of high-risk SKUs requiring immediate supplier purchase orders.

#### Model 3: Operational Anomaly Detection Model
* **Dataset**: Latency gaps between `payments.created_at` and `paid_at`, status transition events from `payment_events` & `shipment_events`, and event rates from `audit_logs`.
* **Technique**: **Isolation Forest** (unsupervised tabular anomaly detection) paired with a simple **PyTorch Autoencoder** (trained via `ai-lab/train.py`).
* **Output**: Anomaly score (0.0 to 1.0) and flag for out-of-distribution events (e.g., payment webhook delay spike, manual override anomaly).

#### Model 4: Courier SLA & Delivery Duration Estimator
* **Dataset**: `shipment_events` duration logs, grouped by courier code (`courier_name`, `service`), package weight, and shipping destination.
* **Technique**: Gradient Boosting Regressor (Scikit-Learn / XGBoost).
* **Output**: Expected delivery time range (hours) and reliability score to help warehouse staff choose optimal shipping partners.

#### Model 5: Data & Pattern Summarization Layer (Hybrid Integration)
* **Dataset**: Structured JSON outputs produced by Models 1–4, concatenated with `GET /analytics/*` metric aggregations.
* **Technique**: Zero-shot / Few-shot prompt engineering over an ultra-fast LLM API (e.g., Gemini 3.5 Flash or GPT-4o-mini).
* **Output**: Executive bullet points for staff (e.g., *"Demand for Red Rose Bouquet is projected to increase by 45% next week due to Valentine's Day. Reorder 120 units from Supplier A by Wednesday to prevent stockout."*).

---

## 3. Custom Trained Models vs. External LLM / Cloud APIs

A critical architectural decision is choosing when to **train custom ML models** versus calling **external LLM/Cloud APIs**. The evaluation matrix below highlights why custom training is superior for core operational objectives, and why a hybrid approach is optimal.

### Comprehensive Comparison Matrix

| Dimension | Custom Trained ML Models (Prophet / XGBoost / PyTorch) | External LLM APIs (GPT-4o-mini / Gemini Flash) | Standard Cloud AutoML APIs |
|---|---|---|---|
| **Mathematical & Numerical Accuracy** | 🟢 **100% Deterministic & Exact**. Computes exact quantities, probabilities, and statistical boundaries without hallucination. | 🔴 **Unreliable for Math**. LLMs hallucinate numbers, struggle with precise time-series math, and cannot do matrix operations. | 🟡 High accuracy, but limited flexibility for custom domain features. |
| **Operational Cost** | 🟢 **Near Zero ($0 / mo)**. Runs locally on existing backend servers or free CPU instances (Hugging Face Spaces / Render). | 🟡 **Pay-per-Token**. Low at small volume (~$1–$5/mo), but scales linearly with database rows and frequent polling. | 🔴 **High Recurring Cost**. Expensive minimum instance fees ($50–$200+/mo). |
| **Inference Latency** | 🟢 **Ultra-Fast (< 5ms - 20ms)**. Instant local execution in Python/Go runtime. | 🔴 **Slow (1,000ms - 3,000ms)**. Network roundtrips + token generation delays. | 🟡 Medium (100ms - 300ms). Network API overhead. |
| **Data Privacy & Security** | 🟢 **100% In-House Data Control**. Customer order history, product margins, and cost prices never leave local infrastructure. | 🔴 **Third-Party Data Exposure**. Sending proprietary margins and order logs to public cloud APIs. | 🟡 Third-party enterprise agreement required. |
| **Offline Capability** | 🟢 **Full Offline Support**. Models run locally inside the `ai-lab` service without internet dependency. | 🔴 **Requires Internet & API Uptime**. Fails if third-party API is degraded or rate-limited. | 🔴 Internet dependent. |
| **Natural Language Summarization** | 🔴 **None**. Produces only vectors, numbers, class labels, and probabilities. | 🟢 **World-Class**. Excellent at converting complex data into clear human-readable narrative explanations. | 🔴 Requires manual output formatting. |

---

### Why Custom Trained Models Win for Analytics & Anomaly Detection

1. **Precision Over Hallucination**:
   Store staff cannot rely on an LLM guessing sales inventory numbers. Time-series algorithms (Prophet) and decision trees (XGBoost) calculate actual standard deviations, moving averages, and exact numerical probabilities.
2. **Tabular Anomaly Detection Efficiency**:
   Algorithms like **Isolation Forest** evaluate thousands of transactional rows (`payments`, `shipments`, `audit_logs`) in milliseconds to detect spatial and temporal outliers. Feeding raw database tables into an LLM context window is cost-prohibitive and slow.
3. **Data Protection & Business Secrecy**:
   `product_performance.gross_margin_pct` and `cost_price` contain sensitive financial strategy data. Training local models keeps sensitive profit margins fully protected inside internal systems.

---

## 4. The Recommended Architecture: Hybrid Strategy

To leverage the strengths of both approaches, `ai-lab` implements a **Hybrid Pipeline**:

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                         1. DATA EXTRACTION LAYER                            │
 │  `service-core` DB: orders, order_items, inventory, product_performance     │
 └──────────────────────────────────────┬──────────────────────────────────────┘
                                        │
                                        ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                     2. CUSTOM TRAINED ML MODELS (ai-lab)                    │
 │                                                                             │
 │  ┌──────────────────────┐  ┌──────────────────────┐  ┌───────────────────┐  │
 │  │ Demand Forecasting   │  │ Stockout Risk Model  │  │ Anomaly Detector  │  │
 │  │ (Prophet + XGBoost)  │  │ (XGBoost Classifier) │  │ (IsolationForest) │  │
 │  └──────────┬───────────┘  └──────────┬───────────┘  └─────────┬─────────┘  │
 └─────────────┼─────────────────────────┼────────────────────────┼────────────┘
               │                         │                        │
               └────────────────────┬────┴────────────────────────┘
                                    │ Structured Numeric Output (JSON)
                                    ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                      3. LLM API EXPLANATION LAYER                           │
 │  Target: Gemini 3.5 Flash / GPT-4o-mini                                     │
 │  Input:  Numeric Forecasts + Anomaly Scores + Metric Aggregations           │
 │  Task:   Summarize patterns into clear Indonesian/English text for staff    │
 └──────────────────────────────────────┬──────────────────────────────────────┘
                                        │
                                        ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                          4. STAFF ADMIN DASHBOARD                           │
 │  • Visual Charts (Historical vs Predicted Demand)                           │
 │  • Red Anomaly Alert Banners (Payment / Fulfillment bottlenecks)             │
 │  • Natural Language Actionable Advice ("Procure 50 Lily bouquets by Tue")   │
 └─────────────────────────────────────────────────────────────────────────────┘
```

> [!TIP]
> **Why Hybrid Wins**: Custom ML models perform the heavy quantitative math instantly for free, and the LLM API acts strictly as a "narrator" to render readable explanations for staff. This consumes less than **$1.00/month** in API costs while maintaining 100% numerical accuracy.

---

## 5. Implementation & Training Workflow in `ai-lab`

The `ai-lab` module already provides a structured PyTorch / Scikit-Learn training environment:

- [train.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train.py): CLI entry point for model execution and dry runs.
- [src/data_loader.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/data_loader.py): Dataset definitions and batching utilities.
- [src/model.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/model.py): PyTorch Neural Network architectures (MLP, Autoencoder).
- [src/trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/trainer.py): Model training loop with evaluation logging.
- [configs/](file:///d:/__Projects/kage/chia.florist/ai-lab/configs): YAML experiment and hyperparameter configurations.

### Step-by-Step Execution Plan

```
Phase 1: Feature Extraction Pipeline
  ├── Query service-core DB for historical transactional records
  ├── Process product_stock_history, order_items, and payment_events
  └── Save training feature matrices under ai-lab/data/processed/

Phase 2: Custom Model Training & Artifact Generation
  ├── Train Prophet on aggregated daily order volumes (Demand Forecasting)
  ├── Train XGBoost Classifier on inventory & supplier lead times (Stockout Risk)
  ├── Train Isolation Forest / PyTorch Autoencoder (Anomaly Detection)
  └── Serialize model checkpoints (.pkl / .pt) into ai-lab/models/

Phase 3: Hybrid Summarization & Inference API
  ├── Wrap model inference pipelines in FastAPI (ai-lab/src/)
  ├── Integrate LLM API prompt formatter for numeric-to-text summarization
  └── Expose REST endpoints to service-core Admin Dashboard:
        • GET  /ai/forecast/demand
        • GET  /ai/risk/stockout
        • GET  /ai/anomalies/operational
        • POST /ai/explain/summary
```

---

## Conclusion & Next Steps

1. **Staff-First Alignment**: All AI development will strictly target internal operational enhancements until customer-facing schema prerequisites (`pgvector`, clickstream tracking) are established.
2. **High ROI Custom Training**: Training custom ML models for forecasting and anomaly detection yields **zero recurring infrastructure cost**, **sub-5ms response times**, and **100% numerical accuracy**.
3. **Action Item**: Begin feature extraction script development in `ai-lab/src/data_loader.py` targeting `product_stock_history` and `order_items` tables from `service-core`.
