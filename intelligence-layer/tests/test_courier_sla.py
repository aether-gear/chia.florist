import os
import json
import pytest
import pandas as pd
import numpy as np
import xgboost as xgb

from src.feature_engineering import CourierSLAFeatureBuilder
from src.courier_sla_trainer import CourierSLATrainer

@pytest.fixture
def dummy_courier_config(tmp_path):
    csv_path = os.path.join(tmp_path, "dummy_courier_sla.csv")
    model_dir = os.path.join(tmp_path, "models")

    np.random.seed(42)
    df = pd.DataFrame({
        "courier_jne": np.random.choice([0.0, 1.0], 100),
        "courier_jnt": np.random.choice([0.0, 1.0], 100),
        "courier_sicepat": np.random.choice([0.0, 1.0], 100),
        "shipping_cost": np.random.uniform(10000, 30000, 100),
        "dispatch_day_of_week": np.random.randint(0, 7, 100).astype(float),
        "dispatch_hour": np.random.randint(8, 20, 100).astype(float),
        "dispatch_is_weekend": np.random.choice([0.0, 1.0], 100),
        "target_label": np.random.uniform(8.0, 48.0, 100) # delivery hours
    })
    df.to_csv(csv_path, index=False)

    return {
        "seed": 42,
        "data": {
            "processed_csv": csv_path,
            "target_column": "target_label",
            "train_val_split": 0.8,
            "courier_columns": ["courier_jne", "courier_jnt", "courier_sicepat"]
        },
        "model": {
            "n_estimators": 20,
            "max_depth": 3,
            "learning_rate": 0.1,
            "early_stopping_rounds": 5,
            "eval_metric": "rmse"
        },
        "output": {
            "model_dir": model_dir,
            "model_name": "test_courier_model.json",
            "evaluation_report_name": "test_courier_eval.json",
            "feature_importance_name": "test_courier_importance.png"
        }
    }


def test_courier_sla_feature_builder():
    raw_shipments = {
        "shipments": [
            {
                "courier_name": "JNE",
                "shipping_cost": 15000.0,
                "created_at": "2026-05-10 10:00:00",
                "delivered_at": "2026-05-11 12:00:00"
            },
            {
                "courier_name": "JNT",
                "shipping_cost": 20000.0,
                "created_at": "2026-05-10 14:00:00",
                "delivered_at": "2026-05-11 20:00:00"
            }
        ]
    }

    builder = CourierSLAFeatureBuilder()
    headers, X, y = builder.build_features(raw_shipments)

    assert len(headers) == 7
    assert X.shape == (2, 7)
    assert y.shape == (2,)
    assert y[0] == pytest.approx(26.0) # 26 hours duration
    assert y[1] == pytest.approx(30.0) # 30 hours duration


def test_data_loading_and_split(dummy_courier_config):
    trainer = CourierSLATrainer(dummy_courier_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    assert len(X_train) == 80
    assert len(X_val) == 20
    assert len(feature_names) == 7
    assert "target_label" not in feature_names


def test_trainer_fit_and_evaluation(dummy_courier_config):
    trainer = CourierSLATrainer(dummy_courier_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()

    model = trainer.train(X_train, y_train, X_val, y_val)
    assert isinstance(model, xgb.XGBRegressor)

    metrics = trainer.evaluate(model, X_train, y_train, X_val, y_val)

    assert "val_mae_hours" in metrics
    assert "val_rmse_hours" in metrics
    assert "val_r2" in metrics
    assert "reliability_scores" in metrics
    assert metrics["val_mae_hours"] >= 0.0

    eval_json_path = os.path.join(dummy_courier_config["output"]["model_dir"], "test_courier_eval.json")
    assert os.path.exists(eval_json_path)


def test_reliability_score_range(dummy_courier_config):
    trainer = CourierSLATrainer(dummy_courier_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()

    model = trainer.train(X_train, y_train, X_val, y_val)
    metrics = trainer.evaluate(model, X_train, y_train, X_val, y_val)

    for c_name, score in metrics["reliability_scores"].items():
        assert 0.0 <= score <= 100.0


def test_model_serialization_and_reloading(dummy_courier_config):
    trainer = CourierSLATrainer(dummy_courier_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    model = trainer.train(X_train, y_train, X_val, y_val)
    model_path = trainer.save_model(model)
    plot_path = trainer.save_feature_importance(model, feature_names)

    assert os.path.exists(model_path)
    assert os.path.exists(plot_path)

    # Reload model
    reloaded_model = xgb.XGBRegressor()
    reloaded_model.load_model(model_path)

    preds_original = model.predict(X_val)
    preds_reloaded = reloaded_model.predict(X_val)

    np.testing.assert_allclose(preds_original, preds_reloaded, rtol=1e-5)
