import os
from typing import List
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    APP_NAME: str = "chia.florist AI Lab API"
    VERSION: str = "1.0.0"
    API_V1_PREFIX: str = "/api/v1"
    MODELS_DIR: str = os.getenv("MODELS_DIR", "models")
    CONFIGS_DIR: str = os.getenv("CONFIGS_DIR", "configs")
    LOG_LEVEL: str = "INFO"
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    ALLOWED_ORIGINS: List[str] = ["*"]

    class Config:
        case_sensitive = True

settings = Settings()
