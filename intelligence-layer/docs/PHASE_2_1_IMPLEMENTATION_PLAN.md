# 📈 Phase 2.1 Implementation Plan: Demand & Sales Forecasting Model

> **Module**: `ai-lab`  
> **Phase**: 2.1 — First Custom Model Training  
> **Date**: 2026-08-10  
> **Predecessor**: [PHASE_1_IMPLEMENTATION_PLAN.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/PHASE_1_IMPLEMENTATION_PLAN.md)  
> **Reference**: [MODEL_TRAINING_REPORT.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/MODEL_TRAINING_REPORT.md)  
> **Status**: Proposed / Pending Review  

---

## Executive Summary

Phase 2.1 trains the **Demand & Sales Forecasting Model** — the most strategically important model in the `ai-lab` stack. It gives staff 7-day ahead sales volume predictions per SKU & shop, enabling proactive procurement scheduling before seasonal demand spikes (Valentine's Day, Mother's Day, Graduation Season).

The model is trained on the **1,080-sample, 13-feature** matrix already produced by the Phase 1 pipeline (`data/processed/demand_forecasting_features.csv`). We use **XGBoost Regressor** as the primary algorithm — it outperforms Prophet for tabular feature-rich data, requires zero GPU infrastructure, serializes to a portable `.json` checkpoint, and runs inference in microseconds on a CPU.

---

## 1. Why XGBoost (Not Prophet) as Primary Model

| Criterion | Prophet | XGBoost Regressor |
|---|---|---|
| **Feature Support** | Date + holidays only (time-series native) | All 13 features: lags, rolling stats, margin, view count, seasonality flags |
| **Inference Speed** | ~200ms (decomposition at runtime) | **< 1ms** (tree traversal) |
| **Portability** | Pickle / JSON (Facebook library dependency) | `.json` (universally portable, no library version pinning) |
| **Explainability** | Trend + seasonality chart | **SHAP feature importance** (tells staff *why* demand is predicted) |
| **Dataset Size Fit** | Works best on 2+ years of data | **Excellent on 1,000–10,000 samples** |
| **Overfitting Risk** | Low (parametric) | Controlled via `max_depth`, `learning_rate`, `n_estimators` + early stopping |

> **Conclusion**: For our 90-day, 13-feature matrix, XGBoost is the right fit. Prophet will be introduced in a later sub-phase when multi-year historical data from live `service-core` is available.

---

## 2. Architecture Overview

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │              Phase 1 Output (ai-lab/data/processed/)                        │
 │  demand_forecasting_features.csv — 1,080 samples × 13 features             │
 │                                                                             │
 │  Features: gross_margin_pct, view_count, day_of_week, month,               │
 │            is_weekend, units_sold_lag_1/7/14/30,                            │
 │            rolling_7d_mean, rolling_7d_std,                                 │
 │            rolling_30d_mean, rolling_30d_std                                │
 │  Target: units_sold_next_7d (continuous float)                              │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │             src/demand_forecasting_trainer.py                               │
 │                                                                             │
 │  1. Load CSV → pandas DataFrame                                             │
 │  2. Train / Val split (80 / 20 chronological)                               │
 │  3. Train XGBRegressor with early stopping on val RMSE                     │
 │  4. Evaluate: MAE, RMSE, R², SHAP importance chart                         │
 │  5. Serialize model to models/demand_forecasting.json                       │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                      models/demand_forecasting.json                         │
 │              (Portable serialized XGBoost checkpoint)                       │
 └─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Planned Components & Code Changes

### 3.1 New Configuration

#### [NEW] [configs/demand_forecasting.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/demand_forecasting.yaml)

Training hyperparameters and I/O configuration for the demand forecasting model:

```yaml
experiment_name: "demand_forecasting_v1"
seed: 42

data:
  processed_csv: "data/processed/demand_forecasting_features.csv"
  target_column: "target_label"
  train_val_split: 0.8   # Chronological (not random) split

model:
  n_estimators: 500
  max_depth: 5
  learning_rate: 0.05
  subsample: 0.8
  colsample_bytree: 0.8
  early_stopping_rounds: 30
  eval_metric: "rmse"

output:
  model_path: "models/demand_forecasting.json"
  feature_importance_path: "models/demand_forecasting_importance.png"
  evaluation_report_path: "models/demand_forecasting_eval.json"
```

---

### 3.2 Core Training Module

#### [NEW] [src/demand_forecasting_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/demand_forecasting_trainer.py)

