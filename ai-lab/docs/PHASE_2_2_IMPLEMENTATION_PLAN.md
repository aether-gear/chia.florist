# 📦 Phase 2.2 Implementation Plan: Inventory Stockout & Safety Reorder Risk Model

> **Module**: `ai-lab`  
> **Phase**: 2.2 — Stockout Risk Classifier  
> **Date**: 2026-08-10  
> **Predecessor**: [PHASE_2_1_IMPLEMENTATION_PLAN.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/PHASE_2_1_IMPLEMENTATION_PLAN.md)  
> **Reference**: [MODEL_TRAINING_REPORT.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/MODEL_TRAINING_REPORT.md)  
> **Status**: Proposed / Pending Review  

---

## Executive Summary

Phase 2.2 trains the **Inventory Stockout & Safety Reorder Risk Classifier** — a binary classifier that predicts which SKUs will hit zero stock before the supplier's lead time expires.

When triggered, staff receive a ranked alert list of high-risk products, enabling preemptive purchase orders before stockouts impact customer-facing availability.

The model is trained on the **60-sample, 6-feature** matrix produced by Phase 1 (`data/processed/stockout_risk_features.csv`). The critical design constraint in this phase is **severe class imbalance**: the Phase 1 dataset contains **56 negatives (no stockout risk)** and only **4 positives (at-risk)** — a 14:1 ratio. A naive model trained without intervention would predict "no stockout" on every sample and score 93.3% accuracy while missing all actual stockouts. This phase uses **SMOTE oversampling** to balance the training data and **precision-recall metrics** (not accuracy) as the primary evaluation criteria.

---

## 1. Key Design Decision: Class Imbalance Handling

| Strategy | Approach | Verdict |
|---|---|---|
| **Naive Training** | Fit classifier on raw imbalanced data | 🔴 Rejected — model ignores minority class entirely |
| **Class Weights** | Pass `scale_pos_weight` to XGBoost | 🟡 Partial fix — does not generate new training signal |
| **SMOTE Oversampling** | Synthesize new minority class samples in training set only | 🟢 **Selected** — creates balanced training distribution without contaminating the val set |
| **Threshold Tuning** | Shift classification threshold below 0.5 on raw output | 🟡 Complement — applied after SMOTE for fine tuning |

> [!IMPORTANT]
> SMOTE is applied **only on the training split** after the train/val separation. The val set is always kept at the original 60-sample distribution to ensure evaluation reflects real-world conditions.

---

## 2. Architecture Overview

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │              Phase 1 Output (ai-lab/data/processed/)                        │
 │  stockout_risk_features.csv — 60 samples × 6 features                      │
 │  Class Distribution: {0: 56, 1: 4}  — 14:1 imbalance ratio                │
 │                                                                             │
 │  Features: stock, reserved_stock, stock_burn_rate_7d,                      │
 │            supplier_lead_time_days, estimated_days_to_stockout,            │
 │            reorder_urgency_ratio                                            │
 │  Target: stockout_within_lead_time (binary 0 / 1)                          │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │             src/stockout_risk_trainer.py                                    │
 │                                                                             │
 │  1. Load CSV → pandas DataFrame                                             │
 │  2. Stratified train/val split (80/20, preserves class ratio in val)        │
 │  3. Apply SMOTE to training set only (balances to 1:1 ratio)               │
 │  4. Train XGBClassifier with scale_pos_weight + early stopping             │
 │  5. Evaluate: Precision, Recall, F1, ROC-AUC, Classification Report        │
 │  6. Serialize model → models/stockout_risk.json                            │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │                      models/stockout_risk.json                              │
 │              (Portable serialized XGBoost classifier checkpoint)            │
 └─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Planned Components & Code Changes

### 3.1 New Configuration

#### [NEW] [configs/stockout_risk.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/stockout_risk.yaml)

Training hyperparameters and I/O configuration for the stockout risk classifier:

```yaml
experiment_name: "stockout_risk_v1"
seed: 42

data:
  processed_csv: "data/processed/stockout_risk_features.csv"
  target_column: "target_label"
  train_val_split: 0.8
  use_smote: true

model:
  n_estimators: 300
  max_depth: 4
  learning_rate: 0.05
  subsample: 0.8
  colsample_bytree: 0.8
  scale_pos_weight: 14
  early_stopping_rounds: 20
  eval_metric: "logloss"

output:
  model_dir: "models"
  model_name: "stockout_risk.json"
  evaluation_report_name: "stockout_risk_eval.json"
  feature_importance_name: "stockout_risk_importance.png"
```

---

### 3.2 Core Training Module

