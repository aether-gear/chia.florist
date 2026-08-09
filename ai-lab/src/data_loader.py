from typing import Tuple, Dict, Any
import torch
from torch.utils.data import Dataset, DataLoader, random_split

class SyntheticClassificationDataset(Dataset):
    """
    A synthetic classification dataset.
    Generates random feature vectors and mock labels for testing and debugging.
    """
    def __init__(self, num_samples: int = 1000, input_dim: int = 784, num_classes: int = 10):
        # Generate random inputs centered around 0 with unit variance
        self.data = torch.randn(num_samples, input_dim)
        
        # Generate target labels. We construct a simple mapping so the model has something to learn
        # Mock projection matrix
        projection = torch.randn(input_dim, num_classes)
        raw_scores = self.data @ projection
        self.labels = torch.argmax(raw_scores, dim=1)

    def __len__(self) -> int:
        return len(self.data)

    def __getitem__(self, idx: int) -> Tuple[torch.Tensor, torch.Tensor]:
        return self.data[idx], self.labels[idx]

def get_dataloaders(config: Dict[str, Any]) -> Tuple[DataLoader, DataLoader]:
    """
    Constructs train and validation DataLoader instances.
    Supports configuring split, batch size, and worker threads.
    """
    # Extract settings from the nested configuration structure
    data_config = config.get("data", {})
    model_config = config.get("model", {})
    
    input_dim = model_config.get("input_dim", 784)
    output_dim = model_config.get("output_dim", 10)
    batch_size = data_config.get("batch_size", 64)
    num_workers = data_config.get("num_workers", 0)
    train_val_split = data_config.get("train_val_split", 0.8)
    
    # Instantiate the dataset
    full_dataset = SyntheticClassificationDataset(
        num_samples=1000,
        input_dim=input_dim,
        num_classes=output_dim
    )
    
    # Compute splits
    train_size = int(len(full_dataset) * train_val_split)
    val_size = len(full_dataset) - train_size
    
    # Split
    train_dataset, val_dataset = random_split(
        full_dataset, 
        [train_size, val_size],
        generator=torch.Generator().manual_seed(config.get("seed", 42))
    )
    
    # Create DataLoaders
    train_loader = DataLoader(
        train_dataset,
        batch_size=batch_size,
        shuffle=True,
        num_workers=num_workers,
        pin_memory=torch.cuda.is_available()
    )
    
    val_loader = DataLoader(
        val_dataset,
        batch_size=batch_size,
        shuffle=False,
        num_workers=num_workers,
        pin_memory=torch.cuda.is_available()
    )
    
    return train_loader, val_loader
