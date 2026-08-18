import pytest
import torch
from src.model import SimpleMLP

def test_mlp_initialization():
    """Verify SimpleMLP structure and parameter allocation."""
    model = SimpleMLP(input_dim=10, hidden_dims=[20, 30], output_dim=5)
    assert isinstance(model, torch.nn.Module)
    
    # Verify linear layer projection logic
    # Total linear layers: 3 (10->20, 20->30, 30->5)
    linear_layers = [m for m in model.modules() if isinstance(m, torch.nn.Linear)]
    assert len(linear_layers) == 3
    assert linear_layers[0].in_features == 10
    assert linear_layers[0].out_features == 20
    assert linear_layers[2].out_features == 5

def test_mlp_forward_pass_2d():
    """Verify SimpleMLP handles typical 2D (batch, features) inputs."""
    batch_size = 16
    input_dim = 10
    output_dim = 5
    
    model = SimpleMLP(input_dim=input_dim, hidden_dims=[20], output_dim=output_dim)
    x = torch.randn(batch_size, input_dim)
    out = model(x)
    
    assert out.shape == (batch_size, output_dim)
    # Check that gradients flow back (not detached)
    assert out.requires_grad

def test_mlp_forward_pass_multidim():
    """Verify SimpleMLP flattens multidimensional (e.g. image) inputs."""
    batch_size = 8
    # 2D features: 28x28 (e.g., MNIST image)
    model = SimpleMLP(input_dim=784, hidden_dims=[64], output_dim=10)
    x = torch.randn(batch_size, 1, 28, 28)
    out = model(x)
    
    assert out.shape == (batch_size, 10)
