import os
import json
import pytest
import pandas as pd
import numpy as np
import xgboost as xgb

from src.demand_forecasting_trainer import DemandForecastingTrainer

@pytest.fixture
def dummy_config(tmp_path):
    csv_path = os.path.join(tmp_path, "dummy_demand.csv")
    model_dir = os.path.join(tmp_path, "models")

    # Create dummy DataFrame with 100 rows
    np.random.seed(42)
    df = pd.DataFrame({
        "gross_margin_pct": np.random.uniform(0.1, 0.5, 100),
        "view_count": np.random.randint(100, 1000, 100),
        "day_of_week": np.random.randint(0, 7, 100),
        "month": np.random.randint(1, 13, 100),
        "units_sold_lag_1": np.random.uniform(0, 10, 100),
        "units_sold_lag_7": np.random.uniform(0, 10, 100),
        "target_label": np.random.uniform(5, 30, 100)
    })
    df.to_csv(csv_path, index=False)

    return {
        "seed": 42,
        "data": {
            "processed_csv": csv_path,
            "target_column": "target_label",
            "train_val_split": 0.8
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
            "model_name": "test_demand_model.json",
            "evaluation_report_name": "test_eval.json",
            "feature_importance_name": "test_importance.png"
        }
    }


def test_data_loading_and_split(dummy_config):
    trainer = DemandForecastingTrainer(dummy_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    assert len(X_train) == 80
    assert len(X_val) == 20
    assert "target_label" not in feature_names
    assert len(feature_names) == 6


def test_chronological_split(dummy_config):
    trainer = DemandForecastingTrainer(dummy_config)
    X_train, y_train, X_val, y_val, _ = trainer.load_and_split_data()

    # Train index range must be strictly smaller than val index range
    assert max(X_train.index) < min(X_val.index)


def test_trainer_fit_and_evaluation(dummy_config):
    trainer = DemandForecastingTrainer(dummy_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    model = trainer.train(X_train, y_train, X_val, y_val)
    assert isinstance(model, xgb.XGBRegressor)

    metrics = trainer.evaluate(model, X_train, y_train, X_val, y_val)
    assert "val_rmse" in metrics
    assert "val_mae" in metrics
    assert "val_r2" in metrics
    assert metrics["val_rmse"] >= 0.0

    eval_json_path = os.path.join(dummy_config["output"]["model_dir"], "test_eval.json")
    assert os.path.exists(eval_json_path)


def test_model_serialization_and_reloading(dummy_config):
    trainer = DemandForecastingTrainer(dummy_config)
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    model = trainer.train(X_train, y_train, X_val, y_val)
    model_path = trainer.save_model(model)
    plot_path = trainer.save_feature_importance(model, feature_names)

    assert os.path.exists(model_path)
    assert os.path.exists(plot_path)

    # Reload model and test inference
    reloaded_model = xgb.XGBRegressor()
    reloaded_model.load_model(model_path)

    preds_original = model.predict(X_val)
    preds_reloaded = reloaded_model.predict(X_val)

    np.testing.assert_allclose(preds_original, preds_reloaded, rtol=1e-5)
