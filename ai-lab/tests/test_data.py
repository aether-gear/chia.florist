import pytest
import torch
from src.data_loader import SyntheticClassificationDataset, get_dataloaders

def test_synthetic_dataset():
    """Verify basic properties of the synthetic dataset class."""
    num_samples = 100
    input_dim = 10
    num_classes = 3
    
    dataset = SyntheticClassificationDataset(
        num_samples=num_samples,
        input_dim=input_dim,
        num_classes=num_classes
    )
    
    assert len(dataset) == num_samples
    
    features, label = dataset[0]
    assert features.shape == (input_dim,)
    assert isinstance(label, torch.Tensor)
    assert label.ndim == 0  # scalar target
    assert 0 <= label.item() < num_classes

def test_dataloaders_split():
    """Verify that get_dataloaders returns correct batch structure and splits."""
    config = {
        "seed": 123,
        "model": {
            "input_dim": 50,
            "output_dim": 5
        },
        "data": {
            "batch_size": 32,
            "num_workers": 0,
            "train_val_split": 0.8
        }
    }
    
    train_loader, val_loader = get_dataloaders(config)
    
    # Check dataset size
    # Base dataset size = 1000. 80% train = 800, 20% val = 200.
    # Batch size = 32. 800 / 32 = 25 batches. 200 / 32 = 7 batches (last batch of 8)
    assert len(train_loader.dataset) == 800
    assert len(val_loader.dataset) == 200
    
    # Iterate a train batch
    features, targets = next(iter(train_loader))
    assert features.shape == (32, 50)
    assert targets.shape == (32,)
