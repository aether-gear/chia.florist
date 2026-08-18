import os
import logging
from typing import Dict, Any, Optional
import xgboost as xgb
import joblib

from app.config import settings

logger = logging.getLogger("model_loader")

class ModelRegistry:
    """
    Singleton Model Registry that loads model checkpoints into memory at server startup.
    """
    _instance: Optional["ModelRegistry"] = None

    def __init__(self, models_dir: Optional[str] = None):
        self.models_dir = models_dir or settings.MODELS_DIR
        self.models: Dict[str, Any] = {}
        self.model_status: Dict[str, str] = {}

    @classmethod
    def get_instance(cls) -> "ModelRegistry":
        if cls._instance is None:
            cls._instance = ModelRegistry()
        return cls._instance

    def load_all_models(self) -> None:
        logger.info(f"Scanning directory '{self.models_dir}' for model checkpoints...")

        # 1. Demand Forecasting Model
        demand_path = os.path.join(self.models_dir, "demand_forecasting.json")
        if os.path.exists(demand_path):
            try:
                model = xgb.XGBRegressor()
                model.load_model(demand_path)
                self.models["demand"] = model
                self.model_status["demand"] = f"loaded ({demand_path})"
                logger.info("Successfully loaded Demand Forecasting XGBoost model.")
            except Exception as e:
                logger.error(f"Failed to load Demand Forecasting model: {e}")
                self.model_status["demand"] = f"error ({e})"
        else:
            self.model_status["demand"] = "missing"

        # 2. Stockout Risk Model
        stockout_path = os.path.join(self.models_dir, "stockout_risk.json")
        if os.path.exists(stockout_path):
            try:
                model = xgb.XGBClassifier()
                model.load_model(stockout_path)
                self.models["stockout"] = model
                self.model_status["stockout"] = f"loaded ({stockout_path})"
                logger.info("Successfully loaded Stockout Risk XGBoost model.")
            except Exception as e:
                logger.error(f"Failed to load Stockout Risk model: {e}")
                self.model_status["stockout"] = f"error ({e})"
        else:
            self.model_status["stockout"] = "missing"

        # 3. Courier SLA Model
        courier_path = os.path.join(self.models_dir, "courier_sla.json")
        if os.path.exists(courier_path):
            try:
                model = xgb.XGBRegressor()
                model.load_model(courier_path)
                self.models["courier"] = model
                self.model_status["courier"] = f"loaded ({courier_path})"
                logger.info("Successfully loaded Courier SLA XGBoost model.")
            except Exception as e:
                logger.error(f"Failed to load Courier SLA model: {e}")
                self.model_status["courier"] = f"error ({e})"
        else:
            self.model_status["courier"] = "missing"

        # 4. Anomaly Detection Model (Isolation Forest)
        anomaly_path = os.path.join(self.models_dir, "anomaly_detector", "isolation_forest.pkl")
        if os.path.exists(anomaly_path):
            try:
                model = joblib.load(anomaly_path)
                self.models["anomaly"] = model
                self.model_status["anomaly"] = f"loaded ({anomaly_path})"
                logger.info("Successfully loaded Isolation Forest Anomaly model.")
            except Exception as e:
                logger.error(f"Failed to load Isolation Forest Anomaly model: {e}")
                self.model_status["anomaly"] = f"error ({e})"
        else:
            self.model_status["anomaly"] = "missing"

    def get_model(self, name: str) -> Any:
        if name not in self.models:
            status = self.model_status.get(name, "unknown")
            raise RuntimeError(f"Model '{name}' is not available. Status: {status}")
        return self.models[name]

    def get_loaded_names(self) -> list[str]:
        return list(self.models.keys())

    def get_status_summary(self) -> Dict[str, str]:
        return self.model_status

model_registry = ModelRegistry.get_instance()
