# AI Training and Testing Lab

This repository contains a modular and production-ready directory structure for training and testing AI/ML models. It is designed to be easily extensible for new datasets, architectures, and experiments.

## Project Structure

```text
ai-lab/
├── configs/                # Hyperparameters and experiment configuration files
├── data/                   # Data directory (split into raw and processed)
├── models/                 # Saved model weights and checkpoints
├── notebooks/              # Jupyter Notebooks for exploration and EDA
├── src/                    # Source code package for ML pipeline
│   ├── data_loader.py      # Dataset definition & batching
│   ├── model.py            # Neural network architecture definition
│   ├── trainer.py          # Training loop, evaluation, and logging
│   └── utils.py            # Random seed seeding, logs, and helper methods
├── tests/                  # Unit and integration tests for model and data pipelines
├── train.py                # Command-line script to run training
├── pyproject.toml          # Tooling configurations (pytest, black, ruff)
└── requirements.txt        # Python dependency list
```

## Setup Instructions

1. **Create Virtual Environment**:
   ```bash
   python -m venv .venv
   ```

2. **Activate Virtual Environment**:
   - **Windows (PowerShell)**:
     ```powershell
     .venv\Scripts\Activate.ps1
     ```
   - **macOS/Linux**:
     ```bash
     source .venv/bin/activate
     ```

3. **Install Dependencies**:
   ```bash
   pip install -r requirements.txt
   ```

## Running Training

To run model training with the default configuration:
```bash
python train.py --config configs/default_config.yaml
```

To run a dry run to check system setup:
```bash
python train.py --dry-run
```

## Running Tests

Verify the neural network architecture and data loader outputs using `pytest`:
```bash
pytest
```

## Versioning & Commit Guidelines

`intelligence-layer` releases are automated via GitHub Actions (Current baseline: `intelligence-layer-v0.4.0`).
FastAPI Docker images are automatically built and published to GHCR upon release.

| Commit Type | Version Bump | Example |
| :--- | :--- | :--- |
| `fix:` / `perf:` | **Patch** (`0.4.0` → `0.4.1`) | `fix(sku): resolve edge case in demand anomaly scorer` |
| `feat:` | **Minor** (`0.4.0` → `0.5.0`) | `feat(ml): integrate chat-to-generate flower board model` |
| `feat!:` / `BREAKING CHANGE:` | **Major** (or minor pre-1.0) | `feat!: update inference input payload contract to v3` |
| `docs:` / `test:` / `chore:` | *No release* | `docs: document 4 ML models for workflows` |

See root [docs/VERSIONING_AND_RELEASES.md](../docs/VERSIONING_AND_RELEASES.md) for full monorepo guidelines.
