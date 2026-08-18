import os
import json
import logging
import joblib
from typing import Dict, Any, Tuple, Optional
import numpy as np
import pandas as pd
import xgboost as xgb
from sklearn.ensemble import IsolationForest
from sklearn.model_selection import train_test_split
from sklearn.metrics import precision_score, recall_score, f1_score, roc_auc_score, classification_report
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

try:
    from imblearn.over_sampling import SMOTE
    SMOTE_AVAILABLE = True
except ImportError:
    SMOTE_AVAILABLE = False

logger = logging.getLogger("anomaly_detection_trainer")

class AnomalyDetectionTrainer:
    """
    Two-stage Operational Anomaly Detection Trainer (Phase 2.3).
    Stage 1: Isolation Forest (unsupervised outlier detection)
    Stage 2: XGBoost Classifier (supervised imbalance-resilient classification)
    Consensus: AND-gate logic combining IF and XGB outputs.
    """
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.seed = config.get("seed", 42)

        data_cfg = config.get("data", {})
        self.csv_path = data_cfg.get("processed_csv", "data/processed/anomaly_detection_features.csv")
        self.target_col = data_cfg.get("target_column", "target_label")
        self.train_val_split = data_cfg.get("train_val_split", 0.8)
        self.use_smote = data_cfg.get("use_smote", True)

        self.if_cfg = config.get("isolation_forest", {})
        self.xgb_cfg = config.get("xgb_classifier", {})

        self.output_cfg = config.get("output", {})
        self.model_dir = self.output_cfg.get("model_dir", "models/anomaly_detector")
        os.makedirs(self.model_dir, exist_ok=True)

    def load_and_split_data(self) -> Tuple[pd.DataFrame, pd.Series, pd.DataFrame, pd.Series, list]:
        if not os.path.exists(self.csv_path):
            raise FileNotFoundError(f"Processed dataset CSV not found at '{self.csv_path}'")

        df = pd.read_csv(self.csv_path)
        if self.target_col not in df.columns:
            raise KeyError(f"Target column '{self.target_col}' not found in CSV.")

        feature_cols = [c for c in df.columns if c != self.target_col]
        X = df[feature_cols]
        y = df[self.target_col].astype(int)

        X_train, X_val, y_train, y_val = train_test_split(
            X, y,
            test_size=(1.0 - self.train_val_split),
            random_state=self.seed,
            stratify=y
        )

        logger.info(
            f"Loaded dataset from '{self.csv_path}'. "
            f"Train: {len(X_train)} (Pos: {sum(y_train)}, Neg: {len(y_train)-sum(y_train)}), "
            f"Val: {len(X_val)} (Pos: {sum(y_val)}, Neg: {len(y_val)-sum(y_val)})"
        )
        return X_train, y_train, X_val, y_val, feature_cols

    def apply_smote(self, X_train: pd.DataFrame, y_train: pd.Series) -> Tuple[pd.DataFrame, pd.Series]:
        if not self.use_smote:
            logger.info("SMOTE disabled in configuration.")
            return X_train, y_train

        if not SMOTE_AVAILABLE:
            logger.warning("imbalanced-learn not installed; skipping SMOTE oversampling.")
            return X_train, y_train

        num_positives = sum(y_train)
        k_neighbors = min(5, num_positives - 1)
        logger.info(f"Applying SMOTE oversampling on XGBoost training split (k_neighbors={k_neighbors})...")
        smote = SMOTE(random_state=self.seed, k_neighbors=k_neighbors)
        X_resampled, y_resampled = smote.fit_resample(X_train, y_train)

        logger.info(
            f"Post-SMOTE train samples: {len(X_resampled)} "
            f"(Pos: {sum(y_resampled)}, Neg: {len(y_resampled)-sum(y_resampled)})"
        )
        return X_resampled, y_resampled

    def train_isolation_forest(self, X_train: pd.DataFrame) -> IsolationForest:
        logger.info("Training Stage 1: Isolation Forest (unsupervised outlier model)...")
        if_model = IsolationForest(
            n_estimators=self.if_cfg.get("n_estimators", 200),
            contamination=self.if_cfg.get("contamination", 0.06),
            max_samples=self.if_cfg.get("max_samples", "auto"),
            random_state=self.seed,
            n_jobs=-1
        )
        if_model.fit(X_train)
        logger.info("Isolation Forest training complete.")
        return if_model

    def train_xgb_classifier(
        self,
        X_train: pd.DataFrame,
        y_train: pd.Series,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> xgb.XGBClassifier:
        logger.info("Training Stage 2: XGBoost Classifier (supervised model)...")
        xgb_model = xgb.XGBClassifier(
            n_estimators=self.xgb_cfg.get("n_estimators", 300),
            max_depth=self.xgb_cfg.get("max_depth", 4),
            learning_rate=self.xgb_cfg.get("learning_rate", 0.05),
            subsample=self.xgb_cfg.get("subsample", 0.8),
            colsample_bytree=self.xgb_cfg.get("colsample_bytree", 0.8),
            scale_pos_weight=self.xgb_cfg.get("scale_pos_weight", 17),
            random_state=self.seed,
            n_jobs=-1,
            early_stopping_rounds=self.xgb_cfg.get("early_stopping_rounds", 20),
            eval_metric=self.xgb_cfg.get("eval_metric", "logloss")
        )

        xgb_model.fit(
            X_train,
            y_train,
            eval_set=[(X_train, y_train), (X_val, y_val)],
            verbose=False
        )

        best_iteration = getattr(xgb_model, "best_iteration", xgb_model.n_estimators)
        logger.info(f"XGBoost Classifier training complete. Best iteration: {best_iteration}")
        return xgb_model

    def evaluate(
        self,
        if_model: IsolationForest,
        xgb_model: xgb.XGBClassifier,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> Dict[str, Any]:
        # 1. Isolation Forest prediction: 1 (inlier) -> 0, -1 (outlier) -> 1
        if_raw = if_model.predict(X_val)
        if_preds = np.where(if_raw == -1, 1, 0)
        if_scores = -if_model.score_samples(X_val) # Higher = more anomalous

        if_prec = float(precision_score(y_val, if_preds, zero_division=0))
        if_rec = float(recall_score(y_val, if_preds, zero_division=0))
        if_f1 = float(f1_score(y_val, if_preds, zero_division=0))
        if_auc = float(roc_auc_score(y_val, if_scores))

        # 2. XGBoost Classifier prediction
        xgb_preds = xgb_model.predict(X_val)
        xgb_probs = xgb_model.predict_proba(X_val)[:, 1]

        xgb_prec = float(precision_score(y_val, xgb_preds, zero_division=0))
        xgb_rec = float(recall_score(y_val, xgb_preds, zero_division=0))
        xgb_f1 = float(f1_score(y_val, xgb_preds, zero_division=0))
        xgb_auc = float(roc_auc_score(y_val, xgb_probs))

        # 3. Consensus prediction (AND-gate logic)
        consensus_preds = np.logical_and(if_preds == 1, xgb_preds == 1).astype(int)
        con_prec = float(precision_score(y_val, consensus_preds, zero_division=0))
        con_rec = float(recall_score(y_val, consensus_preds, zero_division=0))
        con_f1 = float(f1_score(y_val, consensus_preds, zero_division=0))

        metrics = {
            "isolation_forest": {
                "precision": round(if_prec, 4),
                "recall": round(if_rec, 4),
                "f1_score": round(if_f1, 4),
                "roc_auc": round(if_auc, 4)
            },
            "xgb_classifier": {
                "precision": round(xgb_prec, 4),
                "recall": round(xgb_rec, 4),
                "f1_score": round(xgb_f1, 4),
                "roc_auc": round(xgb_auc, 4),
                "best_iteration": int(getattr(xgb_model, "best_iteration", 0))
            },
            "consensus": {
                "precision": round(con_prec, 4),
                "recall": round(con_rec, 4),
                "f1_score": round(con_f1, 4)
            }
        }

        logger.info("Evaluation Results:")
        logger.info(f"  • Isolation Forest:  Prec={if_prec:.4f}, Rec={if_rec:.4f} (Target Rec > 0.60), AUC={if_auc:.4f}")
        logger.info(f"  • XGB Classifier:    Prec={xgb_prec:.4f}, Rec={xgb_rec:.4f} (Target Rec > 0.75), AUC={xgb_auc:.4f} (Target AUC > 0.80)")
        logger.info(f"  • Combined Consensus: Prec={con_prec:.4f} (Target Prec > 0.70), Rec={con_rec:.4f}, F1={con_f1:.4f}")

        # Save eval report
        report_path = os.path.join(self.model_dir, self.output_cfg.get("evaluation_report_name", "anomaly_eval.json"))
        with open(report_path, "w") as f:
            json.dump(metrics, f, indent=2)
        logger.info(f"Saved evaluation report to '{report_path}'")

        return metrics

    def save_feature_importance(self, xgb_model: xgb.XGBClassifier, feature_names: list) -> str:
        importance = xgb_model.feature_importances_
        sorted_idx = np.argsort(importance)

        plt.figure(figsize=(10, 6))
        plt.barh(range(len(sorted_idx)), importance[sorted_idx], align="center", color="#e74c3c")
        plt.yticks(range(len(sorted_idx)), [feature_names[i] for i in sorted_idx])
        plt.xlabel("Feature Importance (Gain)")
        plt.title("XGBoost Operational Anomaly Classifier — Feature Importance")
        plt.tight_layout()

        plot_path = os.path.join(self.model_dir, self.output_cfg.get("feature_importance_name", "anomaly_importance.png"))
        plt.savefig(plot_path, dpi=150)
        plt.close()
        logger.info(f"Saved feature importance plot to '{plot_path}'")
        return plot_path

    def save_models(self, if_model: IsolationForest, xgb_model: xgb.XGBClassifier) -> Tuple[str, str]:
        if_name = self.output_cfg.get("isolation_forest_name", "isolation_forest.pkl")
        xgb_name = self.output_cfg.get("xgb_classifier_name", "xgb_classifier.json")

        if_path = os.path.join(self.model_dir, if_name)
        xgb_path = os.path.join(self.model_dir, xgb_name)

        joblib.dump(if_model, if_path)
        xgb_model.save_model(xgb_path)

        logger.info(f"Saved Stage 1 Isolation Forest checkpoint to '{if_path}'")
        logger.info(f"Saved Stage 2 XGBoost Classifier checkpoint to '{xgb_path}'")
        return if_path, xgb_path
