import argparse
import os
import sys
import yaml
from typing import Dict, Any

from src.utils import setup_logging, set_seed, get_device
from src.data_loader import get_dataloaders
from src.model import SimpleMLP
from src.trainer import Trainer

def parse_args() -> argparse.Namespace:
    """Parses command-line arguments."""
    parser = argparse.ArgumentParser(
        description="CLI utility to train a model in the AI Lab environment."
    )
    parser.add_argument(
        "--config",
        type=str,
        default=os.path.join("configs", "default_config.yaml"),
        help="Path to YAML configuration file."
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Run a single forward and backward pass to test the pipeline and exit."
    )
    parser.add_argument(
        "--device",
        type=str,
        default=None,
        help="Device to run the training on (e.g., 'cuda', 'cpu', 'mps'). Overrides config."
    )
    return parser.parse_args()

def load_config(config_path: str) -> Dict[str, Any]:
    """Loads and parses the YAML configuration file."""
    if not os.path.exists(config_path):
        print(f"Error: Configuration file not found at '{config_path}'")
        sys.exit(1)
        
    with open(config_path, "r") as f:
        try:
            return yaml.safe_load(f)
        except yaml.YAMLError as e:
            print(f"Error parsing YAML config: {e}")
            sys.exit(1)

def main() -> None:
    # 1. Parse CLI arguments
    args = parse_args()
    
    # 2. Load configuration
    config = load_config(args.config)
    
    # 3. Setup logging
    setup_logging()
    import logging
    logger = logging.getLogger("train_entry")
    logger.info(f"Loaded config from '{args.config}'")
    
    # 4. Set reproducibility seeds
    set_seed(config.get("seed", 42))
    
    # 5. Resolve hardware device
    device_pref = args.device or config.get("device", "cuda")
    device = get_device(device_pref)
    
    # 6. Initialize DataLoaders
    logger.info("Initializing Datasets and DataLoaders...")
    train_loader, val_loader = get_dataloaders(config)
    
    # 7. Initialize Model
    model_cfg = config.get("model", {})
    logger.info(
        f"Initializing SimpleMLP model with input_dim={model_cfg.get('input_dim')}, "
        f"hidden_dims={model_cfg.get('hidden_dims')}, "
        f"output_dim={model_cfg.get('output_dim')}..."
    )
    model = SimpleMLP(
        input_dim=model_cfg.get("input_dim", 784),
        hidden_dims=model_cfg.get("hidden_dims", [128, 64]),
        output_dim=model_cfg.get("output_dim", 10),
        dropout=model_cfg.get("dropout", 0.0)
    )
    
    # 8. Initialize Trainer
    trainer = Trainer(
        model=model,
        train_loader=train_loader,
        val_loader=val_loader,
        config=config,
        device=device
    )
    
    # 9. Run Dry-Run or Full Training
    if args.dry_run:
        success = trainer.dry_run()
        sys.exit(0 if success else 1)
    else:
        trainer.fit()

if __name__ == "__main__":
    main()
