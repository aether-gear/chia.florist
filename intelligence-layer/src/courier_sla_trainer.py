import os
import json
import logging
from typing import Dict, Any, Tuple, Optional
import numpy as np
import pandas as pd
import xgboost as xgb
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

logger = logging.getLogger("courier_sla_trainer")

class CourierSLATrainer:
    """
    Trainer and evaluator for XGBoost Courier SLA & Delivery Duration Estimator (Phase 2.4).
    Computes overall MAE/RMSE/R² and per-courier reliability scores (0-100).
    """
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.seed = config.get("seed", 42)

        data_cfg = config.get("data", {})
        self.csv_path = data_cfg.get("processed_csv", "data/processed/courier_sla_features.csv")
        self.target_col = data_cfg.get("target_column", "target_label")
        self.train_val_split = data_cfg.get("train_val_split", 0.8)
        self.courier_cols = data_cfg.get("courier_columns", ["courier_jne", "courier_jnt", "courier_sicepat"])

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
        y = df[self.target_col]

        split_idx = int(len(df) * self.train_val_split)
        X_train, X_val = X.iloc[:split_idx], X.iloc[split_idx:]
        y_train, y_val = y.iloc[:split_idx], y.iloc[split_idx:]

        logger.info(
            f"Loaded dataset from '{self.csv_path}'. "
            f"Train samples: {len(X_train)}, Val samples: {len(X_val)}, Features: {len(feature_cols)}"
        )
        return X_train, y_train, X_val, y_val, feature_cols

    def train(
        self,
        X_train: pd.DataFrame,
        y_train: pd.Series,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> xgb.XGBRegressor:
        logger.info("Initializing XGBRegressor model...")

        model = xgb.XGBRegressor(
            n_estimators=self.model_cfg.get("n_estimators", 500),
            max_depth=self.model_cfg.get("max_depth", 5),
            learning_rate=self.model_cfg.get("learning_rate", 0.05),
            subsample=self.model_cfg.get("subsample", 0.8),
            colsample_bytree=self.model_cfg.get("colsample_bytree", 0.8),
            random_state=self.seed,
            n_jobs=-1,
            early_stopping_rounds=self.model_cfg.get("early_stopping_rounds", 30),
            eval_metric=self.model_cfg.get("eval_metric", "rmse")
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
        model: xgb.XGBRegressor,
        X_train: pd.DataFrame,
        y_train: pd.Series,
        X_val: pd.DataFrame,
        y_val: pd.Series
    ) -> Dict[str, Any]:
        val_preds = model.predict(X_val)

        val_rmse = float(np.sqrt(mean_squared_error(y_val, val_preds)))
        val_mae = float(mean_absolute_error(y_val, val_preds))
        val_r2 = float(r2_score(y_val, val_preds))

        # Per-courier breakdown
        per_courier_metrics = {}
        courier_maes = {}

        for col in self.courier_cols:
            if col in X_val.columns:
                courier_name = col.replace("courier_", "").upper()
                mask = X_val[col] == 1.0
                if mask.sum() > 0:
                    c_y = y_val[mask]
                    c_preds = val_preds[mask]
                    c_mae = float(mean_absolute_error(c_y, c_preds))
                    c_rmse = float(np.sqrt(mean_squared_error(c_y, c_preds)))
                    courier_maes[courier_name] = c_mae
                    per_courier_metrics[courier_name] = {
                        "samples": int(mask.sum()),
                        "mae_hours": round(c_mae, 4),
                        "rmse_hours": round(c_rmse, 4)
                    }

        # Calculate Courier Reliability Score (0 - 100)
        max_mae = max(courier_maes.values()) if courier_maes else 1.0
        reliability_scores = {}
        for c_name, c_mae in courier_maes.items():
            # Higher score = lower MAE relative to max error
            rel_score = round(100.0 * (1.0 - (c_mae / (max_mae * 1.5))), 2)
            reliability_scores[c_name] = max(0.0, min(100.0, rel_score))
            per_courier_metrics[c_name]["reliability_score"] = reliability_scores[c_name]

        metrics = {
            "val_mae_hours": round(val_mae, 4),
            "val_rmse_hours": round(val_rmse, 4),
            "val_r2": round(val_r2, 4),
            "best_iteration": int(getattr(model, "best_iteration", 0)),
            "per_courier": per_courier_metrics,
            "reliability_scores": reliability_scores
        }

        logger.info("Evaluation Results:")
        logger.info(f"  • Val MAE:  {val_mae:.4f} hrs (Target < 5.0 hrs)")
        logger.info(f"  • Val RMSE: {val_rmse:.4f} hrs (Target < 8.0 hrs)")
        logger.info(f"  • Val R²:   {val_r2:.4f} (Target > 0.50)")
        logger.info("Per-Courier Breakdown:")
        for c_name, score in reliability_scores.items():
            logger.info(f"  • {c_name}: MAE = {courier_maes[c_name]:.4f} hrs | Reliability Score = {score}/100")

        # Save eval JSON report
        report_path = os.path.join(self.model_dir, self.output_cfg.get("evaluation_report_name", "courier_sla_eval.json"))
        with open(report_path, "w") as f:
            json.dump(metrics, f, indent=2)
        logger.info(f"Saved evaluation report to '{report_path}'")

        return metrics

    def save_feature_importance(self, model: xgb.XGBRegressor, feature_names: list) -> str:
        importance = model.feature_importances_
        sorted_idx = np.argsort(importance)

        plt.figure(figsize=(10, 6))
        plt.barh(range(len(sorted_idx)), importance[sorted_idx], align="center", color="#27ae60")
        plt.yticks(range(len(sorted_idx)), [feature_names[i] for i in sorted_idx])
        plt.xlabel("Feature Importance (Gain)")
        plt.title("XGBoost Courier SLA & Delivery Duration — Feature Importance")
        plt.tight_layout()

        plot_path = os.path.join(self.model_dir, self.output_cfg.get("feature_importance_name", "courier_sla_importance.png"))
        plt.savefig(plot_path, dpi=150)
        plt.close()
        logger.info(f"Saved feature importance plot to '{plot_path}'")
        return plot_path

    def save_model(self, model: xgb.XGBRegressor) -> str:
        model_name = self.output_cfg.get("model_name", "courier_sla.json")
        model_path = os.path.join(self.model_dir, model_name)
        model.save_model(model_path)
        logger.info(f"Saved native XGBoost model checkpoint to '{model_path}'")
        return model_path