Class `DemandForecastingTrainer` with:
- `load_data(csv_path, target_col)` — Reads Phase 1 CSV, separates features/target, performs a **chronological** train/val split (no random shuffle — time-series integrity must be preserved).
- `train(X_train, y_train, X_val, y_val)` — Initializes and trains `xgboost.XGBRegressor` with early stopping callbacks watching `rmse` on the val set. Returns the fitted model.
- `evaluate(model, X_val, y_val)` — Computes **MAE**, **RMSE**, and **R²** metrics. Logs a summary table. Saves metrics to `models/demand_forecasting_eval.json`.
- `save_feature_importance(model, feature_names)` — Generates a feature importance bar chart saved to `models/demand_forecasting_importance.png`.
- `save_model(model, path)` — Serializes the trained model to `models/demand_forecasting.json` using XGBoost native JSON format.

---

### 3.3 CLI Training Script

#### [NEW] [train_demand.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_demand.py)

CLI entry point modelled after the existing [train.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train.py):

```bash
# Full training run
python train_demand.py --config configs/demand_forecasting.yaml

# Dry run (load data, skip training)
python train_demand.py --config configs/demand_forecasting.yaml --dry-run
```

---

### 3.4 Dependencies

#### [MODIFY] [requirements.txt](file:///d:/__Projects/kage/chia.florist/ai-lab/requirements.txt)

Add Phase 2.1 training dependencies:

```
xgboost>=2.0.0
pandas>=2.0.0
scikit-learn>=1.3.0
matplotlib>=3.7.0
shap>=0.44.0          # optional — SHAP importance plots, degrades gracefully if absent
```

---

### 3.5 Unit Tests

#### [NEW] [tests/test_demand_forecasting.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_demand_forecasting.py)

Tests to cover:
- `test_data_loading` — CSV loads correctly, target column separates, shapes are valid.
- `test_chronological_split` — Train rows are all earlier in index than val rows (no leakage).
- `test_trainer_fit` — Model trains without exception on a small subset.
- `test_evaluation_metrics` — MAE, RMSE, R² values are non-negative, finite floats.
- `test_model_serialization` — Saved `.json` file exists and can be reloaded for inference.

---

## 4. Evaluation Targets

The model should hit these benchmarks on the synthetic 90-day dataset before Phase 2.1 is considered complete:

| Metric | Target |
|---|---|
| **RMSE** | < 3.0 units (acceptable for daily SKU volume) |
| **MAE** | < 2.0 units |
| **R²** | > 0.50 (explains > 50% of variance) |

> [!NOTE]
> Targets will be re-calibrated once real `service-core` data is used. Synthetic data targets are intentionally achievable to validate the pipeline; real data will require hyperparameter re-tuning.

---

## 5. File Summary

| File | Status | Description |
|---|---|---|
| [configs/demand_forecasting.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/demand_forecasting.yaml) | `[NEW]` | Hyperparameters and I/O config for Phase 2.1 |
| [src/demand_forecasting_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/demand_forecasting_trainer.py) | `[NEW]` | Trainer class: load → split → fit → evaluate → save |
| [train_demand.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_demand.py) | `[NEW]` | CLI entry script for Phase 2.1 training |
| [tests/test_demand_forecasting.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_demand_forecasting.py) | `[NEW]` | Unit & integration tests for Phase 2.1 |
| [requirements.txt](file:///d:/__Projects/kage/chia.florist/ai-lab/requirements.txt) | `[MODIFY]` | Add xgboost, pandas, scikit-learn, matplotlib, shap |
| `models/demand_forecasting.json` | `[GENERATED]` | Trained XGBoost model checkpoint |
| `models/demand_forecasting_eval.json` | `[GENERATED]` | Evaluation metrics: MAE, RMSE, R² |
| `models/demand_forecasting_importance.png` | `[GENERATED]` | Feature importance bar chart |

---

## 6. Verification Plan

### Automated Tests
```powershell
pytest tests/test_demand_forecasting.py -v
```

### Manual Verification
1. Run `python train_demand.py --config configs/demand_forecasting.yaml --dry-run` — should exit cleanly.
2. Run `python train_demand.py --config configs/demand_forecasting.yaml` — should complete training and log MAE / RMSE / R².
3. Inspect `models/demand_forecasting.json` (should exist, non-empty).
4. Inspect `models/demand_forecasting_eval.json` to confirm metrics meet targets in §4.
5. Open `models/demand_forecasting_importance.png` to verify feature ranking is plausible (lag features and rolling stats should rank highest).
