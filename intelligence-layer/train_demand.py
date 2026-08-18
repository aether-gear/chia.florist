import argparse
import os
import sys
import logging
import yaml

from src.utils import setup_logging, set_seed
from src.demand_forecasting_trainer import DemandForecastingTrainer

def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="CLI utility to train Phase 2.1 Demand & Sales Forecasting XGBoost model."
    )
    parser.add_argument(
        "--config",
        type=str,
        default=os.path.join("configs", "demand_forecasting.yaml"),
        help="Path to YAML configuration file."
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Load dataset and validate shapes without running full XGBoost model training."
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

def main():
    args = parse_args()
    config = load_config(args.config)

    setup_logging()
    logger = logging.getLogger("train_demand")
    logger.info("Initializing Phase 2.1 Demand & Sales Forecasting Model Training...")

    set_seed(config.get("seed", 42))

    trainer = DemandForecastingTrainer(config)

    # 1. Load & split dataset
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    if args.dry_run:
        logger.info("Dry-run specified. Data loading and chronological splitting validated successfully!")
        logger.info(f"X_train shape: {X_train.shape}, y_train shape: {y_train.shape}")
        logger.info(f"X_val shape:   {X_val.shape}, y_val shape:   {y_val.shape}")
        return

    # 2. Train model
    model = trainer.train(X_train, y_train, X_val, y_val)

    # 3. Evaluate performance
    metrics = trainer.evaluate(model, X_train, y_train, X_val, y_val)

    # 4. Save feature importance plot
    trainer.save_feature_importance(model, feature_names)

    # 5. Serialize native JSON model checkpoint
    model_path = trainer.save_model(model)

    logger.info(f"Phase 2.1 Training completed successfully! Model checkpoint saved to '{model_path}'.")

if __name__ == "__main__":
    main()
