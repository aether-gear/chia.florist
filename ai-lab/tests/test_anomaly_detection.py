import os
import json
import pytest
import joblib
import pandas as pd
import numpy as np
import xgboost as xgb
from sklearn.ensemble import IsolationForest

from src.anomaly_detection_trainer import AnomalyDetectionTrainer

@pytest.fixture
def dummy_anomaly_config(tmp_path):
    csv_path = os.path.join(tmp_path, "dummy_anomaly.csv")
    model_dir = os.path.join(tmp_path, "models", "anomaly_detector")

    # Create dummy DataFrame with 200 normal samples and 20 anomalous samples
    np.random.seed(42)
    normals = pd.DataFrame({
        "amount": np.random.normal(0, 1, 200),
        "time_to_pay_sec": np.random.normal(0, 0.5, 200),
        "is_failed_status": np.zeros(200),
        "is_manual_transfer": np.random.choice([0, 1], 200),
        "target_label": np.zeros(200)
    })
    anomalies = pd.DataFrame({
        "amount": np.random.normal(2, 1, 20),
        "time_to_pay_sec": np.random.normal(3, 1, 20),
        "is_failed_status": np.ones(20),
        "is_manual_transfer": np.random.choice([0, 1], 20),
        "target_label": np.ones(20)
    })

    df = pd.concat([normals, anomalies], ignore_index=True).sample(frac=1.0, random_state=42).reset_index(drop=True)
    df.to_csv(csv_path, index=False)

    return {
        "seed": 42,
        "data": {
            "processed_csv": csv_path,
            "target_column": "target_label",
            "train_val_split": 0.8,
            "use_smote": True
        },
        "isolation_forest": {
            "n_estimators": 50,
            "contamination": 0.09,
            "max_samples": "auto"
        },
        "xgb_classifier": {
            "n_estimators": 20,
            "max_depth": 3,
            "learning_rate": 0.1,
            "scale_pos_weight": 10,
            "early_stopping_rounds": 5,
            "eval_metric": "logloss"
        },
        "output": {
            "model_dir": model_dir,
            "isolation_forest_name": "test_if.pkl",
            "xgb_classifier_name": "test_xgb.json",
            "evaluation_report_name": "test_eval.json",
            "feature_importance_name": "test_importance.png"
        }
    }


def test_data_loading_and_split(dummy_anomaly_config):
    trainer = AnomalyDetectionTrainer(dummy_anomaly_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    assert len(X_train) == 176
    assert len(X_val) == 44
    assert sum(y_train) > 0 # Stratified split
    assert sum(y_val) > 0
    assert len(feature_names) == 4


def test_isolation_forest_training(dummy_anomaly_config):
    trainer = AnomalyDetectionTrainer(dummy_anomaly_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()

    if_model = trainer.train_isolation_forest(X_train)
    assert isinstance(if_model, IsolationForest)

    preds = if_model.predict(X_val)
    assert len(preds) == len(X_val)
    assert set(np.unique(preds)).issubset({-1, 1})


def test_xgb_classifier_training(dummy_anomaly_config):
    trainer = AnomalyDetectionTrainer(dummy_anomaly_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()
    X_res, y_res = trainer.apply_smote(X_train, y_train)

    xgb_model = trainer.train_xgb_classifier(X_res, y_res, X_val, y_val)
    assert isinstance(xgb_model, xgb.XGBClassifier)

    probs = xgb_model.predict_proba(X_val)[:, 1]
    assert len(probs) == len(X_val)
    assert np.all((probs >= 0.0) & (probs <= 1.0))


def test_evaluation_and_consensus(dummy_anomaly_config):
    trainer = AnomalyDetectionTrainer(dummy_anomaly_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()
    X_res, y_res = trainer.apply_smote(X_train, y_train)

    if_model = trainer.train_isolation_forest(X_train)
    xgb_model = trainer.train_xgb_classifier(X_res, y_res, X_val, y_val)

    metrics = trainer.evaluate(if_model, xgb_model, X_val, y_val)

    assert "isolation_forest" in metrics
    assert "xgb_classifier" in metrics
    assert "consensus" in metrics

    assert "recall" in metrics["isolation_forest"]
    assert "roc_auc" in metrics["xgb_classifier"]
    assert "precision" in metrics["consensus"]


def test_model_serialization(dummy_anomaly_config):
    trainer = AnomalyDetectionTrainer(dummy_anomaly_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()
    X_res, y_res = trainer.apply_smote(X_train, y_train)

    if_model = trainer.train_isolation_forest(X_train)
    xgb_model = trainer.train_xgb_classifier(X_res, y_res, X_val, y_val)

    if_path, xgb_path = trainer.save_models(if_model, xgb_model)
    plot_path = trainer.save_feature_importance(xgb_model, feature_names)

    assert os.path.exists(if_path)
    assert os.path.exists(xgb_path)
    assert os.path.exists(plot_path)

    # Reload checkpoints
    reloaded_if = joblib.load(if_path)
    reloaded_xgb = xgb.XGBClassifier()
    reloaded_xgb.load_model(xgb_path)

    np.testing.assert_array_equal(if_model.predict(X_val), reloaded_if.predict(X_val))
    np.testing.assert_allclose(xgb_model.predict_proba(X_val), reloaded_xgb.predict_proba(X_val), rtol=1e-5)
