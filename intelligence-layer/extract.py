import argparse
import os
import sys
import csv
import logging
import yaml
import torch
import numpy as np

from src.utils import setup_logging, set_seed
from src.extractor import get_extractor
from src.feature_engineering import (
    TimeSeriesFeatureBuilder,
    InventoryStockoutFeatureBuilder,
    OperationalAnomalyFeatureBuilder,
    CourierSLAFeatureBuilder,
    StandardScaler
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="CLI utility to execute Phase 1 feature extraction and processing pipeline."
    )
    parser.add_argument(
        "--config",
        type=str,
        default=os.path.join("configs", "feature_extraction.yaml"),
        help="Path to YAML configuration file."
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Run feature extraction in memory without saving artifacts to disk."
    )
    parser.add_argument(
        "--output-dir",
        type=str,
        default=None,
        help="Override output directory for processed features."
    )
    return parser.parse_args()

def load_config(config_path: str) -> dict:
    if not os.path.exists(config_path):
        print(f"Error: Configuration file not found at '{config_path}'")
        sys.exit(1)
    with open(config_path, "r") as f:
        try:
            return yaml.safe_load(f)
        except yaml.YAMLError as e:
            print(f"Error parsing YAML config: {e}")
            sys.exit(1)

def save_csv(filepath: str, headers: list, X: np.ndarray, y: np.ndarray):
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, "w", newline="", encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow(headers + ["target_label"])
        for i in range(len(X)):
            row = list(X[i]) + [float(y[i])]
            writer.writerow(row)

def save_pt(filepath: str, X: np.ndarray, y: np.ndarray):
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    torch.save({
        "X": torch.tensor(X, dtype=torch.float32),
        "y": torch.tensor(y, dtype=torch.float32)
    }, filepath)

def main():
    args = parse_args()
    config = load_config(args.config)

    setup_logging()
    logger = logging.getLogger("extract_pipeline")
    logger.info("Initializing Phase 1 Feature Extraction Pipeline...")

    set_seed(config.get("seed", 42))

    data_cfg = config.get("data", {})
    processed_dir = args.output_dir or data_cfg.get("processed_dir", "data/processed")

    # 1. Extract raw transactional data
    extractor = get_extractor(config)
    raw_data = extractor.extract_raw_data()

    # 2. Build Demand Forecasting Features
    logger.info("Engineering Demand Forecasting feature matrix...")
    ts_cfg = config.get("time_series", {})
    ts_builder = TimeSeriesFeatureBuilder(
        lags=ts_cfg.get("lags", [1, 7, 14, 30]),
        rolling_windows=ts_cfg.get("rolling_windows", [7, 30])
    )
    demand_headers, demand_X, demand_y = ts_builder.build_features(raw_data)

    # 3. Build Stockout Risk Features
    logger.info("Engineering Stockout Risk feature matrix...")
    stock_builder = InventoryStockoutFeatureBuilder()
    stock_headers, stock_X, stock_y = stock_builder.build_features(raw_data)

    # 4. Build Anomaly Detection Features
    logger.info("Engineering Operational Anomaly feature matrix...")
    anomaly_builder = OperationalAnomalyFeatureBuilder()
    anomaly_headers, anomaly_X, anomaly_y = anomaly_builder.build_features(raw_data)

    # 5. Build Courier SLA Features (Phase 2.4 Extension)
    logger.info("Engineering Courier SLA & Delivery Duration feature matrix...")
    courier_builder = CourierSLAFeatureBuilder()
    courier_headers, courier_X, courier_y = courier_builder.build_features(raw_data)

    # 6. Scale features
    scaler = StandardScaler()
    scaled_demand_X = scaler.fit_transform(demand_X)
    scaled_stock_X = scaler.fit_transform(stock_X)
    scaled_anomaly_X = scaler.fit_transform(anomaly_X)

    logger.info("Feature extraction completed successfully:")
    logger.info(f"  • Demand Forecasting Dataset: {demand_X.shape[0]} samples, {demand_X.shape[1]} features")
    logger.info(f"  • Stockout Risk Dataset:      {stock_X.shape[0]} samples, {stock_X.shape[1]} features")
    logger.info(f"  • Operational Anomaly Dataset: {anomaly_X.shape[0]} samples, {anomaly_X.shape[1]} features")
    logger.info(f"  • Courier SLA Dataset:        {courier_X.shape[0]} samples, {courier_X.shape[1]} features")

    if args.dry_run:
        logger.info("Dry-run flag specified. Skipping disk artifact serialization.")
        return

    # 7. Save datasets to disk
    logger.info(f"Serializing processed datasets to '{processed_dir}'...")

    save_csv(os.path.join(processed_dir, "demand_forecasting_features.csv"), demand_headers, scaled_demand_X, demand_y)
    save_pt(os.path.join(processed_dir, "demand_forecasting_features.pt"), scaled_demand_X, demand_y)

    save_csv(os.path.join(processed_dir, "stockout_risk_features.csv"), stock_headers, scaled_stock_X, stock_y)
    save_pt(os.path.join(processed_dir, "stockout_risk_features.pt"), scaled_stock_X, stock_y)

    save_csv(os.path.join(processed_dir, "anomaly_detection_features.csv"), anomaly_headers, scaled_anomaly_X, anomaly_y)
    save_pt(os.path.join(processed_dir, "anomaly_detection_features.pt"), scaled_anomaly_X, anomaly_y)

    save_csv(os.path.join(processed_dir, "courier_sla_features.csv"), courier_headers, courier_X, courier_y)
    save_pt(os.path.join(processed_dir, "courier_sla_features.pt"), courier_X, courier_y)

    logger.info("Phase 1 Feature Extraction Pipeline executed successfully! All artifacts generated.")


if __name__ == "__main__":
    main()
