# 🚨 Phase 2.3 Implementation Plan: Operational Anomaly Detection Model

> **Module**: `ai-lab`  
> **Phase**: 2.3 — Anomaly Detection Model  
> **Date**: 2026-08-11  
> **Predecessor**: [PHASE_2_2_IMPLEMENTATION_PLAN.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/PHASE_2_2_IMPLEMENTATION_PLAN.md)  
> **Reference**: [MODEL_TRAINING_REPORT.md](file:///d:/__Projects/kage/chia.florist/ai-lab/docs/MODEL_TRAINING_REPORT.md)  
> **Status**: Proposed / Pending Review  

---

## Executive Summary

Phase 2.3 trains the **Operational Anomaly Detection Model** — the third model in the Phase 2 stack. It identifies operational failure signals in payment and fulfillment flows in real time, giving staff early warning of events like payment failure spikes, abnormally slow settlement, or manual transfer abuse.

The model is trained on the **1,508-sample, 4-feature** matrix produced by Phase 1 (`data/processed/anomaly_detection_features.csv`). Key facts from the dataset:

- **1,421 normal samples**, **87 anomalous samples** — a **~17:1 ratio** (similar imbalance challenge as Phase 2.2 but with a larger, richer dataset).
- Features are already **standardized** (mean ≈ 0, std ≈ 1) by the Phase 1 pipeline's `StandardScaler`.
- Ground-truth anomaly labels exist (derived from payment failure status and time-to-pay threshold), enabling **supervised** evaluation.

---

## 1. Key Design Decision: Two-Stage Architecture

Phase 2.3 uses a **two-stage approach** that mirrors how anomaly detection works in production:

| Stage | Algorithm | Role |
|---|---|---|
| **Stage 1** | **Isolation Forest** (unsupervised) | Trains with no label information. Produces an anomaly score for every sample. Works purely on the structure of the feature space. |
| **Stage 2** | **XGBoost Classifier** (supervised) | Trains on Phase 1 labels. Uses SMOTE to handle 17:1 imbalance. Produces calibrated probability scores per sample. |

Both models are trained and serialized independently. Their outputs are combined into a **consensus anomaly score** reported to staff: a sample is flagged as anomalous only when **both** models agree (AND-gate logic), minimizing false positives.

> [!NOTE]
> This two-stage design is intentional: when real `service-core` data becomes available, Isolation Forest can run continuously on unlabeled streams, while the XGBoost classifier acts as a supervised refinement layer trained on manually confirmed anomaly labels.

---

## 2. Architecture Overview

```
 ┌─────────────────────────────────────────────────────────────────────────────┐
 │              Phase 1 Output (ai-lab/data/processed/)                        │
 │  anomaly_detection_features.csv — 1,508 samples × 4 features               │
 │  Class Distribution: {Normal: 1421, Anomaly: 87}  — 17:1 imbalance         │
 │                                                                             │
 │  Features (all standardized): amount, time_to_pay_sec,                     │
 │                                is_failed_status, is_manual_transfer         │
 │  Target: is_anomaly_sample (binary 0 / 1)                                  │
 └─────────────────────────────────────────────────────────────────────────────┘
                                      │
                        ┌─────────────┴─────────────┐
                        ▼                           ▼
          ┌─────────────────────────┐  ┌──────────────────────────────┐
          │  Stage 1: IsolationForest│  │  Stage 2: XGBClassifier      │
          │  (unsupervised, no labels│  │  (supervised + SMOTE)        │
          │  needed for training)    │  │                              │
          │  → anomaly_score (float) │  │  → anomaly_prob (0.0–1.0)   │
          └──────────┬──────────────┘  └────────────┬─────────────────┘
                     │                              │
                     └───────────────┬──────────────┘
                                     │ Consensus Report
                                     ▼
                      ┌──────────────────────────────┐
                      │  models/anomaly_detector/     │
                      │  ├── isolation_forest.pkl     │
                      │  ├── xgb_classifier.json      │
                      │  ├── anomaly_eval.json         │
                      │  └── anomaly_importance.png    │
                      └──────────────────────────────┘
```

---

## 3. Planned Components & Code Changes

### 3.1 New Configuration

#### [NEW] [configs/anomaly_detection.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/anomaly_detection.yaml)

```yaml
experiment_name: "anomaly_detection_v1"
seed: 42

data:
  processed_csv: "data/processed/anomaly_detection_features.csv"
  target_column: "target_label"
  train_val_split: 0.8
  use_smote: true

isolation_forest:
  n_estimators: 200
  contamination: 0.06         # Approximate anomaly rate in dataset (~5.8%)
  max_samples: "auto"

xgb_classifier:
  n_estimators: 300
  max_depth: 4
  learning_rate: 0.05
  subsample: 0.8
  colsample_bytree: 0.8
  scale_pos_weight: 17
  early_stopping_rounds: 20
  eval_metric: "logloss"

output:
  model_dir: "models/anomaly_detector"
  isolation_forest_name: "isolation_forest.pkl"
  xgb_classifier_name: "xgb_classifier.json"
  evaluation_report_name: "anomaly_eval.json"
  feature_importance_name: "anomaly_importance.png"
```

---

### 3.2 Core Training Module

#### [NEW] [src/anomaly_detection_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/anomaly_detection_trainer.py)

Class `AnomalyDetectionTrainer` with:
- `load_and_split_data()` — Reads Phase 1 CSV, performs a **stratified** train/val split preserving the 17:1 class ratio.
- `apply_smote(X_train, y_train)` — Reuses the same SMOTE pattern from Phase 2.2 to balance the XGBoost training split.
- `train_isolation_forest(X_train)` — Fits `sklearn.ensemble.IsolationForest` on the unlabeled feature matrix (no `y_train` used). Serializes with `joblib.dump()` to `.pkl`.
- `train_xgb_classifier(X_train_res, y_train_res, X_val, y_val)` — Fits `xgboost.XGBClassifier` with early stopping on validation logloss.
- `evaluate(if_model, xgb_model, X_val, y_val)` — Produces independent metrics for both models, then computes the **consensus score** (flagged anomalous when IF predicts -1 AND XGB predicts 1). Reports Precision, Recall, F1, ROC-AUC for all three.
- `save_feature_importance(xgb_model, feature_names)` — Feature importance chart from the XGBoost classifier.
- `save_models(if_model, xgb_model)` — Serializes both checkpoints.

---

### 3.3 CLI Training Script

#### [NEW] [train_anomaly.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_anomaly.py)

CLI entry point:

```bash
# Full training run (both stages)
python train_anomaly.py --config configs/anomaly_detection.yaml

# Dry run (data loading and split validation only)
python train_anomaly.py --config configs/anomaly_detection.yaml --dry-run
```

---

### 3.4 Unit Tests

#### [NEW] [tests/test_anomaly_detection.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_anomaly_detection.py)

Tests to cover:
- `test_data_loading_and_split` — CSV loads, 4 features present, stratified split applied.
- `test_isolation_forest_training` — IF trains without error, `.pkl` checkpoint serializes and reloads.
- `test_xgb_classifier_training` — XGB classifier fits, SMOTE applied, model serializes.
- `test_evaluation_metrics` — Both models produce valid Precision, Recall, F1, ROC-AUC floats.
- `test_consensus_scoring` — Combined flag is generated correctly for each sample.

---

## 4. Evaluation Targets

Because anomaly detection is operationally critical (missed anomalies = undetected payment failures), **Recall is the primary metric**.

| Stage | Metric | Target | Notes |
|---|---|---|---|
| Isolation Forest | Recall | > 0.60 | Unsupervised — lower recall expected; acts as a low-false-positive gate |
| XGB Classifier | Recall | > 0.75 | Supervised — higher recall expected thanks to SMOTE |
| XGB Classifier | ROC-AUC | > 0.80 | Strong discriminant quality expected on 1,508 samples |
| **Consensus (both)** | **Precision** | **> 0.70** | Combined output should have higher precision (fewer false alerts) |

---

## 5. File Summary

| File | Status | Description |
|---|---|---|
| [configs/anomaly_detection.yaml](file:///d:/__Projects/kage/chia.florist/ai-lab/configs/anomaly_detection.yaml) | `[NEW]` | Two-stage model hyperparameters and I/O config |
| [src/anomaly_detection_trainer.py](file:///d:/__Projects/kage/chia.florist/ai-lab/src/anomaly_detection_trainer.py) | `[NEW]` | IF + XGB trainer class with consensus scoring |
| [train_anomaly.py](file:///d:/__Projects/kage/chia.florist/ai-lab/train_anomaly.py) | `[NEW]` | CLI entry script |
| [tests/test_anomaly_detection.py](file:///d:/__Projects/kage/chia.florist/ai-lab/tests/test_anomaly_detection.py) | `[NEW]` | Full unit & integration test suite |
| `models/anomaly_detector/isolation_forest.pkl` | `[GENERATED]` | Trained Isolation Forest checkpoint |
| `models/anomaly_detector/xgb_classifier.json` | `[GENERATED]` | Trained XGBoost classifier checkpoint |
| `models/anomaly_detector/anomaly_eval.json` | `[GENERATED]` | Evaluation metrics for IF, XGB, and consensus |
| `models/anomaly_detector/anomaly_importance.png` | `[GENERATED]` | XGBoost feature importance bar chart |

---

## 6. Verification Plan

### Automated Tests
```powershell
pytest tests/test_anomaly_detection.py -v
```

### Manual Verification
1. Run `python train_anomaly.py --config configs/anomaly_detection.yaml --dry-run` — should exit cleanly.
2. Run `python train_anomaly.py --config configs/anomaly_detection.yaml` — should train both stages and log metrics for IF, XGB, and consensus.
3. Inspect `models/anomaly_detector/anomaly_eval.json` — confirm metrics meet targets in §4.
4. Open `models/anomaly_detector/anomaly_importance.png` — `time_to_pay_sec` and `is_failed_status` should rank highest.
