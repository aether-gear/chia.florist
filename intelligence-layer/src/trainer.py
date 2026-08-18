import os
import logging
from typing import Dict, Any
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader
from tqdm import tqdm

logger = logging.getLogger(__name__)

class Trainer:
    """
    Standard training & validation runner.
    Orchestrates the training process, evaluates model performance, and manages checkpoints.
    """
    def __init__(
        self,
        model: nn.Module,
        train_loader: DataLoader,
        val_loader: DataLoader,
        config: Dict[str, Any],
        device: torch.device
    ):
        self.model = model.to(device)
        self.train_loader = train_loader
        self.val_loader = val_loader
        self.config = config
        self.device = device
        
        # Training config parameters
        train_cfg = config.get("training", {})
        self.epochs = train_cfg.get("epochs", 5)
        self.lr = train_cfg.get("learning_rate", 0.001)
        self.weight_decay = train_cfg.get("weight_decay", 0.0001)
        self.val_every = train_cfg.get("val_every_n_epochs", 1)
        self.save_dir = train_cfg.get("save_dir", "models")
        self.log_interval = train_cfg.get("log_interval", 10)
        
        # Loss and Optimizer
        self.criterion = nn.CrossEntropyLoss()
        self.optimizer = optim.AdamW(
            self.model.parameters(),
            lr=self.lr,
            weight_decay=self.weight_decay
        )

        os.makedirs(self.save_dir, exist_ok=True)

    def train_epoch(self, epoch: int) -> float:
        """Trains the model for one epoch."""
        self.model.train()
        total_loss = 0.0
        correct = 0
        total = 0
        
        progress_bar = tqdm(self.train_loader, desc=f"Epoch {epoch}/{self.epochs} [Train]")
        
        for batch_idx, (features, targets) in enumerate(progress_bar):
            features, targets = features.to(self.device), targets.to(self.device)
            
            self.optimizer.zero_grad()
            outputs = self.model(features)
            loss = self.criterion(outputs, targets)
            loss.backward()
            self.optimizer.step()
            
            total_loss += loss.item() * features.size(0)
            _, predicted = outputs.max(1)
            total += targets.size(0)
            correct += predicted.eq(targets).sum().item()
            
            # Update progress bar metrics
            if batch_idx % self.log_interval == 0:
                acc = 100.0 * correct / total
                progress_bar.set_postfix({"Loss": f"{loss.item():.4f}", "Acc": f"{acc:.2f}%"})
                
        epoch_loss = total_loss / total
        epoch_acc = 100.0 * correct / total
        logger.info(f"Epoch {epoch} - Train Loss: {epoch_loss:.4f} | Train Acc: {epoch_acc:.2f}%")
        return epoch_loss

    @torch.no_grad()
    def evaluate(self, description: str = "[Val]") -> tuple[float, float]:
        """Evaluates the model on validation dataset."""
        self.model.eval()
        total_loss = 0.0
        correct = 0
        total = 0
        
        for features, targets in self.val_loader:
            features, targets = features.to(self.device), targets.to(self.device)
            outputs = self.model(features)
            loss = self.criterion(outputs, targets)
            
            total_loss += loss.item() * features.size(0)
            _, predicted = outputs.max(1)
            total += targets.size(0)
            correct += predicted.eq(targets).sum().item()
            
        val_loss = total_loss / total
        val_acc = 100.0 * correct / total
        logger.info(f"{description} Loss: {val_loss:.4f} | Acc: {val_acc:.2f}%")
        return val_loss, val_acc

    def fit(self) -> None:
        """Runs the main training and validation loop."""
        logger.info("Starting training loop...")
        best_val_acc = -1.0
        
        for epoch in range(1, self.epochs + 1):
            _ = self.train_epoch(epoch)
            
            if epoch % self.val_every == 0 or epoch == self.epochs:
                val_loss, val_acc = self.evaluate()
                
                # Checkpoint saving
                if val_acc > best_val_acc:
                    best_val_acc = val_acc
                    self.save_checkpoint(epoch, "best_model.pt")
                    
            # Save periodic checkpoint
            self.save_checkpoint(epoch, "latest_model.pt")
            
        logger.info(f"Training completed. Best Val Acc: {best_val_acc:.2f}%")

    def dry_run(self) -> bool:
        """Executes a single forward and backward step to verify pipeline integration."""
        logger.info("Starting dry-run sanity check...")
        try:
            self.model.train()
            features, targets = next(iter(self.train_loader))
            features, targets = features.to(self.device), targets.to(self.device)
            
            self.optimizer.zero_grad()
            outputs = self.model(features)
            loss = self.criterion(outputs, targets)
            loss.backward()
            self.optimizer.step()
            
            logger.info("Dry-run validation forward pass...")
            self.model.eval()
            with torch.no_grad():
                val_features, _ = next(iter(self.val_loader))
                val_features = val_features.to(self.device)
                _ = self.model(val_features)
                
            logger.info("Dry-run passed successfully!")
            return True
        except Exception as e:
            logger.error(f"Dry-run failed: {e}", exc_info=True)
            return False

    def save_checkpoint(self, epoch: int, filename: str) -> None:
        """Saves current state dicts to checkpoint file."""
        filepath = os.path.join(self.save_dir, filename)
        checkpoint = {
            "epoch": epoch,
            "model_state_dict": self.model.state_dict(),
            "optimizer_state_dict": self.optimizer.state_dict(),
            "config": self.config
        }
        torch.save(checkpoint, filepath)
        logger.debug(f"Saved checkpoint to {filepath}")
