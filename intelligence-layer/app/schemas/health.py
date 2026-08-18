from typing import List, Dict
from pydantic import BaseModel

class HealthResponse(BaseModel):
    status: str
    version: str
    loaded_models: List[str]
    model_details: Dict[str, str]