#### [NEW] [src/stockout_risk_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/stockout_risk_trainer.py)

Class `StockoutRiskTrainer` with:
- `load_and_split_data()` — Reads Phase 1 CSV, performs **stratified** train/val split to guarantee at least 1 positive sample in the val set.
- `apply_smote(X_train, y_train)` — Applies `imblearn.over_sampling.SMOTE` to training data only. Falls back gracefully if `imbalanced-learn` is not installed (logs warning, trains on raw imbalanced data).
- `train(X_train, y_train, X_val, y_val)` — Initializes and trains `xgboost.XGBClassifier` with `eval_metric="logloss"` and early stopping.
- `evaluate(model, X_val, y_val)` — Computes **Precision**, **Recall**, **F1**, and **ROC-AUC**. Outputs full `classification_report` to log. Saves metrics JSON.
- `save_feature_importance(model, feature_names)` — Horizontal importance bar chart.
- `save_model(model)` — Serializes to native XGBoost JSON.

---

### 3.3 CLI Training Script

#### [NEW] [train_stockout.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_stockout.py)

CLI entry point modelled after [train_demand.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_demand.py):

```bash
# Full training run
python train_stockout.py --config configs/stockout_risk.yaml

# Dry run (data loading and split validation only)
python train_stockout.py --config configs/stockout_risk.yaml --dry-run
```

---

### 3.4 Dependencies

#### [MODIFY] [requirements.txt](file:///d:/__Projects/kage/chia.florist/ai-lab/requirements.txt)

Add SMOTE dependency:

```
imbalanced-learn>=0.11.0
```

---

### 3.5 Unit Tests

#### [NEW] [tests/test_stockout_risk.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_stockout_risk.py)

Tests to cover:
- `test_data_loading` — CSV loads correctly, 6 feature columns, binary target.
- `test_stratified_split` — Val set contains at least 1 positive sample.
- `test_smote_balancing` — Post-SMOTE training set has equal class counts.
- `test_trainer_fit_and_metrics` — Model trains, F1 and ROC-AUC are valid floats.
- `test_model_serialization` — `.json` checkpoint exists and reloads for inference.

---

## 4. Evaluation Targets

Because this is a high-imbalance classification problem, **accuracy is not the primary metric**.

| Metric | Target | Why It Matters |
|---|---|---|
| **Recall (Sensitivity)** | > 0.70 | Catching 70%+ of real stockout events is critical — false negatives cost more than false positives |
| **Precision** | > 0.50 | Minimize staff alert fatigue from false positives |
| **F1 Score** | > 0.60 | Harmonic mean — balanced signal of both above |
| **ROC-AUC** | > 0.75 | Measures probability rank quality across all thresholds |

> [!NOTE]
> Targets are set for the synthetic 60-sample dataset. With real `service-core` data (more positive stockout events), metrics should significantly improve.

---

## 5. File Summary

| File | Status | Description |
|---|---|---|
| [configs/stockout_risk.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/stockout_risk.yaml) | `[NEW]` | Classifier hyperparameters and SMOTE flag |
| [src/stockout_risk_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/stockout_risk_trainer.py) | `[NEW]` | Stratified split, SMOTE, XGBClassifier, evaluation, serialization |
| [train_stockout.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_stockout.py) | `[NEW]` | CLI entry script |
| [tests/test_stockout_risk.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_stockout_risk.py) | `[NEW]` | Full unit & integration test suite |
| [requirements.txt](file:///d:/__Projects/kage/chia.florist/ai-lab/requirements.txt) | `[MODIFY]` | Add `imbalanced-learn>=0.11.0` |
| `models/stockout_risk.json` | `[GENERATED]` | Trained XGBoost classifier checkpoint |
| `models/stockout_risk_eval.json` | `[GENERATED]` | Evaluation metrics: Precision, Recall, F1, ROC-AUC |
| `models/stockout_risk_importance.png` | `[GENERATED]` | Feature importance bar chart |

---

## 6. Verification Plan

### Automated Tests
```powershell
pytest tests/test_stockout_risk.py -v
```

### Manual Verification
1. Run `python train_stockout.py --config configs/stockout_risk.yaml --dry-run` — should exit cleanly.
2. Run `python train_stockout.py --config configs/stockout_risk.yaml` — should train and log Precision / Recall / F1 / ROC-AUC.
3. Inspect `models/stockout_risk_eval.json` to confirm metrics meet targets in §4.
4. Open `models/stockout_risk_importance.png` — `reorder_urgency_ratio` and `estimated_days_to_stockout` should rank as the top two features.
