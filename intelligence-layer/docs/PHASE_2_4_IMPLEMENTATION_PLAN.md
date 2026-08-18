# 🚚 Phase 2.4 Implementation Plan: Courier SLA & Delivery Duration Estimator

> **Module**: `ai-lab`  
> **Phase**: 2.4 — Courier SLA Model  
> **Date**: 2026-08-11  
> **Predecessor**: [PHASE_2_3_IMPLEMENTATION_PLAN.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/PHASE_2_3_IMPLEMENTATION_PLAN.md)  
> **Reference**: [MODEL_TRAINING_REPORT.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/MODEL_TRAINING_REPORT.md)  
> **Status**: Proposed / Pending Review  

---

## Executive Summary

Phase 2.4 trains the **Courier SLA & Delivery Duration Estimator** — a regression model that predicts **expected delivery duration in hours** for each shipment given the courier, shipping cost, day of dispatch, and order value. Staff use this at dispatch time to compare couriers and select the most reliable option.

**Key difference from Phases 2.1–2.3**: Those phases consumed feature matrices fully prepared by Phase 1. Phase 2.4 requires **extending Phase 1** — the courier SLA feature matrix was never built. This phase adds a `CourierSLAFeatureBuilder` to `src/feature_engineering.py`, registers it in `extract.py`, and then trains the model.

**Available source data** (from Phase 1 extractor):
- `shipments`: 1,396 samples with `courier_name`, `shipping_cost`, `created_at`, `delivered_at`.
- All three major domestic couriers present: `JNE (464)`, `JNT (492)`, `SICEPAT (440)`.

---

## 1. Key Design Decisions

### 1.1 Target Variable
The target is **`delivery_duration_hours`** (continuous float), derived as:
```
(delivered_at - created_at).total_seconds() / 3600
```
This is a **regression** problem. Evaluation uses **MAE** and **RMSE** in hours.

### 1.2 Algorithm: XGBoost Regressor
Same algorithm as Phase 2.1:

| Criterion | Choice |
|---|---|
| **Dataset size** | 1,396 samples — ideal for gradient boosted trees |
| **Feature types** | Mix of numeric (shipping cost) and categorical (courier, day of week) — XGBoost handles both after one-hot encoding |
| **Speed** | Sub-millisecond inference — critical for real-time dispatch recommendations |
| **Interpretability** | Feature importance reveals which courier/feature drives delivery delay |

### 1.3 Per-Courier Reliability Score
In addition to raw duration prediction, the evaluation step computes a **per-courier reliability score** (0–100):
```
reliability[courier] = 100 × (1 - courier_mae / max(all_courier_mae))
```
This gives staff a ranked recommendation like: `JNE: 87, JNT: 74, SICEPAT: 91`.

---

## 2. Architecture Overview

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │              Phase 1 Raw Data (from extractor)                              │
 │  shipments: courier_name, shipping_cost, created_at, delivered_at           │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │        PHASE 1 EXTENSION: src/feature_engineering.py                        │
 │        [MODIFY] Add CourierSLAFeatureBuilder                                │
 │                                                                             │
 │  Features:  courier_jne, courier_jnt, courier_sicepat (one-hot),           │
 │             shipping_cost, dispatch_day_of_week,                            │
 │             dispatch_hour, dispatch_is_weekend                              │
 │  Target:    delivery_duration_hours (continuous float)                      │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │        PHASE 1 EXTENSION: extract.py                                        │
 │        [MODIFY] Register CourierSLAFeatureBuilder in pipeline               │
 │        Output: data/processed/courier_sla_features.csv / .pt               │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │        src/courier_sla_trainer.py                                           │
 │                                                                             │
 │  1. Load courier_sla_features.csv → DataFrame                              │
 │  2. Chronological train/val split (80/20)                                  │
 │  3. Train XGBRegressor with early stopping (val RMSE)                       │
 │  4. Evaluate: MAE (hours), RMSE, R², per-courier MAE & reliability scores  │
 │  5. Serialize → models/courier_sla.json                                    │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │  models/courier_sla.json                                                   │
 │  models/courier_sla_eval.json   (MAE per courier + reliability scores)     │
 │  models/courier_sla_importance.png                                          │
 └─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Planned Components & Code Changes

### 3.1 Phase 1 Extensions

#### [MODIFY] [src/feature_engineering.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/feature_engineering.py)

Add `CourierSLAFeatureBuilder` class:
- Consumes `shipments` records from the raw extractor output.
- Derives features: `shipping_cost`, `dispatch_day_of_week`, `dispatch_hour`, `dispatch_is_weekend`, and one-hot encoded courier columns (`courier_jne`, `courier_jnt`, `courier_sicepat`).
- Target `target_label` = `delivery_duration_hours`.
- Filters out rows where `delivered_at` is `None` (undelivered shipments).

#### [MODIFY] [extract.py](file:///d:/__Projects/kage/chia.florist/ai-lab/extract.py)

Registers `CourierSLAFeatureBuilder` in the main pipeline and persists:
- `data/processed/courier_sla_features.csv`
- `data/processed/courier_sla_features.pt`

---

### 3.2 New Configuration

#### [NEW] [configs/courier_sla.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/courier_sla.yaml)

