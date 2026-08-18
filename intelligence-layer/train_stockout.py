import argparse
import os
import sys
import logging
import yaml

from src.utils import setup_logging, set_seed
from src.stockout_risk_trainer import StockoutRiskTrainer

def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="CLI utility to train Phase 2.2 Inventory Stockout Risk XGBoost classifier."
    )
    parser.add_argument(
        "--config",
        type=str,
        default=os.path.join("configs", "stockout_risk.yaml"),
        help="Path to YAML configuration file."
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Load dataset and validate split shapes without running full model training."
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
    logger = logging.getLogger("train_stockout")
    logger.info("Initializing Phase 2.2 Inventory Stockout Risk Classifier Training...")

    set_seed(config.get("seed", 42))

    trainer = StockoutRiskTrainer(config)

    # 1. Load & split dataset
    X_train, y_train, X_val, y_val, feature_names = trainer.load_and_split_data()

    if args.dry_run:
        logger.info("Dry-run specified. Dataset split validated successfully!")
        logger.info(f"X_train shape: {X_train.shape}, y_train pos: {sum(y_train)}, neg: {len(y_train)-sum(y_train)}")
        logger.info(f"X_val shape:   {X_val.shape}, y_val pos:   {sum(y_val)}, neg: {len(y_val)-sum(y_val)}")
        return

    # 2. Apply SMOTE on training split
    X_train_res, y_train_res = trainer.apply_smote(X_train, y_train)

    # 3. Train model
    model = trainer.train(X_train_res, y_train_res, X_val, y_val)

    # 4. Evaluate performance
    metrics = trainer.evaluate(model, X_train_res, y_train_res, X_val, y_val)

    # 5. Save feature importance plot
    trainer.save_feature_importance(model, feature_names)

    # 6. Serialize native JSON model checkpoint
    model_path = trainer.save_model(model)

    logger.info(f"Phase 2.2 Training completed successfully! Model checkpoint saved to '{model_path}'.")

if __name__ == "__main__":
    main()
