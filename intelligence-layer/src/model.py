from typing import List
import torch
import torch.nn as nn

class SimpleMLP(nn.Module):
    """
    A simple customizable Multi-Layer Perceptron classifier.
    Useful as a baseline or testing model for MLP training.
    """
    def __init__(
        self,
        input_dim: int,
        hidden_dims: List[int],
        output_dim: int,
        dropout: float = 0.0
    ):
        super().__init__()
        
        layers = []
        prev_dim = input_dim
        
        for hidden_dim in hidden_dims:
            layers.append(nn.Linear(prev_dim, hidden_dim))
            layers.append(nn.BatchNorm1d(hidden_dim))
            layers.append(nn.ReLU())
            if dropout > 0.0:
                layers.append(nn.Dropout(dropout))
            prev_dim = hidden_dim
            
        layers.append(nn.Linear(prev_dim, output_dim))
        self.network = nn.Sequential(*layers)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        """
        Forward pass.
        Args:
            x (Tensor): Input tensor of shape (batch_size, input_dim) or (batch_size, ...)
        Returns:
            Tensor: Output logits of shape (batch_size, output_dim)
        """
        # Flatten input if multidimensional (e.g. images)
        if x.dim() > 2:
            x = x.view(x.size(0), -1)
        return self.network(x)
