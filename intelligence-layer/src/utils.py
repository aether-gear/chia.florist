import logging
import random
import numpy as np
import torch

logger = logging.getLogger(__name__)

def setup_logging(level: str = "INFO") -> None:
    """Configures the root logging style."""
    logging.basicConfig(
        level=getattr(logging, level.upper(), logging.INFO),
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

def set_seed(seed: int) -> None:
    """Sets random seeds for reproducibility."""
    random.seed(seed)
    np.random.seed(seed)
    torch.manual_seed(seed)
    if torch.cuda.is_available():
        torch.cuda.manual_seed(seed)
        torch.cuda.manual_seed_all(seed)
        # Guarantees deterministic behavior for PyTorch operations
        torch.backends.cudnn.deterministic = True
        torch.backends.cudnn.benchmark = False
    logger.info(f"Random seed set to {seed}")

def get_device(preferred: str = "cuda") -> torch.device:
    """Resolves and returns the target training device."""
    if preferred == "cuda" and torch.cuda.is_available():
        device = torch.device("cuda")
    elif preferred == "mps" and torch.backends.mps.is_available():
        device = torch.device("mps")
    else:
        device = torch.device("cpu")
    logger.info(f"Using device: {device}")
    return device