```yaml
experiment_name: "courier_sla_v1"
seed: 42

data:
  processed_csv: "data/processed/courier_sla_features.csv"
  target_column: "target_label"
  train_val_split: 0.8
  courier_columns: ["courier_jne", "courier_jnt", "courier_sicepat"]

model:
  n_estimators: 500
  max_depth: 5
  learning_rate: 0.05
  subsample: 0.8
  colsample_bytree: 0.8
  early_stopping_rounds: 30
  eval_metric: "rmse"

output:
  model_dir: "models"
  model_name: "courier_sla.json"
  evaluation_report_name: "courier_sla_eval.json"
  feature_importance_name: "courier_sla_importance.png"
```

---

### 3.3 Core Training Module

#### [NEW] [src/courier_sla_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/courier_sla_trainer.py)

Class `CourierSLATrainer` with:
- `load_and_split_data()` — Loads `courier_sla_features.csv`, chronological train/val split.
- `train(X_train, y_train, X_val, y_val)` — Trains `xgboost.XGBRegressor` with early stopping on val RMSE.
- `evaluate(model, X_train, y_train, X_val, y_val)` — Computes:
  - Overall: MAE (hours), RMSE (hours), R².
  - **Per-courier MAE breakdown** (groups val set by active courier one-hot column).
  - **Per-courier reliability score**: `100 × (1 - courier_mae / max(all courier mae))`.
- `save_feature_importance(model, feature_names)` — Horizontal bar chart.
- `save_model(model)` — Native XGBoost JSON checkpoint.

---

### 3.4 CLI Training Script

#### [NEW] [train_courier_sla.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_courier_sla.py)

```bash
# Rebuild Phase 1 courier SLA dataset then train
python extract.py --config configs/feature_extraction.yaml
python train_courier_sla.py --config configs/courier_sla.yaml

# Dry run
python train_courier_sla.py --config configs/courier_sla.yaml --dry-run
```

---

### 3.5 Unit Tests

#### [NEW] [tests/test_courier_sla.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_courier_sla.py)

Tests to cover:
- `test_courier_sla_feature_builder` — `CourierSLAFeatureBuilder` produces correct shape, one-hot encoding is valid, target is positive float.
- `test_data_loading_and_split` — CSV loads, target column separates, split indices are chronological.
- `test_trainer_fit_and_evaluation` — Regressor trains, MAE/RMSE/R² are valid floats, per-courier scores are computed.
- `test_reliability_score_range` — Reliability scores are in [0, 100] for all couriers.
- `test_model_serialization` — `.json` checkpoint exists and produces identical inference on reload.

---

## 4. Evaluation Targets

| Metric | Target | Notes |
|---|---|---|
| **Val MAE** | < 5.0 hours | Acceptable dispatch-time accuracy for daily operations |
| **Val RMSE** | < 8.0 hours | Penalizes large outlier delivery delays |
| **R²** | > 0.50 | Explains > 50% of delivery time variance |
| **Per-Courier Reliability Score** | All > 50 | Synthetic data is clean — all couriers should be distinguishable |

---

## 5. File Summary

| File | Status | Description |
|---|---|---|
| [src/feature_engineering.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/feature_engineering.py) | `[MODIFY]` | Add `CourierSLAFeatureBuilder` class |
| [extract.py](file:///d:/__Projects/kage/chia.florist/ai-lab/extract.py) | `[MODIFY]` | Register CourierSLAFeatureBuilder and persist `courier_sla_features.*` |
| [configs/courier_sla.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/courier_sla.yaml) | `[NEW]` | Regressor hyperparameters and I/O config |
| [src/courier_sla_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/courier_sla_trainer.py) | `[NEW]` | Trainer class: split → fit → per-courier MAE & reliability scores → save |
| [train_courier_sla.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_courier_sla.py) | `[NEW]` | CLI entry script |
| [tests/test_courier_sla.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_courier_sla.py) | `[NEW]` | Unit & integration test suite |
| `data/processed/courier_sla_features.csv` | `[GENERATED]` | Phase 1 extension output |
| `data/processed/courier_sla_features.pt` | `[GENERATED]` | PyTorch tensor version |
| `models/courier_sla.json` | `[GENERATED]` | Trained XGBoost regressor checkpoint |
| `models/courier_sla_eval.json` | `[GENERATED]` | MAE, RMSE, R², per-courier MAE, reliability scores |
| `models/courier_sla_importance.png` | `[GENERATED]` | Feature importance bar chart |

---

## 6. Verification Plan

### Automated Tests
```powershell
pytest tests/test_courier_sla.py -v
```

### Manual Verification
1. Run `python extract.py --config configs/feature_extraction.yaml` — confirm new `data/processed/courier_sla_features.csv` is created.
2. Run `python train_courier_sla.py --config configs/courier_sla.yaml --dry-run` — should exit cleanly.
3. Run `python train_courier_sla.py --config configs/courier_sla.yaml` — confirm MAE < 5.0 hours and per-courier reliability scores appear in log.
4. Inspect `models/courier_sla_eval.json` — confirm overall + per-courier breakdown.
5. Open `models/courier_sla_importance.png` — `shipping_cost` and `courier_*` flags should rank highest.
