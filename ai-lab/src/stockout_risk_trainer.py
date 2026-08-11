import os
import json
import logging
from typing import Dict, Any, Tuple, Optional
import numpy as np
import pandas as pd
import xgboost as xgb
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

logger = logging.getLogger("stockout_risk_trainer")

class StockoutRiskTrainer:
    """
    Trainer and evaluator for XGBoost-based Inventory Stockout Risk Classifier (Phase 2.2).
    Includes SMOTE oversampling and stratified splitting to handle class imbalance.
    """
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.seed = config.get("seed", 42)
        
        data_cfg = config.get("data", {})
        self.csv_path = data_cfg.get("processed_csv", "data/processed/stockout_risk_features.csv")
        self.target_col = data_cfg.get("target_column", "target_label")
        self.train_val_split = data_cfg.get("train_val_split", 0.8)
        self.use_smote = data_cfg.get("use_smote", True)

        self.model_cfg = config.get("model", {})
        self.output_cfg = config.get("output", {})
        self.model_dir = self.output_cfg.get("model_dir", "models")
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

        # Stratified train/val split to guarantee minority class representation in validation set
        try:
            X_train, X_val, y_train, y_val = train_test_split(
                X, y,
                test_size=(1.0 - self.train_val_split),
                random_state=self.seed,
                stratify=y
            )
        except ValueError:
            # Fallback to simple split if class counts are too low for stratification
            logger.warning("Class counts too low for stratified split; falling back to unstratified split.")
            X_train, X_val, y_train, y_val = train_test_split(
                X, y,
                test_size=(1.0 - self.train_val_split),
                random_state=self.seed
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
        if num_positives < 2:
            logger.warning(f"Not enough positive samples ({num_positives}) for SMOTE; skipping.")
            return X_train, y_train

        k_neighbors = min(3, num_positives - 1)
        logger.info(f"Applying SMOTE oversampling on training split (k_neighbors={k_neighbors})...")
        smote = SMOTE(random_state=self.seed, k_neighbors=k_neighbors)
        X_resampled, y_resampled = smote.fit_resample(X_train, y_train)

        logger.info(
            f"Post-SMOTE train samples: {len(X_resampled)} "
            f"(Pos: {sum(y_resampled)}, Neg: {len(y_resampled)-sum(y_resampled)})"
        )
        return X_resampled, y_resampled

    def train(
        self,
        X_train: pd.DataFrame,
        y_train: pd.Series,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> xgb.XGBClassifier:
        logger.info("Initializing XGBClassifier model...")

        model = xgb.XGBClassifier(
            n_estimators=self.model_cfg.get("n_estimators", 300),
            max_depth=self.model_cfg.get("max_depth", 4),
            learning_rate=self.model_cfg.get("learning_rate", 0.05),
            subsample=self.model_cfg.get("subsample", 0.8),
            colsample_bytree=self.model_cfg.get("colsample_bytree", 0.8),
            scale_pos_weight=self.model_cfg.get("scale_pos_weight", 1.0),
            random_state=self.seed,
            n_jobs=-1,
            early_stopping_rounds=self.model_cfg.get("early_stopping_rounds", 20),
            eval_metric=self.model_cfg.get("eval_metric", "logloss")
        )

        logger.info("Fitting model with early stopping on validation set...")
        model.fit(
            X_train,
            y_train,
            eval_set=[(X_train, y_train), (X_val, y_val)],
            verbose=False
        )

        best_iteration = getattr(model, "best_iteration", model.n_estimators)
        logger.info(f"Training finished. Best iteration: {best_iteration}")
        return model

    def evaluate(
        self,
        model: xgb.XGBClassifier,
        X_train: pd.DataFrame,
        y_train: pd.Series,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> Dict[str, Any]:
        val_preds = model.predict(X_val)
        val_probs = model.predict_proba(X_val)[:, 1] if hasattr(model, "predict_proba") else val_preds

        # Metrics calculation
        prec = float(precision_score(y_val, val_preds, zero_division=0))
        rec = float(recall_score(y_val, val_preds, zero_division=0))
        f1 = float(f1_score(y_val, val_preds, zero_division=0))
        
        try:
            auc = float(roc_auc_score(y_val, val_probs))
        except ValueError:
            auc = 0.5 # Single class fallback in val set

        cls_report = classification_report(y_val, val_preds, output_dict=True, zero_division=0)

        metrics = {
            "val_precision": round(prec, 4),
            "val_recall": round(rec, 4),
            "val_f1": round(f1, 4),
            "val_roc_auc": round(auc, 4),
            "best_iteration": int(getattr(model, "best_iteration", 0)),
            "classification_report": cls_report
        }

        logger.info(f"Evaluation Results:")
        logger.info(f"  • Val Precision: {prec:.4f} (Target > 0.50)")
        logger.info(f"  • Val Recall:    {rec:.4f} (Target > 0.70)")
        logger.info(f"  • Val F1 Score:  {f1:.4f} (Target > 0.60)")
        logger.info(f"  • Val ROC-AUC:   {auc:.4f} (Target > 0.75)")

        # Save eval JSON
        report_path = os.path.join(self.model_dir, self.output_cfg.get("evaluation_report_name", "stockout_risk_eval.json"))
        with open(report_path, "w") as f:
            json.dump(metrics, f, indent=2)
        logger.info(f"Saved evaluation metrics to '{report_path}'")

        return metrics

    def save_feature_importance(self, model: xgb.XGBClassifier, feature_names: list) -> str:
        importance = model.feature_importances_
        sorted_idx = np.argsort(importance)

        plt.figure(figsize=(10, 6))
        plt.barh(range(len(sorted_idx)), importance[sorted_idx], align="center", color="#c0392b")
        plt.yticks(range(len(sorted_idx)), [feature_names[i] for i in sorted_idx])
        plt.xlabel("Feature Importance (Gain)")
        plt.title("XGBoost Stockout Risk Classifier — Feature Importance")
        plt.tight_layout()

        plot_path = os.path.join(self.model_dir, self.output_cfg.get("feature_importance_name", "stockout_risk_importance.png"))
        plt.savefig(plot_path, dpi=150)
        plt.close()
        logger.info(f"Saved feature importance plot to '{plot_path}'")
        return plot_path

    def save_model(self, model: xgb.XGBClassifier) -> str:
        model_name = self.output_cfg.get("model_name", "stockout_risk.json")
        model_path = os.path.join(self.model_dir, model_name)
        model.save_model(model_path)
        logger.info(f"Saved native XGBoost classifier checkpoint to '{model_path}'")
        return model_path
