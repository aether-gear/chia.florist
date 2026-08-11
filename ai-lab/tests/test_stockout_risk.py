import os
import json
import pytest
import pandas as pd
import numpy as np
import xgboost as xgb

from src.stockout_risk_trainer import StockoutRiskTrainer

@pytest.fixture
def dummy_stockout_config(tmp_path):
    csv_path = os.path.join(tmp_path, "dummy_stockout.csv")
    model_dir = os.path.join(tmp_path, "models")
    
    # Create imbalanced DataFrame (90 negatives, 10 positives)
    np.random.seed(42)
    negs = pd.DataFrame({
        "stock": np.random.uniform(20, 100, 90),
        "reserved_stock": np.random.uniform(0, 5, 90),
        "stock_burn_rate_7d": np.random.uniform(0.1, 2.0, 90),
        "supplier_lead_time_days": np.random.choice([2, 3, 5], 90),
        "estimated_days_to_stockout": np.random.uniform(10, 50, 90),
        "reorder_urgency_ratio": np.random.uniform(0.01, 0.3, 90),
        "target_label": np.zeros(90)
    })
    pos = pd.DataFrame({
        "stock": np.random.uniform(0, 5, 10),
        "reserved_stock": np.random.uniform(1, 5, 10),
        "stock_burn_rate_7d": np.random.uniform(3.0, 10.0, 10),
        "supplier_lead_time_days": np.random.choice([3, 5, 7], 10),
        "estimated_days_to_stockout": np.random.uniform(0.1, 2.0, 10),
        "reorder_urgency_ratio": np.random.uniform(1.5, 5.0, 10),
        "target_label": np.ones(10)
    })

    df = pd.concat([negs, pos], ignore_index=True).sample(frac=1.0, random_state=42).reset_index(drop=True)
    df.to_csv(csv_path, index=False)

    return {
        "seed": 42,
        "data": {
            "processed_csv": csv_path,
            "target_column": "target_label",
            "train_val_split": 0.8,
            "use_smote": True
        },
        "model": {
            "n_estimators": 20,
            "max_depth": 3,
            "learning_rate": 0.1,
            "scale_pos_weight": 9.0,
            "early_stopping_rounds": 5,
            "eval_metric": "logloss"
        },
        "output": {
            "model_dir": model_dir,
            "model_name": "test_stockout_model.json",
            "evaluation_report_name": "test_stockout_eval.json",
            "feature_importance_name": "test_stockout_importance.png"
        }
    }


def test_data_loading_and_stratified_split(dummy_stockout_config):
    trainer = StockoutRiskTrainer(dummy_stockout_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    assert len(X_train) == 80
    assert len(X_val) == 20
    assert sum(y_train) > 0 # Stratified split ensures positives in train
    assert sum(y_val) > 0   # Stratified split ensures positives in val
    assert len(feature_names) == 6


def test_smote_balancing(dummy_stockout_config):
    trainer = StockoutRiskTrainer(dummy_stockout_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()

    X_res, y_res = trainer.apply_smote(X_train, y_train)

    assert len(X_res) >= len(X_train)
    # SMOTE should balance classes in training set
    assert sum(y_res) == (len(y_res) - sum(y_res))


def test_trainer_fit_and_evaluation(dummy_stockout_config):
    trainer = StockoutRiskTrainer(dummy_stockout_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()
    X_res, y_res = trainer.apply_smote(X_train, y_train)

    model = trainer.train(X_res, y_res, X_val, y_val)
    assert isinstance(model, xgb.XGBClassifier)

    metrics = trainer.evaluate(model, X_res, y_res, X_val, y_val)
    assert "val_precision" in metrics
    assert "val_recall" in metrics
    assert "val_f1" in metrics
    assert "val_roc_auc" in metrics

    eval_json_path = os.path.join(dummy_stockout_config["output"]["model_dir"], "test_stockout_eval.json")
    assert os.path.exists(eval_json_path)


def test_model_serialization_and_reloading(dummy_stockout_config):
    trainer = StockoutRiskTrainer(dummy_stockout_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()
    X_res, y_res = trainer.apply_smote(X_train, y_train)

    model = trainer.train(X_res, y_res, X_val, y_val)
    model_path = trainer.save_model(model)
    plot_path = trainer.save_feature_importance(model, feature_names)

    assert os.path.exists(model_path)
    assert os.path.exists(plot_path)

    # Reload model and test prediction consistency
    reloaded_model = xgb.XGBClassifier()
    reloaded_model.load_model(model_path)

    probs_original = model.predict_proba(X_val)
    probs_reloaded = reloaded_model.predict_proba(X_val)

    np.testing.assert_allclose(probs_original, probs_reloaded, rtol=1e-5)
